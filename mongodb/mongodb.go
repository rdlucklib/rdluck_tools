package mongodb

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	MgoSession *mgo.Session
	MGO_URL    = ""
	MGO_DB     = ""
)

func GetSession() *mgo.Session {
	if MgoSession == nil {
		var err error
		MgoSession, err = mgo.Dial(MGO_URL)
		if err != nil {
			panic(err)
		}
	} else {
		err := MgoSession.Ping()
		if err != nil {
			MgoSession, err = mgo.Dial(MGO_URL)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		}
	}
	//最大连接数默认为4096
	return MgoSession.Clone()
}

//验证id是否符合mongodbId类型
func CheckIsBsonId(_id string) bool {
	return bson.IsObjectIdHex(_id)
}

//公共方法，获取collection对象
func WitchCollection(collection string, s func(*mgo.Collection) error) error {
	session := GetSession()
	defer session.Close()
	c := session.DB(MGO_DB).C(collection)
	return s(c)
}

func init() {
	MgoSession = GetSession()
}
