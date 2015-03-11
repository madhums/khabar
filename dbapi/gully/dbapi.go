package gully

import (
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/utils"
	"gopkg.in/simversity/gottp.v2"
)

func Get(dbConn *db.MConn, user string, appName string, org string, ident string) *Gully {
	gully := new(Gully)
	if dbConn.GetOne(GullyCollection, utils.M{"app_name": appName,
		"org": org, "user": user, "ident": ident}, gully) != nil {
		return nil
	}

	return gully
}

func Delete(dbConn *db.MConn, doc *utils.M) error {
	return dbConn.Delete(GullyCollection, *doc)
}

func Insert(dbConn *db.MConn, gully *Gully) string {
	return dbConn.Insert(GullyCollection, gully)
}

func GetAll(dbConn *db.MConn, paginator *gottp.Paginator, user string, appName string, org string, ident string) *[]Gully {
	var query utils.M = make(utils.M)
	if paginator != nil {
		query = *utils.GetPaginationToQuery(paginator)
	}
	var result []Gully

	if len(user) > 0 {
		query["user"] = user
	}

	if len(appName) > 0 {
		query["app_name"] = appName
	}

	if len(org) > 0 {
		query["org"] = org
	}

	if len(ident) > 0 {
		query["ident"] = ident
	}

	var limit int
	var skip int
	var limitExists bool
	var skipExists bool

	limit, limitExists = query["limit"].(int)
	skip, skipExists = query["skip"].(int)

	delete(query, "limit")
	delete(query, "skip")

	if !limitExists {
		limit = 30
	}
	if !skipExists {
		skip = 0
	}

	delete(query, "limit")
	delete(query, "skip")

	dbConn.GetCursor(GullyCollection, query).Skip(skip).Limit(limit).All(&result)

	return &result
}

func findPerUser(dbConn *db.MConn, user string, appName string, org string, ident string) *Gully {
	var err error

	var gully *Gully

	gully = Get(dbConn, user, appName, org, ident)
	if gully != nil {
		return gully
	}

	/*
		Curently, Cannot have the case of App setting without organization.
		err = dbConn.GetOne(GullyCollection, utils.M{
			"user":     user,
			"app_name": appName,
			"org":      "",
			"ident":    ident,
		}, gully)

		if err == nil {
			return gully
		}
	*/

	err = dbConn.GetOne(GullyCollection, utils.M{
		"user":     user,
		"app_name": "",
		"org":      org,
		"ident":    ident,
	}, gully)

	if err == nil {
		return gully
	}

	return nil
}

func findPerOrgnaization(dbConn *db.MConn, appName string, org string, ident string) *Gully {
	var err error
	gully := new(Gully)
	err = dbConn.GetOne(GullyCollection, utils.M{
		"user":     "",
		"app_name": appName,
		"org":      org,
		"ident":    ident,
	}, gully)

	if err == nil {
		return gully
	}

	err = dbConn.GetOne(GullyCollection, utils.M{
		"user":     "",
		"app_name": "",
		"org":      org,
		"ident":    ident,
	}, gully)

	if err == nil {
		return gully
	}

	return nil

}

func findGlobal(dbConn *db.MConn, ident string) *Gully {
	var err error
	gully := new(Gully)
	err = dbConn.GetOne(GullyCollection, utils.M{
		"user":     "",
		"app_name": "",
		"org":      "",
		"ident":    ident,
	}, gully)

	if err == nil {
		return gully
	}

	return nil

}

func FindOne(dbConn *db.MConn, user string, appName string, org string, ident string) *Gully {

	var gully *Gully

	gully = findPerUser(dbConn, user, appName, org, ident)

	if gully != nil {
		return gully
	}

	gully = findPerOrgnaization(dbConn, appName, org, ident)

	if gully != nil {
		return gully
	}

	gully = findGlobal(dbConn, ident)

	if gully != nil {
		return gully
	}

	return nil

}
