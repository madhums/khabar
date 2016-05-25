// +build ignore

/**
 * In the command line
 *
 * ```sh
 * MONGODB_URL=mongodb://localhost:27017/notifications_testing go run 003_clear_org_user_settings.go
 * ```
 *
 * `MONGODB_URL` env varilable is not needed if you are running on local, but
 * this can be used to connect to different environments
 */

/**
 * What does this do?
 *
 * - Removes org and user settings from topics collection
 */

package migrations

import (
	// "fmt"
	// "log"
	// "os"

	// "github.com/bulletind/khabar/db"
	// "github.com/bulletind/khabar/utils"

	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

init() {
	//Migrators = make(map[int]Migrator)
	fmt.Println("init 004")
}

/**
 * Main
 */
func main() {
	session, db, _ := Connect()
	defer session.Close()
	fmt.Println("\n", "Closing mongodb connection woeha")
}

// func main2() {
// 	session, db, _ := Connect()
// 	defer session.Close()

// 	for collection, keys := range db.GetIndexes() {
// 		index := mgo.Index{
// 			Key:        strings.Split(keys, " "),
// 			Unique:     true,
// 			DropDups:   true,
// 			Background: true,
// 			Sparse:     true,
// 		}

// 		err := db.C(collection).EnsureIndex(index)
// 		if err != nil {
// 			log.Println("Error creating index:", err)
// 			panic(err)
// 		}
// 	}

// 	Topics := db.C(topicsCollection)
// 	query := bson.M{
// 		"org":  utils.M{"$ne": ""},
// 		"user": utils.M{"$ne": ""},
// 	}

// 	// Remove user preference
// 	change, err := Topics.RemoveAll(query)
// 	handle_errors(err)
// 	fmt.Println(change.Removed, "user preferences  were removed from `", topicsCollection, "` collection")

// 	// Remove org preference
// 	query["user"] = ""
// 	change, err = Topics.RemoveAll(query)
// 	handle_errors(err)
// 	fmt.Println(change.Removed, "org preferences  were removed from `", topicsCollection, "` collection")

// 	fmt.Println("\n", "Closing mongodb connection")
}
