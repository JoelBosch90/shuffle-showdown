import { type GameConfig } from '$lib/types/GameConfig';
import { type Game } from '$lib/types/Game';
import { type Player } from '$lib/types/Player';
import getGame from './API/game/get';
import patchGame from './API/game/patch';
import postGame from './API/game/post';
import postPlayer from './API/player/post';
import patchPlayer from './API/player/patch';
import { SocketConnection } from './API/SocketConnection';

const PLAYER_NAME = 'playerName';
const PLAYER_ID = 'playerId';

export class API {
  private static player: Player | null = null;

  private static setPlayer(player: Player) : Player {
    this.player = player;
    const { id, name } = player;

    localStorage.setItem(PLAYER_ID, id);
    if (name) localStorage.setItem(PLAYER_NAME, name);

    return this.player;
  }

  private static async getPlayerId() : Promise<string> {
    const player = await this.getPlayer();

    if (!player?.id) return '';

    return player.id;
  }

  public static async getPlayer() : Promise<Player | null> {
    if (this.player) return this.player;
  
    const id = localStorage.getItem('playerId');

    if (!id) return null;

    this.player = {
      id,
      // Convert null to undefined if no name is set.
      name: localStorage.getItem('playerName') || undefined,
    };

    return this.player;
  }

  public static async postPlayer(playerName: string) : Promise<Player> {
    const playerId = await this.getPlayerId();

    // If we already have a playerId, we can reuse it an simply rename that player.
    const call = playerId ? patchPlayer(playerName, playerId) : postPlayer(playerName);
    const player = await call;

    return this.setPlayer(player);
  }

  public static async patchPlayer(playerName: string) : Promise<Player> {
    const playerId = await this.getPlayerId();

    if (!playerId) {
      throw new Error('No player ID found');
    }

    return patchPlayer(playerName, playerId);
  }

  public static async postGame(playListString: string) : Promise<Player> {
    const game = await postGame(playListString);

    this.setPlayer({ id: game?.owner?.id });

    return game;
  };

  public static async patchGame(settings: GameConfig) : Promise<Game> {
    const playerId = await this.getPlayerId();

    if (!playerId) {
      throw new Error('No player ID found');
    }

    return patchGame(settings, playerId);
  }

  public static getGame = getGame;
  public static SocketConnection = SocketConnection;
}