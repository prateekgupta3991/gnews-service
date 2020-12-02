package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitDb(ks string, db []string) {
	var err error
	cluster := gocql.NewCluster("172.18.9.140:9042", "172.18.9.140:9043")
	cluster.Keyspace = ks
	if Session, err = cluster.CreateSession(); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("cassandra init done")
}
