import { writable } from 'svelte/store';
import type { ProgressState } from '$lib/types/ProgressState';

export const DEFAULT_PROGRESS: ProgressState = {
  progress: 0,
  max: 100
};

type ProgressStateUpdate = Omit<ProgressState, 'progress'> | Omit<ProgressState, 'max'>;

export const progressState = writable<ProgressState>(DEFAULT_PROGRESS);

export const setProgress = (newState: ProgressStateUpdate) => {
  progressState.update(previousState => ({ ...previousState, ...newState }));
};
