package available_topics

import (
	"log"

	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
)

func GetAppTopics(app_name, org string) []string {
	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{"app_name": app_name}
	topics := []string{}

	err := db.Conn.GetCursor(session, db.AvailableTopicCollection, query).Distinct("ident", topics)
	if err != nil {
		log.Println(err)
	}

	return topics
}

func Get(topic string) (found *db.AvailableTopic, err error) {
	found = new(db.AvailableTopic)

	err = db.Conn.GetOne(db.AvailableTopicCollection, utils.M{"ident": topic}, found)

	if err != nil {
		return nil, err
	}

	return found, nil
}

func Insert(newTopic *db.AvailableTopic) string {
	return db.Conn.Insert(db.AvailableTopicCollection, newTopic)
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.AvailableTopicCollection, *doc)
}
