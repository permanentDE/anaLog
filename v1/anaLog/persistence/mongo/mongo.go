package mongo

import (
	idl "go.iondynamics.net/iDlogger"
	"gopkg.in/mgo.v2"

	"go.permanent.de/anaLog/v1/config"
	//"gopkg.in/mgo.v2/bson"
)

var mgoSession *mgo.Session

func connect2mgo() {
	var err error
	mgoSession, err = mgo.Dial(config.Mongo.Host)

	if err == nil {
		mgoSession.SetMode(mgo.Monotonic, true)
		if config.Mongo.User != "" && config.Mongo.Pass != "" {
			err = mgoSession.Login(&mgo.Credential{
				Username: config.Mongo.User,
				Password: config.Mongo.Pass,
			})
		}
	}
	if err != nil {
		idl.Panic(err)
	}
}

func getMgoSession() *mgo.Session {
	if mgoSession == nil {
		connect2mgo()
	}
	return mgoSession
}

func Close() {
	mgoSession.Close()
}
