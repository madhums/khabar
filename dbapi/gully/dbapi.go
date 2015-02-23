package gully

import (
	"github.com/changer/sc-notifications/db"
)

func Get(dbConn *db.MConn, user string, appName string, org string, ident string) *Gully {
	gully := new(Gully)
	if dbConn.GetOne(GullyCollection, db.M{"app_name": appName,
		"org": org, "user": user, "ident": ident}, gully) != nil {
		return nil
	}
	return gully
}

func Delete(dbConn *db.MConn, doc *db.M) error {
	return dbConn.Delete(GullyCollection, *doc)
}

func Insert(dbConn *db.MConn, gully *Gully) string {
	return dbConn.Insert(GullyCollection, gully)
}

func FindAppropriateGullyForUser(dbConn *db.MConn, user string, appName string, org string, ident string) *Gully {
	var err error
	gully := new(Gully)
	err = dbConn.GetOne(GullyCollection, db.M{
		"user":     user,
		"app_name": appName,
		"org":      org,
		"ident":    ident,
	}, gully)

	if err == nil {
		return gully
	}

	err = dbConn.GetOne(GullyCollection, db.M{
		"user":     user,
		"app_name": appName,
		"ident":    ident,
	}, gully)

	if err == nil {
		return gully
	}

	err = dbConn.GetOne(GullyCollection, db.M{
		"user":  user,
		"org":   org,
		"ident": ident,
	}, gully)

	if err == nil {
		return gully
	}
	return nil
}

func FindAppropriateOrganizationGully(dbConn *db.MConn, appName string, org string, ident string) *Gully {
	var err error
	gully := new(Gully)
	err = dbConn.GetOne(GullyCollection, db.M{
		"app_name": appName,
		"org":      org,
		"ident":    ident,
	}, gully)

	if err == nil {
		return gully
	}

	err = dbConn.GetOne(GullyCollection, db.M{
		"org":   org,
		"ident": ident,
	}, gully)

	if err == nil {
		return gully
	}

	return nil

}

func FindGlobalGully(dbConn *db.MConn, ident string) *Gully {
	var err error
	gully := new(Gully)
	err = dbConn.GetOne(GullyCollection, db.M{
		"ident": ident,
	}, gully)

	if err == nil {
		return gully
	}

	return nil

}

func FindAppropriateGully(dbConn *db.MConn, user string, appName string, org string, ident string) *Gully {

	var gully *Gully

	gully = FindAppropriateGullyForUser(dbConn, user, appName, org, ident)

	if gully != nil {
		return gully
	}

	gully = FindAppropriateOrganizationGully(dbConn, appName, org, ident)

	if gully != nil {
		return gully
	}

	gully = FindGlobalGully(dbConn, ident)

	if gully != nil {
		return gully
	}

	return nil

}
