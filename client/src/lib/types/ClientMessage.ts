import { type ClientMessageType } from '$lib/enums/ClientMessageType';

export type ClientMessage = {
    type: ClientMessageType,
    payload: string | null,
    playerId: string | null,
};