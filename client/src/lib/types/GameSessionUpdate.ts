import type { Player } from '$lib/types/Player';
import type { Round } from '$lib/types/Round';
import type { GameUpdate } from '$lib/types/GameUpdate';

export interface GameSessionUpdate extends GameUpdate {
    sentAt: Date;
    players: Player[];
    rounds: Round[];
}