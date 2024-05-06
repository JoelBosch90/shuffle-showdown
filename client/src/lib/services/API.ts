import { type GameConfig } from '$lib/types/GameConfig';
import getGame from './API/getGame';
import patchGame from './API/patchGame';
import postGame from './API/postGame';
import patchPlayer from './API/patchPlayer';
import { SocketConnection } from './API/SocketConnection';

export class API {
  private static playerId: string | null = null;

  private static setPlayerId(playerId: string) {
    this.playerId = playerId;
    localStorage.setItem('playerId', playerId);
  }

  public static getPlayerId() {
    if (!this.playerId) {
      this.playerId = localStorage.getItem('playerId');
    }

    return this.playerId;
  }

  public static async postGame(playListString: string) {
    const { game } = await postGame(playListString);
    this.setPlayerId(game?.owner?.id);

    return { game };
  };

  public static async patchGame(settings: GameConfig) {
    const playerId = this.getPlayerId();

    if (!playerId) {
      throw new Error('No player ID found');
    }

    return patchGame(settings, playerId);
  }

  public static patchPlayer(playerName: string) {
    const playerId = this.getPlayerId();

    if (!playerId) {
      throw new Error('No player ID found');
    }

    return patchPlayer(playerName, playerId);
  }

  public static getGame = getGame;
  public static SocketConnection = SocketConnection;
}