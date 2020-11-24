package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/refresher/cassandra"
	"github.com/prateekgupta3991/refresher/handlers"
)

var (
	router = gin.Default()
)

func main() {
	fmt.Println("Finally back to GO.")
	CassandraSession := cassandra.Session
	defer CassandraSession.Close()
	router.GET("/login", handlers.Login)
	log.Fatal(router.Run(":8080"))
}
