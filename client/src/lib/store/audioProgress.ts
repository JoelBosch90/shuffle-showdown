import { writable } from 'svelte/store';
import type { AudioProgressState } from '$lib/types/AudioProgressState';

const DEFAULT_AUDIO_PROGRESS: AudioProgressState = {
  progress: 0,
  max: 100,
  isPlaying: false
};

type AudioProgressStateStateUpdate = Partial<AudioProgressState>;

export const audioProgressState = writable<AudioProgressState>(DEFAULT_AUDIO_PROGRESS);

export const setAudioProgress = (newState: AudioProgressStateStateUpdate) => {
  audioProgressState.update(previousState => ({ ...previousState, ...newState }));
};
