import { WebSocketCloseCode } from '$lib/enums/WebSocketCloseCode';
import { type GenericEventCallback, type CloseEventCallback, type MessageEventCallback } from '$lib/types/EventCallbacks';
import type { ClientMessage } from '$lib/types/ClientMessage';
import { ClientMessageType } from '$lib/enums/ClientMessageType';
import { eventToServerMessage } from '$lib/mappers/EventToServerMessage';
import { API } from '../API';

const RECOVERABLE_CLOSE_CODES = [
  WebSocketCloseCode.GOING_AWAY,
  WebSocketCloseCode.ABNORMAL_CLOSURE,
  WebSocketCloseCode.MESSAGE_TOO_BIG,
  WebSocketCloseCode.MISSING_EXTENSION,
  WebSocketCloseCode.INTERNAL_ERROR,
  WebSocketCloseCode.SERVICE_RESTART,
  WebSocketCloseCode.TRY_AGAIN_LATER,
];

export class SocketConnection {
  private static connected = false;
  private static connecting = false;
  private static initialized = false;

  private static gameId: string | undefined = undefined;
  private static host: string | undefined = undefined;
  private static connectionProtocol: string | undefined = undefined;
  private static connection: WebSocket | undefined = undefined;

  private static onOpenCallbacks: Array<GenericEventCallback> = [];
  private static onCloseCallbacks: Array<CloseEventCallback> = [];
  private static onErrorCallbacks: Array<GenericEventCallback> = [];
  private static onMessageCallbacks: Array<MessageEventCallback> = [];

  private static queuedMessages: Array<ClientMessage> = [];

  public static start(gameId?: string) {
    if (gameId && gameId !== SocketConnection.gameId) {
      SocketConnection.close();
      SocketConnection.gameId = gameId;
    }
    if (SocketConnection.connected || SocketConnection.connecting) return;
    if (!SocketConnection.initialized) SocketConnection.initialize();
    SocketConnection.connecting = true;
    SocketConnection.connection = new WebSocket(`${SocketConnection.connectionProtocol}//${SocketConnection.host}/api/v1/ws/${gameId}`);

    SocketConnection.connection.addEventListener('open', SocketConnection.baseOnOpen);
    SocketConnection.connection.addEventListener('close', SocketConnection.baseOnClose);
    SocketConnection.connection.addEventListener('error', SocketConnection.baseOnError);
    SocketConnection.connection.addEventListener('message', SocketConnection.baseOnMessage);

    API.getPlayer().then((player) => {
      // Identify the player to the server.
      SocketConnection.send({
        type: ClientMessageType.Join,
        payload: null,
        playerId: player?.id || null,
      })
    });
  }

  private static initialize() {
    if (SocketConnection.initialized) return;
    SocketConnection.connectionProtocol = window.location.protocol.endsWith('s:') ? 'wss:' : 'ws:';
    SocketConnection.host = window.location.host;
  }

  public static close() {
    if (!SocketConnection.connected || !SocketConnection.connection) return;
    SocketConnection.connection.close();
    SocketConnection.connected = SocketConnection.connecting = false;
  }

  public static send(message: ClientMessage) {
    if (!SocketConnection.connected && !SocketConnection.connecting) SocketConnection.start();
    if (SocketConnection.connected && SocketConnection.connection) SocketConnection.connection.send(JSON.stringify(message));
    else SocketConnection.queuedMessages.push(message);
  }

  private static sendQueuedMessages() {
    const messages = structuredClone(SocketConnection.queuedMessages);

    // Empty the array of queued messages.
    SocketConnection.queuedMessages.splice(0, SocketConnection.queuedMessages.length);

    messages.forEach((message) => SocketConnection.send(message));
  }

  public static onOpen(callback: GenericEventCallback) {
    SocketConnection.onOpenCallbacks.push(callback);
  }

  private static baseOnOpen(event: Event) {
    SocketConnection.connecting = false;
    SocketConnection.connected = true;
    SocketConnection.onOpenCallbacks.forEach(callback => callback(event));

    SocketConnection.sendQueuedMessages();
  }

  public static onClose(callback: CloseEventCallback) {
    SocketConnection.onCloseCallbacks.push(callback);
  }

  private static baseOnClose(event: CloseEvent) {
    SocketConnection.connected = false;
    SocketConnection.onCloseCallbacks.forEach(callback => callback(event));
    
    // Try to reconnect if the connection seems recoverable.
    if (event.code in RECOVERABLE_CLOSE_CODES) SocketConnection.start();
  }

  public static onError(callback: GenericEventCallback) {
    SocketConnection.onErrorCallbacks.push(callback);
  }

  private static baseOnError(event: Event) {
    SocketConnection.onErrorCallbacks.forEach(callback => callback(event));
  }

  public static onMessage(callback: MessageEventCallback) {
    SocketConnection.onMessageCallbacks.push(callback);
  }

  private static baseOnMessage(event: MessageEvent) {
    SocketConnection.onMessageCallbacks.forEach(callback => callback(eventToServerMessage(event)));
  }
}