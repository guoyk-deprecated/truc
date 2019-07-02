package extmgo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/yankeguo/truc/ext/extos"
	"strings"
)

func DialLinkedMongo() (*mgo.Session, error) {
	return DialLinkedMongoWithAlias("mongo")
}

func DialLinkedMongoWithAlias(alias string) (*mgo.Session, error) {
	var host = "localhost"
	var port = "27017"
	var username, password string

	alias = strings.ToUpper(strings.TrimSpace(alias))

	extos.EnvStr(&host, alias+"_PORT_27017_TCP_ADDR", "MONGO_HOST")
	extos.EnvStr(&port, alias+"_PORT_27017_TCP_PORT", "MONGO_PORT")
	extos.EnvStr(&username, alias+"_ENV_MONGO_INITDB_ROOT_USERNAME", "MONGO_USERNAME")
	extos.EnvStr(&password, alias+"_ENV_MONGO_INITDB_ROOT_PASSWORD", "MONGO_PASSWORD")

	if len(username) > 0 && len(password) > 0 {
		return mgo.Dial(fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port))
	} else {
		return mgo.Dial(fmt.Sprintf("mongodb:/%s:%s", host, port))
	}
}
