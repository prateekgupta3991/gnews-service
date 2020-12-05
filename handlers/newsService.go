package handlers

import "github.com/gin-gonic/gin"

type NewsService struct {
	Service interface{}
}

type NewsData interface {
	GetSources(c *gin.Context) []interface{}
	GetHeadlines(c *gin.Context) []interface{}
	GetNews(c *gin.Context) []interface{}
}