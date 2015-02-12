package models

import (
	"github.com/parthdesai/sc-notifications/db"
)

const (
	ChannelCollection = "channels"
)

type Channel struct {
	db.BaseModel   `bson:",inline"`
	UserID         string                 `json:"user_id" bson:"user_id"`
	OrganizationID string                 `json:"org_id" bson:"org_id"`
	ApplicationID  string                 `json:"app_id" bson:"app_id"`
	ChannelData    map[string]interface{} `json:"channel_data" bson:"channel_data" required:"true"`
	Ident          string                 `json:"ident" bson:"ident" required:"true"`
}

func (self *Channel) IsValid() bool {
	if (len(self.UserID) == 0) && (len(self.OrganizationID) == 0) && (len(self.ApplicationID) == 0) {
		return false
	}

	if len(self.Ident) == 0 {
		return false
	}

	if len(self.ChannelData) == 0 {
		return false
	}

	return true
}

func (self *Channel) FindAppropriateChannelForUser(dbConn *db.MConn) bool {
	var hasData bool
	hasData = dbConn.Get(ChannelCollection, db.M{
		"user_id": self.UserID,
		"app_id":  self.ApplicationID,
		"org_id":  self.OrganizationID,
		"ident":   self.Ident,
	}).Next(self)

	if hasData {
		return true
	}

	hasData = dbConn.Get(ChannelCollection, db.M{
		"user_id": self.UserID,
		"app_id":  self.ApplicationID,
		"ident":   self.Ident,
	}).Next(self)

	if hasData {
		return true
	}

	hasData = dbConn.Get(ChannelCollection, db.M{
		"user_id": self.UserID,
		"org_id":  self.OrganizationID,
		"ident":   self.Ident,
	}).Next(self)

	if hasData {
		return true
	}
	return false
}

func (self *Channel) FindAppropriateOrganizationChannel(dbConn *db.MConn) bool {
	var hasData bool
	hasData = dbConn.Get(ChannelCollection, db.M{
		"app_id": self.ApplicationID,
		"org_id": self.OrganizationID,
		"ident":  self.Ident,
	}).Next(self)

	if hasData {
		return true
	}

	hasData = dbConn.Get(ChannelCollection, db.M{
		"org_id": self.OrganizationID,
		"ident":  self.Ident,
	}).Next(self)

	if hasData {
		return true
	}

	return false

}

func (self *Channel) FindGlobalChannel(dbConn *db.MConn) bool {
	var hasData bool
	hasData = dbConn.Get(ChannelCollection, db.M{
		"ident": self.Ident,
	}).Next(self)

	if hasData {
		return true
	}

	return false

}

func (self *Channel) FindAppropriateChannel(dbConn *db.MConn) bool {

	if self.FindAppropriateChannelForUser(dbConn) {
		return true
	}

	if self.FindAppropriateOrganizationChannel(dbConn) {
		return true
	}

	if self.FindGlobalChannel(dbConn) {
		return true
	}

	return false

}

func (self *Channel) GetFromDatabase(dbConn *db.MConn) bool {
	return dbConn.Get(ChannelCollection, db.M{"app_id": self.ApplicationID,
		"org_id": self.OrganizationID, "user_id": self.UserID, "ident": self.Ident}).Next(self)
}

func (self *Channel) DeleteFromDatabase(dbConn *db.MConn) error {
	return dbConn.Delete(ChannelCollection, db.M{"app_id": self.ApplicationID,
		"org_id": self.OrganizationID, "user_id": self.UserID, "ident": self.Ident})
}

func (self *Channel) InsertIntoDatabase(dbConn *db.MConn) string {
	return dbConn.Insert(ChannelCollection, self)
}
