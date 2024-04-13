export default async (playlistId: number) => {
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