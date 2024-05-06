export default async (playerName: string, playerId: string) => {
  const requestBody = {
    playerName
  }

  return fetch(`/api/v1/player/${playerId}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestBody)
  });
}