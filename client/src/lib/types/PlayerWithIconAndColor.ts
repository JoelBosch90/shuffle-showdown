import type { Player } from './Player';

export interface PlayerWithIconAndColor extends Player {
  icon: string;
  color: string;
}