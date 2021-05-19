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
		webhookObj := new(entities.Webhook)
		err := json.Unmarshal(body, &webhookObj)
		if err != nil {
			fmt.Printf("Could not process the webhook. Error encountered : %v", err.Error())
		} else {
			if webhookObj.Ok {
				for _, val := range webhookObj.Res {
					if subscriber, err := t.TelegramDbClient.GetUserByTgDetils(int(val.Msg.From.Id), val.Msg.From.UserName); err != nil {
						fmt.Printf("New subscriber with Id : %s and Username : %s", val.Msg.From.Id, val.Msg.From.UserName)
						m := entities.UserDetails{
							ID:         int64(val.Msg.From.Id),
							Name:       val.Msg.From.FirstName,
							TelegramId: val.Msg.From.UserName,
							ChatId:     int32(val.Msg.Chat.Id),
						}
						t.TelegramDbClient.InsertUser(m)
					} else {
						fmt.Printf("Subscriber with Id : %s and Username : %s", subscriber.ID, subscriber.TelegramId)
					}
					// send a reply to subscriber or m
				}
			}
		}
	}
}

func (t *Telegram) Notify(c *gin.Context) {
	// send a message to telegram as bot reply
	fmt.Printf("I will notify tomorrow")
}
