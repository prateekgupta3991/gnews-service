package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitDb(ks string, db []string) {
	var err error
	cluster := gocql.NewCluster("192.168.76.203:9042")
	cluster.Keyspace = ks
	cluster.Consistency = gocql.One
	// cluster.Timeout = 10000
	if Session, err = cluster.CreateSession(); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("cassandra init done")
}
