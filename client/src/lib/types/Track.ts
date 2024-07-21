import type { Artist } from './Artist';
import type { Album } from './Album';
export interface Track {
    id?: string;
    name?: string;
    artists?: Artist[];
    album?: Album;
    previewUrl: string;
}