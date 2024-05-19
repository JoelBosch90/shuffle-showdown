import { goto } from '$app/navigation';
import { API } from '$lib/services/API';
import { type Player } from '$lib/types/Player';
import { type GameUpdate } from '$lib/types/GameUpdate';
import { isPlayerKickedMessage, isGameUpdateMessage, type ServerMessage, type GameUpdateMessage } from '$lib/types/ServerMessage';
import { ClientMessageType } from '$lib/enums/ClientMessageType';

type GameUpdateCallback = (update: { game: GameUpdate | null, me: Player | null }) => void;

export class GameSession {
	private lastUpdate: GameUpdate | null = null;
	private me: Player | null = null;
	private updateCallbacks: GameUpdateCallback[] = [];

    constructor(private gameId: string) {}

	public initialize = async () => {		
		this.lastUpdate = await API.getGame(this.gameId).catch(() => goto('/game')) ?? null;

		if (!this.lastUpdate) goto(`/game/${this.gameId}/configure`);

		this.me = await API.getPlayer().catch(() => goto(`/game/${this.gameId}/join`)) ?? null;

		if (!this.me) {
			return goto(`/game/${this.gameId}/join`);
		}

        API.onSocketMessage(this.handleMessage);
        API.startSocketConnection(this.gameId);
	}

	public onUpdate = (callback: GameUpdateCallback) => this.updateCallbacks.push(callback);
    
    private handleMessage = (message: ServerMessage) => {
		if (isPlayerKickedMessage(message)) return goto('/game');
		if (isGameUpdateMessage(message)) return this.handleUpdate(message);
	};

	private handleUpdate = (message: GameUpdateMessage) => {
		this.lastUpdate = message.payload;
		this.me = message.payload?.players.find(({ id }) => this.me?.id === id) ?? this.me;

		this.updateCallbacks.forEach((callback) => callback({
			game: message.payload,
			me: this.me,
		}));
	}

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
}