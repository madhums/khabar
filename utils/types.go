package utils

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v3/utils"
)

type M bson.M

type R struct {
	StatusCode int
	Data       interface{}
	Message    string
}

func (self R) SendOverWire() utils.Q {
	return utils.Q{
		"status":  self.StatusCode,
		"data":    self.Data,
		"message": self.Message,
	}
}
