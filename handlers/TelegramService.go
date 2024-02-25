package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/prateekgupta3991/refresher/clients"
	"github.com/prateekgupta3991/refresher/entities"
	"github.com/prateekgupta3991/refresher/repository"
)

type Telegram struct {
	TelegramClient   *clients.TelegramClient
	TelegramDbClient *repository.UserDbSession
	GNewsService     *GNewsService
}

func NewTelegram(hClient *http.Client, cassession *gocql.Session, newsServ *GNewsService) *Communication {
	return &Communication{
		IMService: &Telegram{
			TelegramClient:   clients.InitTelegramClient(hClient),
			TelegramDbClient: repository.NewUserDbSession(cassession),
			GNewsService:     newsServ,
		},
	}
}

func (t *Telegram) PushedUpdates(c *gin.Context) {
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		fmt.Printf("Error encountered : %v", err.Error())
		c.JSON(http.StatusBadRequest, "Bad request")
	} else {
		webhookObj := new(entities.Result)
		err := json.Unmarshal(body, &webhookObj)
		fmt.Println(webhookObj)
		if err != nil {
			fmt.Printf("Could not process the webhook. Error encountered : %v\n", err.Error())
			c.JSON(http.StatusBadRequest, "BadRequest")
		} else {
			// fmt.Printf("%+v", webhookObj)
			if webhookObj.Query.Id != "" {
				fmt.Printf("Inline reply received for subscriber with Id : %d and Username : %s and ChatId : %d\n", webhookObj.Query.From.Id, webhookObj.Query.Msg.Chat.UserName, webhookObj.Query.Msg.Chat.Id)
				msgTxt, err := t.fetchTopNewsRequest("5", []string{webhookObj.Query.Data})
				if err != nil {
					fmt.Printf("Exception due to - %v", err.Error())
					c.JSON(http.StatusInternalServerError, "InternalServiceError")
				}
				cid := strconv.Itoa(int(webhookObj.Query.Msg.Chat.Id))
				for _, txt := range msgTxt {
					t.CallTelegramSendApi(cid, txt, entities.ButtonsInMessage{[][]entities.Button{}})
				}
				c.JSON(http.StatusOK, "OK")
				return
			} else if subscriber, err := t.TelegramDbClient.GetUserByTgDetils(int64(webhookObj.Msg.From.Id), webhookObj.Msg.From.UserName); err != nil || subscriber.ID == 0 {
				fmt.Printf("New subscriber with Id : %d and Username : %s and ChatId : %d\n", webhookObj.Msg.From.Id, webhookObj.Msg.From.UserName, webhookObj.Msg.Chat.Id)
				m := entities.UserDetails{
					ID:         int64(webhookObj.Msg.From.Id),
					Name:       webhookObj.Msg.From.FirstName,
					TelegramId: webhookObj.Msg.From.UserName,
					ChatId:     int64(webhookObj.Msg.Chat.Id),
				}
				if err := t.TelegramDbClient.InsertUser(m); err != nil {
					fmt.Printf("Failure while persisting subscriber with Id : %d and Username : %s - %s\n", subscriber.ID, subscriber.TelegramId, err.Error())
					c.JSON(http.StatusInternalServerError, "InternalServiceError")
				}
			} else {
				fmt.Printf("Subscriber found with Id : %d and Username : %s\n", subscriber.ID, subscriber.TelegramId)
			}

			cid := strconv.Itoa(int(webhookObj.Msg.Chat.Id))
			// t.CallTelegramSendApi(cid, "Your news feed is updated", entities.ButtonsInMessage{})
			reqId, _ := c.Get("uuid")
			if sources, ok := t.fetchNewsSources(reqId, []string{"en"}, []string{"in"}); !ok {
				fmt.Printf("Unable to converse about source preference during requestId : %s", reqId)
				c.JSON(http.StatusInternalServerError, "InternalServiceError")
			} else {
				buttonsList := make([][]entities.Button, 1)
				srcAsButtons := entities.ButtonsInMessage{buttonsList}
				srcListStr := "Select the news source"
				for _, val := range sources.SourceList.Sources {
					buttonsList[0] = append(buttonsList[0], entities.Button{val.Name, "", val.Id})
				}
				t.CallTelegramSendApi(cid, srcListStr, srcAsButtons)
			}
			c.JSON(http.StatusOK, "OK")
		}
	}
}

