import { WebSocketCloseCode } from '$lib/enums/WebSocketCloseCode';
import { type GenericEventCallback, type CloseEventCallback, type MessageEventCallback } from '$lib/types/EventCallbacks';
import type { ClientMessage } from '$lib/types/ClientMessage';
import { ClientMessageType } from '$lib/enums/ClientMessageType';

const RECOVERABLE_CLOSE_CODES = [
  WebSocketCloseCode.GOING_AWAY,
  WebSocketCloseCode.ABNORMAL_CLOSURE,
  WebSocketCloseCode.MESSAGE_TOO_BIG,
  WebSocketCloseCode.MISSING_EXTENSION,
  WebSocketCloseCode.INTERNAL_ERROR,
  WebSocketCloseCode.SERVICE_RESTART,
  WebSocketCloseCode.TRY_AGAIN_LATER,
];
const DEFAULT_RETRIES = 5;
const RETRY_MAX_WAIT_SECONDS = 180;

export class SocketConnection {
  private connected = false;
  private connecting = false;
  private retries = DEFAULT_RETRIES;

  private host: string | undefined = undefined;
  private connectionProtocol: string | undefined = undefined;
  private connection: WebSocket | undefined = undefined;

  private onOpenCallbacks: Array<GenericEventCallback> = [];
  private onCloseCallbacks: Array<CloseEventCallback> = [];
  private onErrorCallbacks: Array<GenericEventCallback> = [];
  private onMessageCallbacks: Array<MessageEventCallback> = [];

  private queuedMessages: Array<ClientMessage> = [];

  constructor(public readonly gameId: string, private playerId: string) {
    this.connectionProtocol = window.location.protocol.endsWith('s:') ? 'wss:' : 'ws:';
    this.host = window.location.host;

    this.baseOnOpen = this.baseOnOpen.bind(this);
    this.baseOnClose = this.baseOnClose.bind(this);
    this.baseOnError = this.baseOnError.bind(this);
    this.baseOnMessage = this.baseOnMessage.bind(this);
  }

  private async awaitRetries(retries: number) {
    if (retries === DEFAULT_RETRIES) return;

    const exponentialWaitInSeconds = RETRY_MAX_WAIT_SECONDS * (1 / retries ** 2);

    return new Promise<void>((resolve) => setTimeout(resolve, exponentialWaitInSeconds));
  }

  public async start(retries?: number) {
    if (this.connected || this.connecting) return;
    this.connecting = true;

    if (retries === 0) return;
    this.retries = retries ?? DEFAULT_RETRIES;
    await this.awaitRetries(this.retries);

    this.connection = new WebSocket(`${this.connectionProtocol}//${this.host}/api/v1/ws/${this.gameId}`);
    this.connection.addEventListener('open', this.baseOnOpen);
    this.connection.addEventListener('close', this.baseOnClose);
    this.connection.addEventListener('error', this.baseOnError);
    this.connection.addEventListener('message', this.baseOnMessage);

    this.send({ type: ClientMessageType.Join, payload: null });
  }

  public close() {
    if (!this.connected || !this.connection) return;
    this.connection.close();
    this.connected = this.connecting = false;
  }

  public send(message: Omit<ClientMessage, 'playerId'>) {
    const messageWithPlayerId = {
      ...message,
      playerId: this.playerId,
    };
    if (!this.connected && !this.connecting) this.start();
    if (this.connected && this.connection) this.connection.send(JSON.stringify(messageWithPlayerId));
    else this.queuedMessages.push(messageWithPlayerId);
  }

  private sendQueuedMessages() {
    const messages = structuredClone(this.queuedMessages);

    // Empty the array of queued messages.
    this.queuedMessages.splice(0, this.queuedMessages.length);

    messages.forEach((message) => this.send(message));
  }

  public onOpen(callback: GenericEventCallback) {
    this.onOpenCallbacks.push(callback);
  }

  private baseOnOpen(event: Event) {
    this.connecting = false;
    this.connected = true;
    this.onOpenCallbacks.forEach(callback => callback(event));

    this.sendQueuedMessages();
  }

  public onClose(callback: CloseEventCallback) {
    this.onCloseCallbacks.push(callback);
  }

  private baseOnClose(event: CloseEvent) {
    this.connected = false;
    this.onCloseCallbacks.forEach(callback => callback(event));

    // Try to reconnect if the connection seems recoverable.
    if (event.code in RECOVERABLE_CLOSE_CODES) this.start(this.retries--);
  }

  public onError(callback: GenericEventCallback) {
    this.onErrorCallbacks.push(callback);
  }

  private baseOnError(event: Event) {
    this.onErrorCallbacks.forEach(callback => callback(event));
  }

  public onMessage(callback: MessageEventCallback) {
    this.onMessageCallbacks.push(callback);
  }

  private baseOnMessage(event: MessageEvent) {
    this.onMessageCallbacks.forEach(callback => callback({
      gameId: this.gameId,
      ...JSON.parse(event.data)
    }));
  }
}