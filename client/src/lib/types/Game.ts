import { type Player } from './Player';
import { type Playlist } from './Playlist';

export interface Game {
  id: string;
  playlistId: number;
  songsToWin: number;
  titleRequired: boolean;
  artistRequired: boolean;
  configured: boolean;
  owner: Player;
  playlist: Playlist;
  players: Player[];
}