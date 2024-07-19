import { type MostPlayedPlaylist } from "$lib/types/MostPlayedPlaylist";

export default async (): Promise<MostPlayedPlaylist[]> => {
  const response = await fetch(`/api/v1/playlists/most-played`);

  if (!response.ok) throw Error("Failed to get playlists.");

  const { playlists } = await response.json();

  return playlists;
}