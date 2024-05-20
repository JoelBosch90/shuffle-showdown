import type { Artist } from './Artist';
export interface Track {
    id?: string;
    name?: string;
    artists?: Artist[];
    releaseYear?: Number;
    previewUrl: string;
}