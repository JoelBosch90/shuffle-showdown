import { type Player } from './Player';
import { type Playlist } from './Playlist';

export interface GameUpdate {
  id: string;
  createdAt: Date;
  updatedAt: Date;
  playlist: Playlist;
  songsToWin: number;
  owner: Player;
  hasStarted: boolean;
  hasFinished: boolean;
}
