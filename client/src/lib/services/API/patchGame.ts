import { type GameConfig } from '$lib/types/GameConfig';

export default async (settings: GameConfig, playerId: string) => {
  const requestBody = {
    settings,
    playerId
  }

  return fetch(`/api/v1/game/${settings.gameId}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestBody)
  });
}