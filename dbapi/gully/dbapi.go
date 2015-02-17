package gully

import (
	"github.com/parthdesai/sc-notifications/db"
)

func GetFromDatabase(dbConn *db.MConn, userID string, applicationID string, organizationID string, ident string) *Gully {
	gully := new(Gully)
	if !dbConn.Get(GullyCollection, db.M{"app_id": applicationID,
		"org_id": organizationID, "user_id": userID, "ident": ident}).Next(gully) {
		return nil
	}
	return gully
}

func DeleteFromDatabase(dbConn *db.MConn, gully *Gully) error {
	return dbConn.Delete(GullyCollection, db.M{"app_id": gully.ApplicationID,
		"org_id": gully.OrganizationID, "user_id": gully.UserID, "ident": gully.Ident})
}

func InsertIntoDatabase(dbConn *db.MConn, gully *Gully) string {
	return dbConn.Insert(GullyCollection, gully)
}

func FindAppropriateGullyForUser(dbConn *db.MConn, userID string, applicationID string, organizationID string, ident string) *Gully {
	var hasData bool
	gully := new(Gully)
	hasData = dbConn.Get(GullyCollection, db.M{
		"user_id": userID,
		"app_id":  applicationID,
		"org_id":  organizationID,
		"ident":   ident,
	}).Next(gully)

	if hasData {
		return gully
	}

	hasData = dbConn.Get(GullyCollection, db.M{
		"user_id": userID,
		"app_id":  applicationID,
		"ident":   ident,
	}).Next(gully)

	if hasData {
		return gully
	}

	hasData = dbConn.Get(GullyCollection, db.M{
		"user_id": userID,
		"org_id":  organizationID,
		"ident":   ident,
	}).Next(gully)

	if hasData {
		return gully
	}
	return nil
}

func FindAppropriateOrganizationGully(dbConn *db.MConn, applicationID string, organizationID string, ident string) *Gully {
	var hasData bool
	gully := new(Gully)
	hasData = dbConn.Get(GullyCollection, db.M{
		"app_id": applicationID,
		"org_id": organizationID,
		"ident":  ident,
	}).Next(gully)

	if hasData {
		return gully
	}

	hasData = dbConn.Get(GullyCollection, db.M{
		"org_id": organizationID,
		"ident":  ident,
	}).Next(gully)

	if hasData {
		return gully
	}

	return nil

}

func FindGlobalGully(dbConn *db.MConn, ident string) *Gully {
	var hasData bool
	gully := new(Gully)
	hasData = dbConn.Get(GullyCollection, db.M{
		"ident": ident,
	}).Next(gully)

	if hasData {
		return gully
	}

	return nil

}

func FindAppropriateGully(dbConn *db.MConn, userID string, applicationID string, organizationID string, ident string) *Gully {

	var gully *Gully

	gully = FindAppropriateGullyForUser(dbConn, userID, applicationID, organizationID, ident)

	if gully != nil {
		return gully
	}

	gully = FindAppropriateOrganizationGully(dbConn, applicationID, organizationID, ident)

	if gully != nil {
		return gully
	}

	gully = FindGlobalGully(dbConn, ident)

	if gully != nil {
		return gully
	}

	return nil

}
