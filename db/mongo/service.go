package mongo

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

var (
	mgoDB *mgo.Session
	err   error
)

// ConnectMgo MongoDB Connect
func ConnectMgo() *mgo.Session {
	mongoHost := viper.GetString("mongo.host")
	mongoUser := viper.GetString("mongo.user")
	mongoPass := viper.GetString("mongo.pass")

	connString := fmt.Sprintf("%v:%v@%v", mongoUser, mongoPass, mongoHost)
	mgoDB, err = mgo.Dial(connString)
	if err != nil {
		log.Printf("dial mongodb server with connection string %q: %v", connString, err)
		//return
	}

	return mgoDB
}

// MgoManager return MongoDB Session
func MgoManager() *mgo.Session {
	return mgoDB
}
