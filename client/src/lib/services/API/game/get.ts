import { type Game } from "$lib/types/Game";

export default async (gameId: string) : Promise<Game> => {
  const response = await fetch(`/api/v1/game/${gameId}`);

  if (!response.ok) throw Error("Failed to get game.");

  const { game } = await response.json();

  return game;
}