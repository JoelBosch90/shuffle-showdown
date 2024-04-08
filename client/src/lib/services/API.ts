import { type GameConfig } from '../types/GameConfig';

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
    return fetch(`/api/v1/game/${gameId}`);
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