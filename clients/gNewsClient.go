package clients

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prateekgupta3991/refresher/entities"
)

func InitGNewsClient() *GClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	gClient := &GClient {
		HttpClient: &http.Client{
			Transport: tr,
		},
	}

    return gClient
}

type GClient struct {
	HttpClient *http.Client
}

type Apis interface {
	GetSources() (*entities.NewsSource, error)
	GetHeadlines() (*entities.TopHeadline, error)
	GetEverything() (*entities.Everything, error)
}

func (c *GClient) GetSources() (*entities.NewsSource, error) {
	req, err := http.NewRequest("GET", "https://newsapi.org/v2/sources?language=en&country=in", nil); if err != nil {
		return nil, err
	} else {
		req.Header.Add("X-Api-Key", "17d9a468d74748d3a39175d524747e95")
		resp, err := c.HttpClient.Do(req); if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		ns := new(entities.NewsSource)
		body, err := ioutil.ReadAll(resp.Body); if err != nil {
			return nil, err
		} else {
			err := json.Unmarshal(body, &ns); if err != nil {
				return nil, err
			} else {
				return ns, nil
			}
		}
	}
}

func (c *GClient) GetHeadlines() (*entities.TopHeadline, error) {
	req, err := http.NewRequest("GET", "https://newsapi.org/v2/top-headlines?country=in&pageSize=5&page=0", nil); if err != nil {
		return nil, err
	} else {
		req.Header.Add("X-Api-Key", "17d9a468d74748d3a39175d524747e95")
		resp, err := c.HttpClient.Do(req); if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		ns := new(entities.TopHeadline)
		body, err := ioutil.ReadAll(resp.Body); if err != nil {
			return nil, err
		} else {
			err := json.Unmarshal(body, &ns); if err != nil {
				return nil, err
			} else {
				return ns, nil
			}
		}
	}
}

func (c *GClient) GetEverything() (*entities.Everything, error) {
	req, err := http.NewRequest("GET", "https://newsapi.org/v2/everything?sources=associated-press,financial-post&q=trump", nil); if err != nil {
		return nil, err
	} else {
		req.Header.Add("X-Api-Key", "17d9a468d74748d3a39175d524747e95")
		resp, err := c.HttpClient.Do(req); if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		ns := new(entities.Everything)
		body, err := ioutil.ReadAll(resp.Body); if err != nil {
			return nil, err
		} else {
			err := json.Unmarshal(body, &ns); if err != nil {
				return nil, err
			} else {
				return ns, nil
			}
		}
	}
}