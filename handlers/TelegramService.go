package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/refresher/clients"
	"github.com/prateekgupta3991/refresher/entities"
	"github.com/prateekgupta3991/refresher/repository"
)

type Telegram struct {
	TelegramClient   *clients.TelegramClient
	TelegramDbClient *repository.UserDbSession
}

func NewTelegram() *Communication {
	return &Communication{
		IMService: &Telegram{
			TelegramClient:   clients.InitTelegramClient(),
			TelegramDbClient: repository.GetNewUserDbSession(),
		},
	}
}

func (t *Telegram) PushedUpdates(c *gin.Context) {
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		fmt.Printf("Error encountered : %v", err.Error())
	} else {
		webhookObj := new(entities.Result)
		err := json.Unmarshal(body, &webhookObj)
		if err != nil {
			fmt.Printf("Could not process the webhook. Error encountered : %v", err.Error())
		} else {
			if subscriber, err := t.TelegramDbClient.GetUserByTgDetils(int(webhookObj.Msg.From.Id), webhookObj.Msg.From.UserName); err != nil || subscriber.ID == 0 {
				fmt.Printf("New subscriber with Id : %d and Username : %s", webhookObj.Msg.From.Id, webhookObj.Msg.From.UserName)
				m := entities.UserDetails{
					ID:         int64(webhookObj.Msg.From.Id),
					Name:       webhookObj.Msg.From.FirstName,
					TelegramId: webhookObj.Msg.From.UserName,
					ChatId:     int32(webhookObj.Msg.Chat.Id),
				}
				if err := t.TelegramDbClient.InsertUser(m); err != nil {
					fmt.Printf("Failure while persisting subscriber with Id : %d and Username : %s - %s", subscriber.ID, subscriber.TelegramId, err.Error())
				}
			} else {
				fmt.Printf("Subscriber found with Id : %d and Username : %s", subscriber.ID, subscriber.TelegramId)
			}
			t.CallTelegramSendApi(strconv.Itoa(int(webhookObj.Msg.Chat.Id)), "tell me")
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
					for _, usr := range usrDetails {
						cid := strconv.Itoa(int(usr.ChatId))
						t.CallTelegramSendApi(cid, reply.Text)
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
