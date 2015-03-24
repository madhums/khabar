package gully

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
)

const BLANK = ""

func Get(dbConn *db.MConn, user, appName, org, ident string) (gully *db.Gully, err error) {
	gully = new(db.Gully)
	err = dbConn.GetOne(
		db.GullyCollection,
		utils.M{
			"app_name": appName,
			"org":      org,
			"user":     user,
			"ident":    ident,
		},
		gully,
	)

	return
}

func Delete(dbConn *db.MConn, doc *utils.M) error {
	return dbConn.Delete(db.GullyCollection, *doc)
}

func Insert(dbConn *db.MConn, gully *db.Gully) string {
	return dbConn.Insert(db.GullyCollection, gully)
}

func GetAll(dbConn *db.MConn, user, appName, org string) (*[]db.Gully, error) {
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

	session := dbConn.Session.Copy()
	defer session.Close()

	err := dbConn.GetCursor(session, db.GullyCollection, query).All(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func findPerUser(dbConn *db.MConn, user, appName, org, ident string) (gully *db.Gully, err error) {

	gully, err = Get(dbConn, user, appName, org, ident)
	if err != nil {
		gully, err = Get(dbConn, user, BLANK, org, ident)
		if err != nil {
			err = dbConn.GetOne(db.GullyCollection, utils.M{
				"user":     user,
				"app_name": BLANK,
				"org":      org,
				"ident":    ident,
			}, gully)
		}
	}

	/*
		Curently, Cannot have the case of App setting without organization.
		err = dbConn.GetOne(db.GullyCollection, utils.M{
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

func findPerOrgnaization(dbConn *db.MConn, appName, org, ident string) (gully *db.Gully, err error) {

	gully = new(db.Gully)
	err = dbConn.GetOne(db.GullyCollection, utils.M{
		"user":     BLANK,
		"app_name": appName,
		"org":      org,
		"ident":    ident,
	}, gully)

	if err != nil {
		err = dbConn.GetOne(db.GullyCollection, utils.M{
			"user":     BLANK,
			"app_name": BLANK,
			"org":      org,
			"ident":    ident,
		}, gully)
	}

	return
}

func findGlobal(dbConn *db.MConn, appName, ident string) (gully *db.Gully, err error) {
	gully = new(db.Gully)
	err = dbConn.GetOne(db.GullyCollection, utils.M{
		"user":     BLANK,
		"app_name": appName,
		"org":      BLANK,
		"ident":    ident,
	}, gully)

	if err != nil {
		err = dbConn.GetOne(db.GullyCollection, utils.M{
			"user":     BLANK,
			"app_name": BLANK,
			"org":      BLANK,
			"ident":    ident,
		}, gully)

	}

	return

}

func FindOne(dbConn *db.MConn, user, appName, org, ident string) (gully *db.Gully, err error) {

	gully, err = findPerUser(dbConn, user, appName, org, ident)
	if err != nil {
		gully, err = findPerOrgnaization(dbConn, appName, org, ident)
		if err != nil {
			gully, err = findGlobal(dbConn, appName, ident)
		}
	}

	return
}
