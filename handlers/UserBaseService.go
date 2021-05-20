package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/refresher/entities"
	"github.com/prateekgupta3991/refresher/repository"
)

type UserBaseService struct {
	UserDbClient *repository.UserDbSession
}

func NewUserBaseService() *UserService {
	return &UserService{
		UserServ: &UserBaseService{
			UserDbClient: repository.GetNewUserDbSession(),
		},
	}
}

// subscribed user if new user func
func (u *UserBaseService) Subscribed(c *gin.Context) {
	var respList []entities.UserDetails
	if results, err := u.UserDbClient.GetAllUser(); err != nil {
		fmt.Errorf("Error while fetching subscribers. Returning empty response")
	} else {
		respList = results
	}
	c.JSON(http.StatusOK, respList)
}

func (u *UserBaseService) Subscribe(c *gin.Context) {
	// fmt.Println("Stay tuned. This will be done someday")
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		fmt.Printf("Error encountered : %v", err.Error())
	} else {
		usrDet := new(entities.UserDetails)
		err := json.Unmarshal(body, &usrDet)
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
