export default async (playListString: string) => {
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
      countryCode
    })
  });
  
  if (!response.ok) throw Error("Failed to create game.");

  return response.json();
}