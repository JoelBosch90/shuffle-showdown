package spotify

import (
	spotifyModels "api/lib/spotify/models"
	"encoding/json"
	"net/http"
	"strconv"
)

func getCountryParam(params []Param) Param {
	for _, param := range params {
		if param.Name == "market" {
			return param
		}
	}

	return Param{}
}

func AddAdditionalTracks(playListInfo *spotifyModels.Playlist, path string, playListHeaders []Header, playListParams []Param) ([]spotifyModels.Item, error) {
	allAdditionalTrackItems := []spotifyModels.Item{}
	tracksParams := []Param{
		getCountryParam(playListParams),
		{Name: "fields", Value: GetSpotifyModelFields(spotifyModels.Tracks{})},
		{Name: "limit", Value: strconv.Itoa(playListInfo.Tracks.Limit)},
	}
	offset := playListInfo.Tracks.Offset + playListInfo.Tracks.Limit

	// Get all the other Track Items
	for playListInfo.Tracks.Total > offset {
		requestParams := append(tracksParams, Param{Name: "offset", Value: strconv.Itoa(offset)})

		additionalTracksResponse, additionalTracksRequestError := ApiRequest(http.MethodGet, path+"/tracks", playListHeaders, requestParams)
		if additionalTracksRequestError != nil {
			return []spotifyModels.Item{}, additionalTracksRequestError
		}

		// Add the new tracks to the list
		var additionalTracks spotifyModels.Tracks
		additionalTracksDecoder := json.NewDecoder(additionalTracksResponse.Body)
		additionalTracksDecodeError := additionalTracksDecoder.Decode(&additionalTracks)
		if additionalTracksDecodeError != nil {
			return []spotifyModels.Item{}, additionalTracksDecodeError
		}

		allAdditionalTrackItems = append(allAdditionalTrackItems, additionalTracks.Items...)

		offset += playListInfo.Tracks.Limit
	}

	return allAdditionalTrackItems, nil
}
