import { type GameUpdate } from "$lib/types/GameUpdate";

export default async (gameId: string) : Promise<GameUpdate> => {
  const response = await fetch(`/api/v1/game/${gameId}`);

  if (!response.ok) throw Error("Failed to get game.");

  const { game } = await response.json();

  return game;
}