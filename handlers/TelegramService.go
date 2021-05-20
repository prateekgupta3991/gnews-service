package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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
			// fix this
			// t.Notify(c)
		}
	}
}

func (t *Telegram) Notify(c *gin.Context) {
	var qp map[string][]string = make(map[string][]string)
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		fmt.Printf("Error encountered : %v", err.Error())
		c.JSON(http.StatusOK, err.Error())
		return
	} else {
		// fix the contract
		reply := new(entities.TelegramReplyMsg)
		err := json.Unmarshal(body, &reply)
		if err != nil {
			fmt.Printf("Could not unmarshal the body. Error encountered : %v", err.Error())
			c.JSON(http.StatusOK, err.Error())
			return
		} else {
			// fix this please
			// if usrDet, err := t.TelegramDbClient.GetUserByTgDetils(reply.ChatId, reply.UserName); err != nil || usrDet.ID == 0 {
			// 	fmt.Printf("Could not find user by username. Error encountered : %v", err.Error())
			// 	c.JSON(http.StatusOK, err.Error())
			// 	return
			// } else {
			msgId := strconv.Itoa(int(reply.ChatId))
			qp["chat_id"] = []string{msgId}
			qp["text"] = []string{reply.Text}
			t.TelegramClient.Send(qp)
			// }
			c.JSON(http.StatusOK, "OK")
		}
	}
}
