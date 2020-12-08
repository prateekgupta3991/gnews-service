package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/refresher/clients"
	"github.com/prateekgupta3991/refresher/validations"
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
	qp := c.Request.URL.Query()
	ok, err := validations.RequestQParams(qp, "source"); if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
	}
	s, err := g.Client.GetSources(qp); if err != nil {
		fmt.Printf("Error encountered : %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": err})
	} else {
		c.JSON(http.StatusOK, s)
	}
}

func (g *GNewsService) GetHeadlines(c *gin.Context) {
	qp := c.Request.URL.Query()
	ok, err := validations.RequestQParams(qp, "headlines"); if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
	}
	s, err := g.Client.GetHeadlines(qp); if err != nil {
		fmt.Printf("Error encountered : %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": err})
	} else {
		c.JSON(http.StatusOK, s)
	}
}

func (g *GNewsService) GetNews(c *gin.Context) {
	qp := c.Request.URL.Query()
	ok, err := validations.RequestQParams(qp, "all"); if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
	}
	s, err := g.Client.GetEverything(qp); if err != nil {
		fmt.Printf("Error encountered : %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": err})
	} else {
		c.JSON(http.StatusOK, s)
	}
}