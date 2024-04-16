export default async (playListString: string) => {
  const playList = encodeURIComponent(playListString);

  const response = await fetch('/api/v1/game', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ playList })
  });
  
  if (!response.ok) throw Error("Failed to create game.");

  return response.json();
}