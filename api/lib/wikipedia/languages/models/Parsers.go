package models

import (
	"api/lib/wikipedia/models"
)

type ReleaseYearParser func(category models.Category) (uint, error)
type ConfirmArtistParser func(category models.Category, artistName string) bool
type RecognizeRedirectParser func(category models.Category) bool

type CategoryParsers struct {
	ReleaseYearParsers       []ReleaseYearParser
	ConfirmArtistParsers     []ConfirmArtistParser
	RecognizeRedirectParsers []RecognizeRedirectParser
}
