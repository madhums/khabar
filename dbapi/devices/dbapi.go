package devices

import (
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
)

func Get(token string) (err error, found *db.Device) {
	found = new(db.Device)
	err = db.Conn.GetOne(db.DeviceCollection, utils.M{"token": token}, &found)
	return err, found
}

func Insert(item *db.Device) string {
	item.PrepareSave()
	return db.Conn.Insert(db.DeviceCollection, item)
}

func Update(item *db.Device) error {
	item.PrepareSave()
	return db.Conn.Update(db.DeviceCollection, utils.M{"_id": item.Id}, utils.M{"$set": utils.M{"arn": item.EndpointArn}})
}
