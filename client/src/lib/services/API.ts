import { SocketConnection } from './API/SocketConnection';
import { REST } from './API/REST';
import { handleError } from '$lib/helpers/handleError';

export class API {
  private static socketConnection: SocketConnection | null = null;

  private static connectToSocket = async (gameId: string) : Promise<SocketConnection> => {
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

  public static getSocketConnection = handleError(API.connectToSocket);
  public static getGame = handleError(REST.getGame);
  public static getPlayer = handleError(REST.getPlayer);
  public static patchGame = handleError(REST.patchGame);
  public static postGame = handleError(REST.postGame);
  public static patchPlayer = handleError(REST.patchPlayer);
  public static postPlayer = handleError(REST.postPlayer);
}