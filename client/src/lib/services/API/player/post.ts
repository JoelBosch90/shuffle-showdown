import { type Player } from "$lib/types/Player";

export default async (name: string) : Promise<Player> => {
  const response = await fetch('/api/v1/player', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        name
    })
  });
  
  if (!response.ok) throw Error("Failed to create player.");

  const { player } = await response.json();

  return player;
}