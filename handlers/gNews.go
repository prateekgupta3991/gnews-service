package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/refresher/clients"
)

type GNewsService struct {
	Client *clients.GClient
}

func GetNewGNews() *NewsService {
	return &NewsService{
		Service: &GNewsService{
			Client: clients.InitGNewsClient(),
		},
	}
}

func (g *GNewsService) GetSources(c *gin.Context) {
	s, err := g.Client.GetSources(); if err != nil {
		fmt.Errorf("Error encountered : %v", err.Error())
	} else {
		c.JSON(http.StatusOK, s)
	}
}

func (g *GNewsService) GetHeadlines(c *gin.Context) {
	s, err := g.Client.GetHeadlines(); if err != nil {
		fmt.Errorf("Error encountered : %v", err.Error())
	} else {
		c.JSON(http.StatusOK, s)
	}
}

func (g *GNewsService) GetNews(c *gin.Context) {
	s, err := g.Client.GetEverything(); if err != nil {
		fmt.Errorf("Error encountered : %v", err.Error())
	} else {
		c.JSON(http.StatusOK, s)
	}
}