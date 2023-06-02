package cassandra

import (
	"fmt"
	"os"

	"github.com/gocql/gocql"
)

func NewDbSession(ks string, db []string) *gocql.Session {
	var Session *gocql.Session
	var err error
	dbHost := os.Getenv("DATABASE_HOST")
	if dbHost == "" {
		fmt.Println("DATABASE_HOST env variable not set")
		panic("error")
	}
	cluster := gocql.NewCluster(dbHost)
	cluster.Keyspace = ks
	if Session, err = cluster.CreateSession(); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("cassandra init done")
	return Session
}
