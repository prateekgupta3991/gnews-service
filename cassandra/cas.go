package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

func NewDbSession(ks string, db []string) *gocql.Session {
	var Session *gocql.Session
	var err error
	cluster := gocql.NewCluster("172.17.0.2:9042")
	cluster.Keyspace = ks
	if Session, err = cluster.CreateSession(); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("cassandra init done")
	return Session
}
