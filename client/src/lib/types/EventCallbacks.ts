import { type ServerMessage } from '$lib/types/ServerMessage';

export type GenericEventCallback = (event: Event) => void;
export type CloseEventCallback = (event: CloseEvent) => void;
export type MessageEventCallback = (event: ServerMessage) => void;