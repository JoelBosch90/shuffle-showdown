import { type Player } from './Player';
import { ServerMessageType } from '$lib/enums/ServerMessageType';

export interface ServerMessage {
    type: ServerMessageType;
    payload: unknown;
}

export interface PlayersUpdateMessage extends ServerMessage {
    type: ServerMessageType.PlayersUpdate;
    payload: Player[];
}

export const isPlayersUpdateMessage = (message: ServerMessage) : message is PlayersUpdateMessage => {
    return message.type === ServerMessageType.PlayersUpdate;
}

export interface ErrorMessage extends ServerMessage {
    type: ServerMessageType.Error;
    payload: string;
}

export const isErrorMessage = (message: ServerMessage) : message is ErrorMessage => {
    return message.type === ServerMessageType.Error;
}

export interface PlayerKickedMessage extends ServerMessage {
    type: ServerMessageType.PlayerKicked;
    payload: string;
}

export const isPlayerKickedMessage = (message: ServerMessage) : message is PlayerKickedMessage => {
    return message.type === ServerMessageType.PlayerKicked;
}