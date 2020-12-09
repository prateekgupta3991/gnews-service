package entities

import "time"

type NewsBySource struct {
	SourceId string `json:"sid"`
	CreatedAt time.Time `json:"created_at"`
	TitleHash string `json:"title_hash"`
	NewsAuthor string `json:"nauthor"`
	NewsContent string `json:"ncontent"`
	NewsDescription string `json:"ndesc"`
	NewsPublishedAt string `json:"npublished_at"`
	NewsTitle string `json:"ntitle"`
	NewsUrl string `json:"nurl"`
	NewsUrlToImage string `json:"nurl_to_image"`
	SourceCategory string `json:"scategory"`
	SourceDescription string `json:"sdesc"`
	SourceCountry string `json:"scountry"`
	SourceLanguage string `json:"slang"`
	SourceName string `json:"sname"`
	SourceUrl string `json:"surl"`
}