import { SocketConnection } from './API/SocketConnection';
import getGame from './API/getGame';
import patchGame from './API/patchGame';
import postGame from './API/postGame';

export class API {
  public static postGame = postGame;
  public static getGame = getGame;
  public static patchGame = patchGame;
  public static SocketConnection = SocketConnection;
}