import { writable } from 'svelte/store';

export interface GameConfig {
  gameId: string;
  songsToWin: number;
  titleRequired: boolean;
  artistRequired: boolean;
  releaseYearRequired: boolean;
}

export const gameConfigStore = writable<GameConfig>();