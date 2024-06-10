import { type GameConfig } from '$lib/types/GameConfig';
import { type GameUpdate } from '$lib/types/GameUpdate';
import { type Player } from '$lib/types/Player';
import getGame from './REST/game/get';
import patchGame from './REST/game/patch';
import postGame from './REST/game/post';
import postPlayer from './REST/player/post';
import patchPlayer from './REST/player/patch';

const PLAYER_NAME = 'playerName';
const PLAYER_ID = 'playerId';

export class REST {
  private static player: Player | null = null;

  private static setPlayer(player: Player) : Player {
    REST.player = player;
    const { id, name } = player;

    localStorage.setItem(PLAYER_ID, id);
    if (name) localStorage.setItem(PLAYER_NAME, name);

    return REST.player;
  }

  public static async getPlayerId() : Promise<string> {
    const player = await REST.getPlayer();

    if (!player?.id) return '';

    return player.id;
  }

  public static async getPlayer() : Promise<Player | null> {
    if (REST.player) return REST.player;

    const id = localStorage.getItem('playerId');

    if (!id) return null;

    REST.player = {
      id,
      // Convert null to undefined if no name is set.
      name: localStorage.getItem('playerName') || undefined,
    };

    return REST.player;
  }

  public static async postPlayer(playerName: string) : Promise<Player> {
    const playerId = await REST.getPlayerId();

    // If we don't have a player yet, create one.
    if (!playerId) return REST.setPlayer(await postPlayer(playerName));

    // If we already have a player, try to reuse that id.
    try {
      return REST.setPlayer(await patchPlayer(playerName, playerId));
    } catch (error) {
      return REST.setPlayer(await postPlayer(playerName));
    }
  }

  public static async patchPlayer(playerName: string) : Promise<Player> {
    const playerId = await REST.getPlayerId();

    if (!playerId) {
      return REST.postPlayer(playerName);
    }

    try {
      return await patchPlayer(playerName, playerId);
    } catch (error) {
      return await REST.postPlayer(playerName);
    }
  }

  public static async postGame(playListString: string) : Promise<Player> {
    const playerId = await REST.getPlayerId();
    const game = await postGame(playListString, playerId || undefined);

    REST.setPlayer({ id: game?.owner?.id });

    return game;
  };

  public static async patchGame(settings: GameConfig) : Promise<GameUpdate> {
    const playerId = await REST.getPlayerId();

    if (!playerId) {
      throw new Error('No player ID found');
    }

    return patchGame(settings, playerId);
  }

  public static getGame = getGame;
}