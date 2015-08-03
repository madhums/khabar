package db

import (
	"errors"
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/bulletind/khabar/utils"
)

var Conn *MConn

type MConn struct {
	Session *mgo.Session
	Dbname  string
}

func (self *MConn) getCursor(session *mgo.Session, table string,
	query utils.M) *mgo.Query {

	fields, err1 := query["fields"].(utils.M)
	delete(query, "fields")
	if !err1 {
		fields = utils.M{}
	}

	sort, err2 := query["sort"].(string)
	delete(query, "sort")
	if !err2 {
		sort = "$natural"
	}

	skip, err3 := query["skip"].(int)
	delete(query, "skip")
	if !err3 {
		skip = 0
	}

	limit, err4 := query["limit"].(int)
	delete(query, "limit")
	if !err4 {
		limit = 0
	}

	cursor := self.GetCursor(session, table, query)
	return cursor.Limit(limit).Skip(skip).Sort(sort).Select(fields)
}

func (self *MConn) findAndApply(
	table string, query utils.M, change mgo.Change, result interface{},
) (*mgo.ChangeInfo, error) {
	//Create a Session Copy and be responsible for Closing it.
	session := self.Session.Copy()
	db := session.DB(self.Dbname)
	defer session.Close()

	change.ReturnNew = true

	coll := db.C(table)
	info, err := coll.Find(query).Apply(change, result)
	if err != nil {
		log.Println("Error Applying Changes", table, err)
	}
	return info, err
}

func (self *MConn) FindAndUpdate(
	table string, query utils.M, doc utils.M, result interface{},
) (*mgo.ChangeInfo, error) {
	change := mgo.Change{
		Update: doc,
		Upsert: true,
	}
	return self.findAndApply(table, query, change, result)
}

func (self *MConn) GetCursor(session *mgo.Session, table string,
	query utils.M) *mgo.Query {
	db := session.DB(self.Dbname)

	coll := db.C(table)
	out := coll.Find(query)

	return out
}

func (self *MConn) Get(session *mgo.Session, table string,
	query utils.M) *mgo.Iter {
	return self.getCursor(session, table, query).Iter()
}

func (self *MConn) GetOne(table string, query utils.M,
	result interface{}) error {
	//Create a Session Copy and be responsible for Closing it.
	session := self.Session.Copy()
	defer session.Close()

	cursor := self.getCursor(session, table, query)
	err := cursor.One(result)
	if err != nil {
		log.Println("Error fetching", table, err)
	}

	return err
}

func (self *MConn) Count(table string, query utils.M) int {
	//Create a Session Copy and be responsible for Closing it.
	session := self.Session.Copy()
	defer session.Close()

	cursor := self.getCursor(session, table, query).Select(utils.M{"_id": 1})
	count, err := cursor.Count()
	if err != nil {
		log.Println("Error Counting", table, err)
	}

	return count
}

func (self *MConn) Upsert(table string, query utils.M, doc utils.M) error {

	//Create a Session Copy and be responsible for Closing it.
	session := self.Session.Copy()
	db := session.DB(self.Dbname)
	defer session.Close()

	var err error
	if len(doc) == 0 {
		err = errors.New(
			"Empty upsert is blocked. Refer to " +
				"https://github.com/Simversity/blackjack/issues/1051",
		)
	} else {
		coll := db.C(table)
		_, err = coll.Upsert(query, doc)
	}

	if err != nil {
		log.Println("Error Upserting:", table, err)
	}
	return err
}

func AlterDoc(doc *utils.M, operator string, operation utils.M) {
	spec := *doc
	if spec[operator] != nil {
		op, _ := spec[operator].(utils.M)
		for key, value := range op {
			operation[key] = value
		}
	}
	spec[operator] = operation
}

func (self *MConn) Update(table string, query utils.M, doc utils.M) error {

	//Create a Session Copy and be responsible for Closing it.
	session := self.Session.Copy()
	db := session.DB(self.Dbname)
	defer session.Close()

	coll := db.C(table)
	var update_err error
	if len(doc) == 0 {
		update_err = errors.New(
			"Empty Update is blocked. Refer to " +
				"https://github.com/Simversity/blackjack/issues/1051",
		)
	} else {
		AlterDoc(&doc, "$set", utils.M{"updated_on": utils.EpochNow()})
		_, update_err = coll.UpdateAll(query, doc)
	}

	if update_err != nil {
		log.Println("Error Updating:", table, update_err)
	}
	return update_err
}

func (self *MConn) Delete(table string, query utils.M) error {
	//Create a Session Copy and be responsible for Closing it.
	session := self.Session.Copy()
	db := session.DB(self.Dbname)
	defer session.Close()

	var delete_err error

	coll := db.C(table)

	_, delete_err = coll.RemoveAll(query)

	if delete_err != nil {
		log.Println("Error Deleting:", table, delete_err)
	}

	return delete_err
}

func InArray(key string, arrays ...[]string) bool {
	for _, val := range arrays {
		for _, one := range val {
			if key == one {
				return true
			}
		}
	}
	return false
}

func (self *MConn) InsertMulti(table string, arguments ...interface{}) (error, *mgo.BulkResult) {
	session := self.Session.Copy()
	db := session.DB(self.Dbname)
	defer session.Close()

	b := db.C(table).Bulk()
	b.Insert(arguments...)
	bulkResult, err := b.Run()
	if err != nil {
		return err, bulkResult
	}
	return nil, bulkResult
}

func (self *MConn) Insert(table string, arguments ...interface{}) (_id string) {
	//Create a Session Copy and be responsible for Closing it.
	session := self.Session.Copy()
	db := session.DB(self.Dbname)
	defer session.Close()

	var out interface{}
	if len(arguments) > 1 {
		out = arguments[1]
	} else {
		out = nil
	}

	doc := arguments[0]

	coll := db.C(table)
	err := coll.Insert(doc)
	if err != nil {
		panic(err)
	}

	if out != nil {
		stream, merr := bson.Marshal(doc)
		if merr == nil {
			bson.Unmarshal(stream, out)
		}
	}

	return
}

var cached = struct {
	sync.RWMutex
	sessions map[string]*mgo.Session
}{sessions: map[string]*mgo.Session{}}

func GetConn(db_name string, address string, creds ...string) *MConn {
	//Check if the connection has been stored already.
	var session *mgo.Session
	var ok bool

	cached.RLock()
	session, ok = cached.sessions[db_name]
	cached.RUnlock()

	if !ok {
		var username, password string

		if len(creds) > 0 {
			username = creds[0]
			if len(creds) > 1 {
				password = creds[1]
			}
		}

		// Timeout is the amount of time to wait for a server to respond when
		// first connecting and on follow up operations in the session. If
		// timeout is zero, the call may block forever waiting for a connection
		// to be established.

		// FailFast will cause connection and query attempts to fail faster when
		// the server is unavailable, instead of retrying until the configured
		// timeout period. Note that an unavailable server may silently drop
		// packets instead of rejecting them, in which case it's impossible to
		// distinguish it from a slow server, so the timeout stays relevant.

		info := mgo.DialInfo{
			Addrs:    []string{address},
			Timeout:  60 * time.Second,
			FailFast: true,
			Database: db_name,
			Username: username,
			Password: password,
		}

		var err error
		session, err = mgo.DialWithInfo(&info)
		if err != nil {
			panic(err)
		}

		//Save the Session for Later use.

		cached.Lock()
		cached.sessions[db_name] = session
		cached.Unlock()
	}

	//Return only a Session & the name. Let the Consumer make a Session.Copy()
	//to ensure that database state is resumed.

	return &MConn{session, db_name}
}
