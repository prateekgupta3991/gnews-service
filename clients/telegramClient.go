package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/prateekgupta3991/refresher/entities"
)

func InitTelegramClient(hClient *http.Client) *TelegramClient {
	telegramClient := &TelegramClient{
		HttpClient: hClient,
	}

	return telegramClient
}

type TelegramClient struct {
	HttpClient *http.Client
}

type Send interface {
	Send() (*entities.Webhook, error)
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

func (c *TelegramClient) Send(qp map[string][]string) (*entities.Webhook, error) {
	// TODO take bot token from conf
	url := prepareUrl("https://api.telegram.org/bot6945346087:AAFamOovT__-Sw20my4y5qrMLm5iUiABJAk/sendMessage", qp)
	if req, err := http.NewRequest("POST", url, nil); err != nil {
		fmt.Println(url)
		return nil, err
	} else {
		resp, err := c.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		wh := new(entities.Webhook)
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
		    fmt.Printf("Error while sending news for chatId - %s  error - %v\n", qp["chat_id"], err.Error())
			return nil, err
		} else {
			// fix the unmarshalling here
			if err := json.Unmarshal(body, &wh); err != nil || !wh.Ok {
			    fmt.Printf("Marshalling Error while sending news for chatId - %s  error - %v\n", qp["chat_id"], err.Error())
				return nil, err
			} else {
    			fmt.Printf("Successfully sent news for chatId - %s\n", qp["chat_id"])
				return wh, nil
			}
		}
	}
}
