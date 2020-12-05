package main

import (
	"fmt"
	"log"
	"strings"

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
		log.Panic("Error during config initialisation : " + err.Error())
	}
	hosts := strings.Split(con.CasDb, ",")
	cassandra.InitDb(con.Keyspace, hosts)
	CassandraSession := cassandra.Session
	defer CassandraSession.Close()
	router.GET("/login", handlers.Login)
	gnews := handlers.GetNewGNews().Service.(*handlers.GNewsService)
	router.GET("/sources", gnews.GetSources)
	router.GET("/headlines", gnews.GetHeadlines)
	router.GET("/news", gnews.GetNews)
	log.Fatal(router.Run(":" + con.ServerPort))
}
