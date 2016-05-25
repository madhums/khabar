/**
 * In the command line
 *
 * ```sh
 * MONGODB_URL=mongodb://localhost:27017/notifications_testing go run 001_modify_available_topics.go
 * ```
 *
 * `MONGODB_URL` env varilable is not needed if you are running on local, but
 * this can be used to connect to different environments
 */

/**
 * What does this do?
 *
 * - Adds `channels` array to `topics_available` collection
 * - Adds `value` property to `topics` collection which determines if the
 * 	 setting is a `default` one on an org level
 */

package migrations

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"path/filepath"
)

type Migrator interface {
	Up(key string, val string)
	Down(key string) string
	Index() int
	Name() string
}

var Migrators map[int]Migrator

func Migrate() {
	fmt.Println("\n", os.Args)

	cwd := os.Getenv("PWD")
	transDir := cwd + "/migrations"

	log.Println("Directory for migrations :" + transDir)

	filepath.Walk(transDir, func(path string, _ os.FileInfo, err error) error {
		fileExt := filepath.Ext(path)
		list := filepath.SplitList(path)
		log.Println(list)
		log.Print("Skipping translation file:" + path + " " +
			"File Extension:" + fileExt + " " + list[len(list)-1])
		return nil
	})
}

func migrationIndex(path string) bool, int {
	fName := filepath.Base(path)
	extName := filepath.Ext(path)
	bName := fName[:len(fName)-len(extName)]
	if (len(bName)> 3) {
		
	}
}

/**
 * Main
 */

func main() {
	session, _ := Connect()
	//MigrateAvailable(db)
	//MigrateTopics(db)
	session.Close()
	fmt.Println("\n", "Closing mongodb connection")
}

/**
 * Connect to mongo
 */

func Connect() (*mgo.Session, *mgo.Database) {
	uri := os.Getenv("MONGODB_URL")

	if uri == "" {
		uri = "mongodb://localhost:27017/notifications_testing"
	}

	mInfo, err := mgo.ParseURL(uri)
	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	session.SetSafe(&mgo.Safe{})
	fmt.Println("Connected to", uri, "\n")

	sess := session.Clone()

	return session, sess.DB(mInfo.Database)
}

/**
 * Handle errors
 */

func handle_errors(err error) {
	if err != nil {
		log.Printf("Error %v\n", err)
		os.Exit(1)
	}
}
