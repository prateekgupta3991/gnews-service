package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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
		}
	}
}

func (t *Telegram) Notify(c *gin.Context) {
	// send a message to telegram as bot reply
	fmt.Printf("I will notify tomorrow")
}
