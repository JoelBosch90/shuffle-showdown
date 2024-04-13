import { type GameConfig } from '../types/GameConfig';

export interface Game {
  id: string;
  playlistId: number;
  songsToWin: number;
  titleRequired: boolean;
  artistRequired: boolean;
  configured: boolean;
}

export class API {

  public static postGame = async (playlistId: number) => {
    const response = await fetch('/api/v1/game', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ playlist: playlistId })
		});
    
		if (!response.ok) throw Error("Failed to create game.");

    return response.json();
  }

  public static getGame = async (gameId: string) => {
    const response = await fetch(`/api/v1/game/${gameId}`);

    if (!response.ok) throw Error("Failed to get game.");

    return response.json() as Promise<Game>;
  }

  public static patchGame = async (settings: GameConfig) => {
    return fetch(`/api/v1/game/${settings.gameId}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(settings)
    });
  }
}