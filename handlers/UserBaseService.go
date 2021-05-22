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
	rId, _ := c.Get("uuid")
	fmt.Printf("The request with uuid %s is served succesfully \n", rId)
	c.JSON(http.StatusOK, respList)
}

func (u *UserBaseService) Subscribe(c *gin.Context) {
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		fmt.Printf("Error encountered : %v", err.Error())
	} else {
		usrDet := new(entities.UserDetails)
		err := json.Unmarshal(body, &usrDet)
		if err != nil {
			fmt.Printf("Could not process the subscription request. Error encountered : %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
		} else {
			if err := u.CheckAndPersist(usrDet); err != nil {
				c.JSON(http.StatusOK, gin.H{"Error: ": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"Success": "true"})
		}
	}
}

func (u *UserBaseService) CheckAndPersist(usrDet *entities.UserDetails) error {
	if subscriber, err := u.UserDbClient.GetUserByTgDetils(int(usrDet.ID), usrDet.TelegramId); err != nil || subscriber.ID == 0 {
		fmt.Printf("New subscriber with Id : %d and Username : %s", usrDet.ID, usrDet.TelegramId)
		m := entities.UserDetails{
			ID:         usrDet.ID,
			Name:       usrDet.Name,
			TelegramId: usrDet.TelegramId,
			ChatId:     usrDet.ChatId,
		}
		if err := u.UserDbClient.InsertUser(m); err != nil {
			fmt.Printf("Failure while persisting subscriber with Id : %d and Username : %s - %s", subscriber.ID, subscriber.TelegramId, err.Error())
			return err
		}
	} else {
		fmt.Printf("Subscriber found with Id : %d and Username : %s", subscriber.ID, subscriber.TelegramId)
		return err
	}
	return nil
}