func (t *Telegram) Notify(c *gin.Context) {
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		fmt.Printf("Error encountered : %v\n", err.Error())
		c.JSON(http.StatusOK, err.Error())
		return
	} else {
		reply := new(entities.TelegramReplyMsg)
		err := json.Unmarshal(body, &reply)
		if err != nil {
			fmt.Printf("Could not unmarshal the body. Error encountered : %v\n", err.Error())
			c.JSON(http.StatusOK, err.Error())
			return
		} else {
			fmt.Printf("inputs %v\n", reply)
			if reply.ChatId == 0 && strings.EqualFold(reply.UserName, "") {
				// send to every subscriber
				if usrDetails, err := t.TelegramDbClient.GetAllUser(); err != nil || len(usrDetails) == 0 {
					fmt.Printf("Could not find user by username. Error encountered : %v\n", err.Error())
					c.JSON(http.StatusOK, err.Error())
					return
				} else {
					msgTxt, err := t.fetchTopNewsRequest("5", []string{"google-news-in", "the-times-of-india", "the-hindu"})
					if err != nil {
						fmt.Printf("Exception due to - %v\n", err.Error())
					}
					for _, usr := range usrDetails {
						if usr.ChatId != 0 {
							cid := strconv.Itoa(int(usr.ChatId))
							for _, txt := range msgTxt {
								t.CallTelegramSendApi(cid, txt, entities.ButtonsInMessage{[][]entities.Button{}})
							}
						}
					}
				}
			} else if reply.ChatId != 0 && !strings.EqualFold(reply.UserName, "") {
				// dont fetch usd. directly use the chatId
				cid := strconv.Itoa(int(reply.ChatId))
				t.CallTelegramSendApi(cid, reply.Text, entities.ButtonsInMessage{})
			} else if strings.EqualFold(reply.UserName, "") {
				// dont fetch usd. directly use the chatId
				cid := strconv.Itoa(int(reply.ChatId))
				t.CallTelegramSendApi(cid, reply.Text, entities.ButtonsInMessage{})
			} else {
				// fetch usd by username
				if usrDet, err := t.TelegramDbClient.GetUserByTgUn(reply.UserName); err != nil || usrDet.ID == 0 {
					fmt.Printf("Could not find user by username. Error encountered : %v\n", err.Error())
					c.JSON(http.StatusOK, err.Error())
					return
				} else {
					cid := strconv.Itoa(int(usrDet.ChatId))
					t.CallTelegramSendApi(cid, reply.Text, entities.ButtonsInMessage{[][]entities.Button{}})
				}
			}
			c.JSON(http.StatusOK, "OK")
		}
	}
}

func (t *Telegram) CallTelegramSendApi(chatId, text string, buttons entities.ButtonsInMessage) error {
	var qp map[string]interface{} = make(map[string]interface{})
	qp["chat_id"] = []string{chatId}
	qp["text"] = []string{text}
	if replyMarkupJSON, err := json.Marshal(buttons); err != nil {
		fmt.Println("Error while marshalling buttons json")
	} else {
		fmt.Println(string(replyMarkupJSON))
		qp["reply_markup"] = string(replyMarkupJSON)
	}
	fmt.Printf("Sending news for chatId - %s\n %s\n", chatId, text)
	if _, err := t.TelegramClient.Send(qp); err != nil {
		fmt.Printf("Error while sending msg %s\n", err)
		return err
	}
	return nil
}

func (t *Telegram) fetchTopNewsRequest(newsCount string, sources []string) ([]string, error) {
	qp := make(map[string][]string)
	qp["top"] = []string{newsCount}
	qp["sources"] = sources
	nc, _ := strconv.Atoi(newsCount)
	if news, err := t.GNewsService.GetTopNewsBySourceFromDb(qp, nc); err != nil {
		fmt.Printf("Error while fetching news. Error encountered : %v\n", err.Error())
		return nil, err
	} else {
		msgTxt := make([]string, nc+1)
		i := 1
		msgTxt[0] = fmt.Sprintf("Tada...Here is the top %s news for you", newsCount)
		for _, val := range news {
			msgTxt[i] = val.NewsDescription + " - " + val.NewsUrl
			fmt.Println(msgTxt[i])
			i++
		}
		return msgTxt, nil
	}
}

func (t *Telegram) fetchNewsSources(requestId interface{}, language, countries []string) (*entities.NewsSource, bool) {
	qp := make(map[string][]string)
	qp["language"] = language
	qp["country"] = countries
	if sources, ok := t.GNewsService.GetNewsSources(qp); !ok {
		fmt.Printf("Unable to converse about source preference during requestId : %s\n", requestId)
		return nil, false
	} else {
		fmt.Printf("Got following sources %v for requestId : %s\n", sources, requestId)
		return sources, true
	}
}
