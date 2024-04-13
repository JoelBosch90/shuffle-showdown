import { type Game } from "$lib/types/Game";

export default async (gameId: string) => {
  const response = await fetch(`/api/v1/game/${gameId}`);

  if (!response.ok) throw Error("Failed to get game.");

  return response.json() as Promise<Game>;
}