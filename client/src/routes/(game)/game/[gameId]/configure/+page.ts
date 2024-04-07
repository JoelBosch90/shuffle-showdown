import { error } from '@sveltejs/kit';

export function load({ params }) {
  if (!params.gameId) return error(400, 'No game ID provided');

  return { gameId: params.gameId };
}