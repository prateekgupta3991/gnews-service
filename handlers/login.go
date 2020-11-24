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
	results := cassandra.Session.Query("select * from test").Iter()
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
