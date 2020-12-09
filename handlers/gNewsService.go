package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/refresher/clients"
	"github.com/prateekgupta3991/refresher/entities"
	"github.com/prateekgupta3991/refresher/repository"
	"github.com/prateekgupta3991/refresher/validations"
)

type GNewsService struct {
	Client *clients.GClient
	DbClient *repository.DbSession
}

func GetNewGNews() *NewsService {
	return &NewsService{
		Service: &GNewsService{
			Client: clients.InitGNewsClient(),
			DbClient: repository.GetNewDbSession(),
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
		var newsBySourceList []entities.NewsBySource
		for _, val := range s.ArticleList.Articles {
			// if pTime, err := time.Parse("2006-01-02T15:04:05Z", val.PublishedAt);
			newsBySourceList = append(newsBySourceList, entities.NewsBySource{
				SourceId: val.Source.Id,
				TitleHash: "srg5egebr5156",
				NewsAuthor: val.Author,
				NewsContent: val.Content,
				NewsDescription: val.Description,
				SourceName: val.Source.Name,
				NewsPublishedAt: val.PublishedAt,
			})
		}
		g.DbClient.InsertTopNews(newsBySourceList[0])
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
		if dbResp, err := g.DbClient.GetTopNews("associated-press", 3); err == nil {
			c.JSON(http.StatusOK, dbResp)
		}
		fmt.Println(s)
		// c.JSON(http.StatusOK, s)
	}
}