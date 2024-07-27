import type { Artist } from './Artist';
export interface Track {
  id?: string;
  name?: string;
  artists?: Artist[];
  releaseYear?: number;
  releaseMonth?: number;
  releaseDay?: number;
  previewUrl: string;
}