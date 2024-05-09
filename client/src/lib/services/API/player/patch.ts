import { type Player } from "$lib/types/Player";

export default async (name: string, playerId: string) : Promise<Player> => {
  const requestBody = {
    name
  }

  const response = await fetch(`/api/v1/player/${playerId}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestBody)
  });
  
  if (!response.ok) throw Error("Failed to set player name.");

  const { player } = await response.json();

  return player;
}