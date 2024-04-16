package lib

import (
	"net/url"
	"strings"
)

func ExtractSpotifyPlayListId(encodedUrl string) string {
	decodedUrl, decodingError := url.QueryUnescape(encodedUrl)
	if decodingError != nil {
		return encodedUrl
	}

	// First, we look for the path in the URL that contains the playlist ID.
	urlParts := strings.Split(decodedUrl, "/playlist/")
	if len(urlParts) < 2 {
		return decodedUrl
	}

	// If there are any query parameters, we remove them.
	return strings.Split(urlParts[1], "?")[0]
}
