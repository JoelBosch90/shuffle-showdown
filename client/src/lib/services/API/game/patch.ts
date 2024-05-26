import { type GameConfig } from '$lib/types/GameConfig';
import { type GameUpdate } from "$lib/types/GameUpdate";

export default async (settings: GameConfig, playerId: string) : Promise<GameUpdate> => {
  const requestBody = {
    settings,
    playerId
  }

  const response = await fetch(`/api/v1/game/${settings.gameId}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestBody)
  });
  
  if (!response.ok) throw Error("Failed to configure game.");

  const { game } = await response.json();

  return game;
}