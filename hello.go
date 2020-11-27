package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/refresher/cassandra"
	"github.com/prateekgupta3991/refresher/configs"
	"github.com/prateekgupta3991/refresher/handlers"
)

var (
	router = gin.Default()
)

func main() {
	fmt.Println("Finally back to GO.")
	var con *configs.Conf
	var err error
	if con, err = configs.InitConfig("./configs/conf.dev.json"); err != nil {
		log.Panic("Error during config initialisation")
	}
	CassandraSession := cassandra.Session
	defer CassandraSession.Close()
	router.GET("/login", handlers.Login)
	log.Fatal(router.Run(":" + con.ServerPort))
}
