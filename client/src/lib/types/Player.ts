import type { WonTrack } from './WonTrack';

export interface Player {
  id: string;
  name?: string;
  isConnected?: boolean;
  isOwner?: boolean;
  wonTracks?: WonTrack[];
}