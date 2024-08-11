import { type GameSessionUpdate } from './GameSessionUpdate';
import { type AnswerSelectionUpdate } from './AnswerSelectionUpdate';
import { type PlaylistLoadingUpdate } from './PlaylistLoadingUpdate';
import { ServerMessageType } from '$lib/enums/ServerMessageType';

export interface ServerMessage {
  gameId: string;
  type: ServerMessageType;
  payload: unknown;
}

export interface ErrorMessage extends ServerMessage {
  type: ServerMessageType.Error;
  payload: {
    message: string;
  };
}

export const isErrorMessage = (message: ServerMessage): message is ErrorMessage => {
  return message.type === ServerMessageType.Error;
}

export interface PlayerKickedMessage extends ServerMessage {
  type: ServerMessageType.PlayerKicked;
  payload: {
    playerId: string;
  };
}

export const isPlayerKickedMessage = (message: ServerMessage): message is PlayerKickedMessage => {
  return message.type === ServerMessageType.PlayerKicked;
}

export interface GameSessionUpdateMessage extends ServerMessage {
  type: ServerMessageType.GameSessionUpdate;
  payload: GameSessionUpdate;
}

export const isGameSessionUpdateMessage = (message: ServerMessage): message is GameSessionUpdateMessage => {
  return message.type === ServerMessageType.GameSessionUpdate;
}

export interface AnswerSelectionUpdateMessage extends ServerMessage {
  type: ServerMessageType.AnswerSelectionUpdate;
  payload: AnswerSelectionUpdate;
}

export const isAnswerSelectionUpdateMessage = (message: ServerMessage): message is AnswerSelectionUpdateMessage => {
  return message.type === ServerMessageType.AnswerSelectionUpdate;
}

export interface PlaylistLoadingUpdateMessage extends ServerMessage {
  type: ServerMessageType.PlaylistLoadingUpdate;
  payload: PlaylistLoadingUpdate;
}

export const isPlaylistLoadingUpdateMessage = (message: ServerMessage): message is PlaylistLoadingUpdateMessage => {
  return message.type === ServerMessageType.PlaylistLoadingUpdate;
}