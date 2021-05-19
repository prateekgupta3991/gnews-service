package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/prateekgupta3991/refresher/entities"
)

func InitTelegramClient() *TelegramClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	telegramClient := &TelegramClient{
		HttpClient: &http.Client{
			Transport: tr,
		},
	}

	return telegramClient
}

type TelegramClient struct {
	HttpClient *http.Client
}

type Send interface {
	Send() error
}

func prepareBaseUrl(base string, qp map[string][]string) string {
	var url strings.Builder
	url.WriteString(base)
	url.WriteString("?")
	for idx, val := range qp {
		for _, value := range val {
			url.WriteString(fmt.Sprintf("%s=%s", idx, value))
			url.WriteString("&")
		}
	}
	return strings.TrimRight(url.String(), "&")
}

func (c *TelegramClient) Send() error {
	url := prepareUrl("http://newsapi.org/v2/everything", qp)
	if req, err := http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	} else {
		req.Header.Add("X-Api-Key", "17d9a468d74748d3a39175d524747e95")
		resp, err := c.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		ns := new(entities.Everything)
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			err := json.Unmarshal(body, &ns)
			if err != nil {
				return nil, err
			} else {
				return ns, nil
			}
		}
	}
}
