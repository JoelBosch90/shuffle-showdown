export interface PlaylistLoadingUpdate {
  sentAt: Date;
  finishedLoading: boolean;
  finishedChecking: boolean;
  tracksChecked: number;
  tracksVerified: number;
  tracksLoaded: number;
  tracksTotal: number;
}