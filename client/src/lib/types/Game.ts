export interface Game {
  id: string;
  playlistId: number;
  songsToWin: number;
  titleRequired: boolean;
  artistRequired: boolean;
  configured: boolean;
}