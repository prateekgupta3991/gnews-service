package entities

import(

)

type NewsResponse struct {
	Status string `json:"status"`
	TotalResults int `json:"totalResults"`
}

type Source struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Url string `json:"url"`
	Category string `json:"category"`
	Language string `json:"language"`
	Country string `json:"country"`
}

type SourceList struct {
	Sources []Source `json:"sources"`
}

type ArticleSource struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	Source ArticleSource `json:"source"`
	Autho string `json:"author"`
	Title string `json:"title"`
	Description string `json:"description"`
	Url string `json:"url"`
	UrlToImage string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content string `json:"content"`
}

type ArticleList struct {
	Articles []Article `json:"articles"`
}

type TopHeadline struct {
	NewsResponse
	ArticleList
}

type Everything struct {
	NewsResponse
	ArticleList
}

type NewsSource struct {
	NewsResponse
	SourceList 
}