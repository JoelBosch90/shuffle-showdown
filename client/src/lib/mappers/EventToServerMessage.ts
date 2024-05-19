import { type ServerMessage } from '$lib/types/ServerMessage';

export const eventToServerMessage = (event: MessageEvent) : ServerMessage => {
    console.log('eventToServerMessage', event);
    // Parse the event data.
    const { type, payload: payloadJson } = JSON.parse(event.data);

    // Parse the payload JSON.
    const payload = JSON.parse(payloadJson);

    return { type, payload };
}