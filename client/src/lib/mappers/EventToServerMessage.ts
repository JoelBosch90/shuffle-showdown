import { type ServerMessage } from '$lib/types/ServerMessage';

export const eventToServerMessage = (event: MessageEvent) : ServerMessage => {
    // Parse the event data.
    const { type, content: contentJson } = JSON.parse(event.data);

    // Parse the content JSON.
    const content = JSON.parse(contentJson);

    return { type, content };
}