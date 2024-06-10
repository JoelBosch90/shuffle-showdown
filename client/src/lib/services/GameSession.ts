import { goto } from '$app/navigation';
import { API } from '$lib/services/API';
import { type Player } from '$lib/types/Player';
import { type GameUpdate } from '$lib/types/GameUpdate';
import { type GameSessionUpdate } from '$lib/types/GameSessionUpdate';
import { type Answer } from '$lib/types/Answer';
import { isPlayerKickedMessage, isGameSessionUpdateMessage, type ServerMessage, type GameSessionUpdateMessage, isErrorMessage } from '$lib/types/ServerMessage';
import { ClientMessageType } from '$lib/enums/ClientMessageType';
import type { SocketConnection } from './API/SocketConnection';
import { showToast } from '$lib/store/toasts';
import { ToastType } from '$lib/enums/ToastType';
import { WebSocketCloseCode } from '$lib/enums/WebSocketCloseCode';

type GameUpdateCallback = (update: { game: GameSessionUpdate | null, me: Player | null }) => void;

export class GameSession {
	private lastUpdate: GameUpdate | GameSessionUpdate | null = null;
	private me: Player | null = null;
	private updateCallbacks: GameUpdateCallback[] = [];
	private connection: SocketConnection | null = null;

  constructor(private gameId: string) {}

	public initialize = async () => {
		this.me = await API.getPlayer().catch(() => goto(`/${this.gameId}/join`)) ?? null;
		if (!this.me) return goto(`/${this.gameId}/join`);

		this.connection = await API.getSocketConnection(this.gameId);
    this.connection?.onMessage(this.handleMessage);
    this.connection?.onError(() => this.handleError(new Error("An error occurred with the connection.")));
    this.connection?.onClose(this.handleClose);
    this.connection?.start();
		this.requestUpdate();
	};

	public onUpdate = (callback: GameUpdateCallback) => this.updateCallbacks.push(callback);

  private handleMessage = (message: ServerMessage) => {
    if (isErrorMessage(message)) return this.handleError(new Error(message.payload));
		if (isPlayerKickedMessage(message)) return goto("/");
		if (isGameSessionUpdateMessage(message)) return this.handleUpdate(message);
	};

  private handleClose = (event: CloseEvent) => {
    if (event.code === WebSocketCloseCode.NORMAL) return showToast({
      type: ToastType.Info,
      message: "The connection was closed.",
    });
    else showToast({
      type: ToastType.Warning,
      message: "The connection was closed unexpectedly.",
    });
  };

  private handleError = (error: Error) => {
    showToast({
      type: ToastType.Error,
      message: error.message,
    });
  };

	private handleUpdate = (message: GameSessionUpdateMessage) => {
		this.lastUpdate = message.payload;
		this.me = message.payload?.players.find(({ id }) => this.me?.id === id) ?? this.me;

		localStorage.setItem(`me-${this.gameId}`, JSON.stringify(this.me));
		localStorage.setItem(`lastUpdate-${this.gameId}`, JSON.stringify(this.lastUpdate));

		this.updateCallbacks.forEach((callback) => callback({
			game: message.payload,
			me: this.me,
		}));
	};

	public getCachedUpdate = () => {
		if (this.lastUpdate && this.me) return {
			game: this.lastUpdate,
			me: this.me
		};

		return {
			game: JSON.parse(localStorage.getItem(`lastUpdate-${this.gameId}`) ?? 'null'),
			me: JSON.parse(localStorage.getItem(`me-${this.gameId}`) ?? 'null')
		};
	};

	public kickPlayer = (playerToKick: Player) => {
		if (this.connection) this.connection.send({
			type: ClientMessageType.KickPlayer,
			payload: playerToKick.id,
		});
	};

	public startGame = () => {
		if (this.connection) this.connection.send({
			type: ClientMessageType.StartGame,
			payload: null,
		});
	};

	public submitAnswer = (answer: Answer) => {
		if (this.connection) this.connection.send({
			type: ClientMessageType.SubmitAnswer,
			payload: JSON.stringify(answer),
		});
	}

	private requestUpdate = () => {
		if (this.connection) this.connection.send({
			type: ClientMessageType.UpdateRequest,
			payload: null,
		});
	}
}