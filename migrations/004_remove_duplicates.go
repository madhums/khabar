/**
 * In the command line
 *
 * ```sh
 * MONGODB_URL=mongodb://localhost:27017/notifications_testing go run migrations/004_remove_duplicates.go
 * ```
 *
 * `MONGODB_URL` env varilable is not needed if you are running on local, but
 * this can be used to connect to different environments
 */

/**
 * What does this do?
 *
 * - Removes duplicates for indexes we are going to enforce
 */

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bulletind/khabar/db"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/**
 * Main
 */

func main() {
	session, database, _ := Connect()
	defer session.Close()

	for collection, keys := range db.GetIndexes() {
		var fields = bson.M{}

		//grouping
		for _, field := range strings.Split(keys, " ") {
			fields[field] = "$" + field
		}

		var pipes = make([]bson.M, 3)
		pipes[0] = bson.M{"$group": bson.M{"_id": fields, "dups": bson.M{"$sum": 1}}}
		pipes[1] = bson.M{"$match": bson.M{"dups": bson.M{"$gte": 2}}}
		pipes[2] = bson.M{"$sort": bson.M{"dups": -1}}

		pipe := database.C(collection).Pipe(pipes)
		result := []bson.M{}
		err := pipe.All(&result)
		if err != nil {
			fmt.Println(err, pipes)
		}

		// loop through results
		for _, row := range result {
			dups := row["dups"].(int)

			// extra ordinary
			if dups > 2 {
				fmt.Println("Number of dups: ", dups, row["_id"])
			}

			// remove item per item
			for i := 0; i < dups-1; i++ {
				errDelete := database.C(collection).Remove(row["_id"])
				if errDelete != nil {
					fmt.Println(errDelete, row["_id"])
				}
			}
		}

		fmt.Println(collection, " --duplicates-- ", len(result))
	}
}

/**
 * Connect to mongo
 */

func Connect() (*mgo.Session, *mgo.Database, string) {
	uri := os.Getenv("MONGODB_URL")

	if uri == "" {
		uri = "mongodb://localhost:27017/notifications_testing"
	}

	const SSL_SUFFIX = "?ssl=true"
	connStringForDb := uri
	if strings.HasSuffix(uri, SSL_SUFFIX) {
		connStringForDb = strings.TrimSuffix(uri, SSL_SUFFIX)
	}

	mInfo, _ := mgo.ParseURL(connStringForDb)
	conn := db.GetConn(uri, mInfo.Database)
	session := conn.Session
	session.SetSafe(&mgo.Safe{})
	fmt.Println("Connected to", uri, "\n")

	sess := session.Clone()

	return session, sess.DB(mInfo.Database), mInfo.Database
}
