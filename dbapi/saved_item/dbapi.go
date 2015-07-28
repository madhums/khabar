package saved_item

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
)

func Insert(coll string, savedItem *db.SavedItem) string {
	savedItem.PrepareSave()
	return db.Conn.Insert(coll, savedItem)
}

func Get(coll string, query *utils.M) (savedItem *db.SavedItem, err error) {
	savedItem = new(db.SavedItem)

	err = db.Conn.GetOne(coll, *query, savedItem)

	if err != nil {
		return nil, err
	}

	return savedItem, nil
}

func GetSentOrganizations(coll string, email string) (string, []string) {
	session := db.Conn.Session.Copy()
	defer session.Close()

	orgs := []string{}
	userId := ""

	iter := db.Conn.GetCursor(session, coll, utils.M{"details.context.email": email}).Sort("-_id").Iter()

	var one struct {
		Details struct {
			Organization string `bson:"org"`
			User         string `bson:"user"`
		} `bson:"details"`
	}

	for iter.Next(&one) {
		if userId == "" {
			userId = one.Details.User
		}

		if one.Details.User != userId {
			continue
		}

		if !db.InArray(one.Details.Organization, orgs) {
			orgs = append(orgs, one.Details.Organization)
		}
	}

	return userId, orgs
}
