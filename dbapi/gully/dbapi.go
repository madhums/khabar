package gully

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
)

const BLANK = ""

func Get(user, appName, org,
	ident string) (gully *db.Gully, err error) {
	gully = new(db.Gully)
	err = db.Conn.GetOne(
		db.GullyCollection,
		utils.M{
			"app_name": appName,
			"org":      org,
			"user":     user,
			"ident":    ident,
		},
		gully,
	)

	if err != nil {
		return nil, err
	}

	return
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.GullyCollection, *doc)
}

func Insert(gully *db.Gully) string {
	return db.Conn.Insert(db.GullyCollection, gully)
}

func GetAll(user, appName, org string) (*[]db.Gully, error) {
	var query utils.M = make(utils.M)

	var result []db.Gully

	if len(user) > 0 {
		query["user"] = user
	}

	if len(appName) > 0 {
		query["app_name"] = appName
	}

	if len(org) > 0 {
		query["org"] = org
	}

	session := db.Conn.Session.Copy()
	defer session.Close()

	err := db.Conn.GetCursor(session, db.GullyCollection, query).All(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func findPerUser(user, appName, org,
	ident string) (gully *db.Gully, err error) {

	gully, err = Get(user, appName, org, ident)
	if err != nil {
		gully, err = Get(user, BLANK, org, ident)
		if err != nil {
			err = db.Conn.GetOne(db.GullyCollection, utils.M{
				"user":     user,
				"app_name": BLANK,
				"org":      org,
				"ident":    ident,
			}, gully)
		}
	}

	/*
		Curently, Cannot have the case of App setting without organization.
		err = db.Conn.GetOne(db.GullyCollection, utils.M{
			"user":     user,
			"app_name": appName,
			"org":      BLANK,
			"ident":    ident,
		}, gully)

		if err == nil {
			return gully
		}
	*/

	return
}

func findPerOrgnaization(appName, org,
	ident string) (gully *db.Gully, err error) {

	gully = new(db.Gully)
	err = db.Conn.GetOne(db.GullyCollection, utils.M{
		"user":     BLANK,
		"app_name": appName,
		"org":      org,
		"ident":    ident,
	}, gully)

	if err != nil {
		err = db.Conn.GetOne(db.GullyCollection, utils.M{
			"user":     BLANK,
			"app_name": BLANK,
			"org":      org,
			"ident":    ident,
		}, gully)
	}

	return
}

func findGlobal(appName, ident string) (gully *db.Gully, err error) {
	gully = new(db.Gully)
	err = db.Conn.GetOne(db.GullyCollection, utils.M{
		"user":     BLANK,
		"app_name": appName,
		"org":      BLANK,
		"ident":    ident,
	}, gully)

	if err != nil {
		err = db.Conn.GetOne(db.GullyCollection, utils.M{
			"user":     BLANK,
			"app_name": BLANK,
			"org":      BLANK,
			"ident":    ident,
		}, gully)

	}

	return

}

func FindOne(user, appName, org, ident string) (gully *db.Gully, err error) {

	gully, err = findPerUser(user, appName, org, ident)
	if err != nil {
		gully, err = findPerOrgnaization(appName, org, ident)
		if err != nil {
			gully, err = findGlobal(appName, ident)
		}
	}

	return
}
