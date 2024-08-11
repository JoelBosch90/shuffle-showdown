export interface PlaylistLoadingUpdate {
  sentAt: Date;
  liveUpdate: boolean;
  finishedLoading: boolean;
  finishedChecking: boolean;
  tracksChecked: number;
  tracksVerified: number;
  tracksLoaded: number;
  tracksTotal: number;
}