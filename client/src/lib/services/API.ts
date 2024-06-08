import { SocketConnection } from './API/SocketConnection';
import { REST } from './API/REST';

export class API {
  private static socketConnection: SocketConnection | null = null;

  public static getSocketConnection = async (gameId: string) : Promise<SocketConnection> => {
    if (API.socketConnection) {
      if (API.socketConnection.gameId === gameId) return API.socketConnection;

      API.socketConnection.close();
    };

    const player = await REST.getPlayer();
    if (!player) throw new Error('No player found.');

    API.socketConnection = new SocketConnection(gameId, player.id);
    await API.socketConnection.start();

    return API.socketConnection;
  };

  public static getGame = REST.getGame;
  public static getPlayer = REST.getPlayer;
  public static patchGame = REST.patchGame;
  public static postGame = REST.postGame;
  public static patchPlayer = REST.patchPlayer;
  public static postPlayer = REST.postPlayer;
}