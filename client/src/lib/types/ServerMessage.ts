import { type Player } from './Player';
import { ServerMessageType } from '$lib/enums/ServerMessageType';

export interface ServerMessage {
    type: ServerMessageType;
    content: unknown;
}

export interface PlayersUpdateMessage extends ServerMessage {
    type: ServerMessageType.PlayersUpdate;
    content: Player[];
}

export const isPlayersUpdateMessage = (message: ServerMessage) : message is PlayersUpdateMessage => {
    return message.type === ServerMessageType.PlayersUpdate;
}

export interface ErrorMessage extends ServerMessage {
    type: ServerMessageType.Error;
    content: string;
}

export const isErrorMessage = (message: ServerMessage) : message is ErrorMessage => {
    return message.type === ServerMessageType.Error;
}