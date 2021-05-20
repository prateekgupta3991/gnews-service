package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/prateekgupta3991/refresher/cassandra"
	"github.com/prateekgupta3991/refresher/entities"
)

// Login func
func Login(c *gin.Context) {
	var u entities.User
	m := map[string]interface{}{}
	results := cassandra.Session.Query("select * from emp").Iter()
	var uid gocql.UUID
	var uname string
	for results.MapScan(m) {
		uid = m["id"].(gocql.UUID)
		uname = m["name"].(string)
		break
	}
	u.ID = uid.String()
	u.Name = uname
	c.JSON(http.StatusOK, u)
}

// subscribed user if new user func
func SubscribedUser(c *gin.Context) {
	var respList []entities.UserDetails
	m := map[string]interface{}{}
	results := cassandra.Session.Query("select * from user").Iter()
	var uid int64
	var uname string
	var tId string
	var cId int32
	for results.MapScan(m) {
		uid = m["id"].(int64)
		uname = m["name"].(string)
		tId = m["t_un"].(string)
		cId = m["chat_id"].(int32)
		var u entities.UserDetails
		u.ID = uid
		u.Name = uname
		u.TelegramId = tId
		u.ChatId = cId
		respList = append(respList, u)
	}
	c.JSON(http.StatusOK, respList)
}
