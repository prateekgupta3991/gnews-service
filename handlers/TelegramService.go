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
		if err != nil {
			fmt.Printf("Could not process the webhook. Error encountered : %v\n", err.Error())
		} else {
			if subscriber, err := t.TelegramDbClient.GetUserByTgDetils(int(webhookObj.Msg.From.Id), webhookObj.Msg.From.UserName); err != nil || subscriber.ID == 0 {
				fmt.Printf("New subscriber with Id : %d and Username : %s\n", webhookObj.Msg.From.Id, webhookObj.Msg.From.UserName)
				m := entities.UserDetails{
					ID:         int64(webhookObj.Msg.From.Id),
					Name:       webhookObj.Msg.From.FirstName,
					TelegramId: webhookObj.Msg.From.UserName,
					ChatId:     int32(webhookObj.Msg.Chat.Id),
				}
				if err := t.TelegramDbClient.InsertUser(m); err != nil {
					fmt.Printf("Failure while persisting subscriber with Id : %d and Username : %s - %s\n", subscriber.ID, subscriber.TelegramId, err.Error())
				}
			} else {
				fmt.Printf("Subscriber found with Id : %d and Username : %s\n", subscriber.ID, subscriber.TelegramId)
			}
			cid := strconv.Itoa(int(webhookObj.Msg.Chat.Id))
			t.CallTelegramSendApi(cid, "Your news feed is updated")

			qp := make(map[string][]string)
			qp["language"] = []string{"en"}
			qp["country"] = []string{"in"}
			if sources, ok := t.GNewsService.GetNewsSources(qp); !ok {
				rId, _ := c.Get("uuid")
				fmt.Printf("Unable to converse about source preference during requestId : %s", rId)
			} else {
				var srcListStr string
				for _, val := range sources.SourceList.Sources {
					srcListStr = fmt.Sprintf("%s - %s", val.Name, val.Url)
					fmt.Printf("String representation of the sources value : %s \n", srcListStr)
					t.CallTelegramSendApi(cid, srcListStr)
				}
			}
			c.JSON(http.StatusOK, "OK")
		}
	}
}

func (t *Telegram) Notify(c *gin.Context) {
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		fmt.Printf("Error encountered : %v", err.Error())
		c.JSON(http.StatusOK, err.Error())
		return
	} else {
		reply := new(entities.TelegramReplyMsg)
		err := json.Unmarshal(body, &reply)
		if err != nil {
			fmt.Printf("Could not unmarshal the body. Error encountered : %v", err.Error())
			c.JSON(http.StatusOK, err.Error())
			return
		} else {
			if reply.ChatId == 0 && strings.EqualFold(reply.UserName, "") {
				// send to every subscriber
				if usrDetails, err := t.TelegramDbClient.GetAllUser(); err != nil || len(usrDetails) == 0 {
					fmt.Printf("Could not find user by username. Error encountered : %v", err.Error())
					c.JSON(http.StatusOK, err.Error())
					return
				} else {
					msgTxt := make([]string, 6)
					msgTxt[0] = reply.Text
					if strings.EqualFold(reply.Text, "") {
						qp := make(map[string][]string)
						// qp["top"] = []string{"5"}
						qp["sources"] = []string{"google-news-in"}
						if news, err := t.GNewsService.GetTopNewsBySourceFromDb(qp, 5); err != nil {
							fmt.Printf("Error while content creation. Error encountered : %v", err.Error())
							return
						} else {
							i := 1
							msgTxt[0] = "Tada...Here is the top 5 news for you"
							for _, val := range news {
								msgTxt[i] = val.NewsDescription + " - " + val.NewsUrl
								i++
							}
						}
					}
					for _, usr := range usrDetails {
						// if usr.ID == 1367340022 || usr.ID == 1815027583 {
						cid := strconv.Itoa(int(usr.ChatId))
						for _, txt := range msgTxt {
							t.CallTelegramSendApi(cid, txt)
						}
						// }
					}
				}
			} else if reply.ChatId != 0 && !strings.EqualFold(reply.UserName, "") {
				// dont fetch usd. directly use the chatId
				cid := strconv.Itoa(int(reply.ChatId))
				t.CallTelegramSendApi(cid, reply.Text)
			} else if strings.EqualFold(reply.UserName, "") {
				// dont fetch usd. directly use the chatId
				cid := strconv.Itoa(int(reply.ChatId))
				t.CallTelegramSendApi(cid, reply.Text)
			} else {
				// fetch usd by username
				if usrDet, err := t.TelegramDbClient.GetUserByTgUn(reply.UserName); err != nil || usrDet.ID == 0 {
					fmt.Printf("Could not find user by username. Error encountered : %v", err.Error())
					c.JSON(http.StatusOK, err.Error())
					return
				} else {
					cid := strconv.Itoa(int(usrDet.ChatId))
					t.CallTelegramSendApi(cid, reply.Text)
				}
			}
			c.JSON(http.StatusOK, "OK")
		}
	}
}

func (t *Telegram) CallTelegramSendApi(chatId, text string) error {
	var qp map[string][]string = make(map[string][]string)
	qp["chat_id"] = []string{chatId}
	qp["text"] = []string{text}
	if _, err := t.TelegramClient.Send(qp); err != nil {
		return err
	}
	return nil
}
