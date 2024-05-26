import { type GameUpdate } from "$lib/types/GameUpdate";

export default async (playListString: string, playerId?: string) : Promise<GameUpdate> => {
  const playList = encodeURIComponent(playListString);

  // By lack of better methods without doing any external requests, we try to 
  // estimate the country code based off of the navigator.lanuage.
  const countryCode = navigator.language.slice(-2).toUpperCase();

  const response = await fetch('/api/v1/game', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      playList,
      countryCode,
      playerId
    })
  });
  
  if (!response.ok) throw Error("Failed to create game.");

  const { game } = await response.json();

  return game;
}