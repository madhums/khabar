package gully

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2/bson"
)

const BLANK = ""

//CAUTION: This call does not filter out sensitive information,
//Since it is required by the application.
//DO NOT DIRECTLY WRITE THIS OUTPUT TO USER.
func Get(user, org, ident string) (*db.Gully, error) {
	var gully = new(db.Gully)
	var err error

	err = db.Conn.GetOne(
		db.GullyCollection,
		utils.M{
			"org":   org,
			"user":  user,
			"ident": ident,
		},
		gully,
	)

	if err != nil {
		return nil, err
	}

	return gully, nil
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.GullyCollection, *doc)
}

func Insert(gully *db.Gully) string {
	return db.Conn.Insert(db.GullyCollection, gully)
}

func GetAll(user, org string) (*[]db.Gully, error) {
	var query utils.M = make(utils.M)

	var result []db.Gully

	//Currently, we do not allow to query user level channels
	//Making user value blank so that only channels which are not customized
	//For users will be returned.
	/**if len(user) > 0 {
		query["user"] = user
	} **/
	query["user"] = BLANK

	if len(org) > 0 {
		query["org"] = org
	}

	session := db.Conn.Session.Copy()
	defer session.Close()

	err := db.Conn.GetCursor(session, db.GullyCollection, query).
		Select(bson.M{"data": 0}).All(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func findPerUser(user, org, ident string) (*db.Gully, error) {
	var gully *db.Gully
	var err error

	gully, err = Get(user, org, ident)
	if err == nil {
		return gully, err
	}

	gully, err = Get(user, BLANK, ident)
	return gully, err
}

func findPerOrgnaization(org, ident string) (*db.Gully, error) {
	return Get(BLANK, org, ident)
}

func findGlobal(ident string) (*db.Gully, error) {
	return Get(BLANK, BLANK, ident)
}

func FindOne(user, org, ident string) (gully *db.Gully, err error) {

	gully, err = findPerUser(user, org, ident)
	if err != nil {
		gully, err = findPerOrgnaization(org, ident)
		if err != nil {
			gully, err = findGlobal(ident)
		}
	}

	return
}
