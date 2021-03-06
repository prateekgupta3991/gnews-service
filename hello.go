package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	CassandraSession := cassandra.NewDbSession(con.Keyspace, hosts)
	defer CassandraSession.Close()

	router.Use(guidMiddleware())

	hClient := GetHttpClient()

	usr := handlers.NewUserBaseService(CassandraSession).UserServ.(*handlers.UserBaseService)
	router.POST("/subscribers", usr.Subscribe)
	router.GET("/subscribers/all", usr.Subscribed)
	gnews := handlers.NewGNews(hClient, CassandraSession).Service.(*handlers.GNewsService)
	router.GET("/sources", gnews.GetSources)
	router.GET("/headlines", gnews.GetHeadlines)
	router.GET("/news", gnews.GetNews)
	tgm := handlers.NewTelegram(hClient, CassandraSession, gnews).IMService.(*handlers.Telegram)
	router.POST("/tgm/updates", tgm.PushedUpdates)
	router.POST("/tgm/reply", tgm.Notify)
	log.Fatal(router.Run(":" + con.ServerPort))
}

func guidMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Set("uuid", uuid)
		fmt.Printf("The request with uuid %s is started \n", uuid)
		c.Next()
		fmt.Printf("The request with uuid %s is served \n", uuid)
	}
}

func GetHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	return &http.Client{
		Transport: tr,
	}
}
