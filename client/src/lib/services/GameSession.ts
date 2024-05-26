import { goto } from '$app/navigation';
import { API } from '$lib/services/API';
import { type Player } from '$lib/types/Player';
import { type GameUpdate } from '$lib/types/GameUpdate';
import { type GameSessionUpdate } from '$lib/types/GameSessionUpdate';
import { type Answer } from '$lib/types/Answer';
import { isPlayerKickedMessage, isGameSessionUpdateMessage, type ServerMessage, type GameSessionUpdateMessage } from '$lib/types/ServerMessage';
import { ClientMessageType } from '$lib/enums/ClientMessageType';

type GameUpdateCallback = (update: { game: GameSessionUpdate | null, me: Player | null }) => void;

export class GameSession {
	private lastUpdate: GameUpdate | GameSessionUpdate | null = null;
	private me: Player | null = null;
	private updateCallbacks: GameUpdateCallback[] = [];

    constructor(private gameId: string) {}

	public initialize = async () => {
		this.me = await API.getPlayer().catch(() => goto(`/game/${this.gameId}/join`)) ?? null;
		if (!this.me) return goto(`/game/${this.gameId}/join`);

        API.onSocketMessage(this.handleMessage);
        API.startSocketConnection(this.gameId);
		this.requestUpdate();
	};

	public onUpdate = (callback: GameUpdateCallback) => this.updateCallbacks.push(callback);
    
    private handleMessage = (message: ServerMessage) => {
		if (isPlayerKickedMessage(message)) return goto('/game');
		if (isGameSessionUpdateMessage(message)) return this.handleUpdate(message);
	};

	private handleUpdate = (message: GameSessionUpdateMessage) => {
		this.lastUpdate = message.payload;
		this.me = message.payload?.players.find(({ id }) => this.me?.id === id) ?? this.me;

		localStorage.setItem('me', JSON.stringify(this.me));
		localStorage.setItem('lastUpdate', JSON.stringify(this.lastUpdate));

		this.updateCallbacks.forEach((callback) => callback({
			game: message.payload,
			me: this.me,
		}));
	}

	public getCachedUpdate = () => {
		if (this.lastUpdate && this.me) return {
			game: this.lastUpdate,
			me: this.me
		};

		return {
			game: JSON.parse(localStorage.getItem('lastUpdate') ?? 'null'),
			me: JSON.parse(localStorage.getItem('me') ?? 'null')
		};
	};

	public kickPlayer = (playerToKick: Player) => {
		API.sendSocketMessage({
			type: ClientMessageType.KickPlayer,
			payload: playerToKick.id,
		});
	};

	public startGame = () => {
		API.sendSocketMessage({
			type: ClientMessageType.StartGame,
			payload: null,
		});
	};

	public submitAnswer = (answer: Answer) => {
		API.sendSocketMessage({
			type: ClientMessageType.SubmitAnswer,
			payload: JSON.stringify(answer),
		});
	}

	private requestUpdate = () => {
		API.sendSocketMessage({
			type: ClientMessageType.UpdateRequest,
			payload: null,
		});
	}
}