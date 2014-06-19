package domain

import (
  "os"
	"labix.org/v2/mgo"
)

var mgoSession   *mgo.Session

func GetSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(os.Getenv("OPENSHIFT_MONGODB_DB_URL"))
		if err != nil {
			panic(err) // no, not really
		}
	}
	return mgoSession.Clone()
}