package models

import (
	wikipediaModels "api/lib/wikipedia/models"
)

type Pack struct {
	Language            Language                                                                          `json:"language"`
	ApiUrl              string                                                                            `json:"apiUrl"`
	GetLinkSuggestions  func(page wikipediaModels.Page, trackTitle string, artistNames []string) []string `json:"getLinkSuggestions"`
	GetTitleSuggestions func(trackTitle string, artistNames []string) []string                            `json:"getTitleSuggestions"`
	GetReleaseYear      func(response wikipediaModels.Response, artistNames []string) (uint, error)       `json:"getReleaseYear"`
}
