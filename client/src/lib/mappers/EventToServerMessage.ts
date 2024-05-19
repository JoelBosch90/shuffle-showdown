import { type ServerMessage } from '$lib/types/ServerMessage';

export const eventToServerMessage = (event: MessageEvent) : ServerMessage => {
    const { type, payload: payloadJson } = JSON.parse(event.data);

    // The payload is always separately JSON-encoded string, so we need to parse it.
    const payload = JSON.parse(payloadJson);

    return { type, payload };
}