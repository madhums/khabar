package db

import (
	"github.com/parthdesai/sc-notifications/utils/time"
	"gopkg.in/mgo.v2/bson"
)

type BaseInterface interface {
	PrepareSave()
}

type BaseModel struct {
	Id         bson.ObjectId `json:"_id,omitempty" bson:"_id" required:"true"`
	CreatedOn  int64         `json:"created_on" bson:"created_on" required:"true"`
	ModifiedOn int64         `json:"updated_on" bson:"updated_on" required:"true"`
}

func (self *BaseModel) PrepareSave() {
	if !self.Id.Valid() {
		self.Id = bson.NewObjectId()
	}

	if !(self.CreatedOn > 0) {
		self.CreatedOn = time.EpochNow()
	}

	self.ModifiedOn = time.EpochNow()
}
