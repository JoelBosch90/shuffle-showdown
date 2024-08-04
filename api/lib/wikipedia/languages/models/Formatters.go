package models

import (
	"api/lib/wikipedia/models"
)

type LinkBasedFormatter func(page models.Page, trackTitle string, artistNames []string) string
type TitleBasedFormatter func(trackTitle string, artistNames []string) string

type RedirectFormatters struct {
	LinkBasedFormatters  []LinkBasedFormatter
	TitleBasedFormatters []TitleBasedFormatter
}
