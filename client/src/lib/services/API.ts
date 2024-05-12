import { type GameConfig } from '$lib/types/GameConfig';
import { type Game } from '$lib/types/Game';
import { type Player } from '$lib/types/Player';
import { type ClientMessage } from '$lib/types/ClientMessage';
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

    // If we don't have a player yet, create one.
    if (!playerId) return this.setPlayer(await postPlayer(playerName));

    // If we already have a player, try to reuse that id.
    try {
      return this.setPlayer(await patchPlayer(playerName, playerId));
    } catch (error) {
      return this.setPlayer(await postPlayer(playerName));
    }
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

  public static sendSocketMessage(message: Omit<ClientMessage, 'playerId'>) {
    SocketConnection.send({
      ...message,
      playerId: this.player?.id || null,
    });
  }

  public static startSocketConnection = SocketConnection.start;
  public static closeSocketConnection = SocketConnection.close;
  public static onSocketMessage = SocketConnection.onMessage;
  public static onSocketOpen = SocketConnection.onOpen;
  public static onSocketClose = SocketConnection.onClose;
  public static onSocketError = SocketConnection.onError;
  public static getGame = getGame;
}