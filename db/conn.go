package db

import (
	"errors"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/parthdesai/sc-notifications/utils/time"
)

var DbConnection *MConn

type M bson.M

func Convert(doc M, out interface{}) {
	stream, err := bson.Marshal(doc)
	if err == nil {
		bson.Unmarshal(stream, out)
	} else {
		panic(err)
	}
}

type MConn struct {
	db *mgo.Database
}

func (self *MConn) getCursor(table string, query M) *mgo.Query {

	fields, err1 := query["fields"].(M)
	delete(query, "fields")
	if !err1 {
		fields = M{}
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

	cursor := self.GetCursor(table, query)
	return cursor.Limit(limit).Skip(skip).Sort(sort).Select(fields)
}

type MapReduce mgo.MapReduce

func (self *MConn) MapReduce(table string, query M, result interface{}, job *MapReduce) (*mgo.MapReduceInfo, error) {
	coll := self.db.C(table)
	realJob := mgo.MapReduce{Map: job.Map, Reduce: job.Reduce, Finalize: job.Finalize, Scope: job.Scope, Verbose: true}
	return coll.Find(query).MapReduce(&realJob, result)
}

func (self *MConn) DropIndex(table string, key ...string) error {
	coll := self.db.C(table)
	return coll.DropIndex(key...)
}

func (self *MConn) DropIndices(table string) error {
	collection := self.db.C(table)
	indexes, err := collection.Indexes()
	if err == nil {
		for _, index := range indexes {
			err = collection.DropIndex(index.Key...)
			if err != nil {
				return err
			}
		}
	}

	if err != nil {
		panic(err)
	}
	return nil
}

func (self *MConn) findAndApply(table string, query M, change mgo.Change, result interface{}) error {
	change.ReturnNew = true

	coll := self.db.C(table)
	_, err := coll.Find(query).Apply(change, result)
	if err != nil {
		log.Println("Error Applying Changes", table, err)
	}
	return err
}

func (self *MConn) FindAndUpsert(table string, query M, doc M, result interface{}) error {
	change := mgo.Change{
		Update: doc,
		Upsert: true,
	}
	return self.findAndApply(table, query, change, result)
}

func (self *MConn) FindAndUpdate(table string, query M, doc M, result interface{}) error {
	change := mgo.Change{
		Update: doc,
		Upsert: false,
	}
	return self.findAndApply(table, query, change, result)
}

func (self *MConn) EnsureIndex(table string, index mgo.Index) error {
	coll := self.db.C(table)
	return coll.EnsureIndex(index)
}

func (self *MConn) GetCursor(table string, query M) *mgo.Query {
	coll := self.db.C(table)
	out := coll.Find(query)

	//explanation := M{}
	//err := out.Explain(explanation)
	//if err == nil {
	//    log.Println(table, query, explanation)
	//}

	return out
}

func (self *MConn) Get(table string, query M) *mgo.Iter {
	return self.getCursor(table, query).Iter()
}

func (self *MConn) HintedGetOne(table string, query M, result interface{}, hint string) error {
	cursor := self.getCursor(table, query).Hint(hint)
	err := cursor.One(result)
	if err != nil {
		log.Println("Error fetching", table, err)
	}
	return err
}

func (self *MConn) GetOne(table string, query M, result interface{}) error {
	cursor := self.getCursor(table, query)
	err := cursor.One(result)
	if err != nil {
		//log.Println("Error fetching", table, err)
	}
	return err
}

func (self *MConn) InternalConn() *mgo.Database {
	return self.db
}

func (self *MConn) HintedCount(table string, query M, hint string) int {
	cursor := self.getCursor(table, query).Select(M{"_id": 1}).Hint(hint)
	count, err := cursor.Count()
	if err != nil {
		log.Println("Error Counting", table, err)
	}
	return count
}

func (self *MConn) Count(table string, query M) int {
	cursor := self.getCursor(table, query).Select(M{"_id": 1})
	count, err := cursor.Count()
	if err != nil {
		log.Println("Error Counting", table, err)
	}
	return count
}

func (self *MConn) Upsert(table string, query M, doc M) error {

	var err error
	if len(doc) == 0 {
		err = errors.New(
			"Empty upsert is blocked. Refer to " +
				"https://github.com/Simversity/blackjack/issues/1051",
		)
	} else {
		coll := self.db.C(table)
		_, err = coll.Upsert(query, doc)
	}

	if err != nil {
		log.Println("Error Upserting:", table, err)
	}
	return err
}

func AlterDoc(doc *M, operator string, operation M) {
	spec := *doc
	if spec[operator] != nil {
		op, _ := spec[operator].(M)
		for key, value := range op {
			operation[key] = value
		}
	}
	spec[operator] = operation
}

func (self *MConn) Update(table string, query M, doc M) error {

	coll := self.db.C(table)
	var update_err error
	if len(doc) == 0 {
		update_err = errors.New(
			"Empty Update is blocked. Refer to " +
				"https://github.com/Simversity/blackjack/issues/1051",
		)
	} else {
		AlterDoc(&doc, "$set", M{"updated_on": time.EpochNow()})
		_, update_err = coll.UpdateAll(query, doc)
	}

	if update_err != nil {
		log.Println("Error Updating:", table, update_err)
	}
	return update_err
}

func (self *MConn) Delete(table string, query M) error {

	var delete_err error

	coll := self.db.C(table)

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

func (self *MConn) Insert(table string, arguments ...interface{}) (_id string) {

	var out interface{}
	if len(arguments) > 1 {
		out = arguments[1]
	} else {
		out = nil
	}

	doc := arguments[0]

	coll := self.db.C(table)
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

func (self *MConn) InsertMulti(table string, docs ...M) error {
	var interfaceSlice []interface{} = make([]interface{}, len(docs))
	for i, d := range docs {

		var ok bool
		if _, ok = d["_id"]; !ok {
			d["_id"] = bson.NewObjectId()
		}

		if _, ok = d["created_on"]; !ok {
			d["created_on"] = time.EpochNow()
		}

		interfaceSlice[i] = d
	}
	coll := self.db.C(table)
	err := coll.Insert(interfaceSlice...)
	if err != nil {
		log.Println("Error Multi Inserting:", table, err)
	}
	return err
}

func (self *MConn) Aggregate(table string, doc []M) *mgo.Pipe {
	coll := self.db.C(table)
	return coll.Pipe(doc)
}

var cachedConnections = map[string]*mgo.Session{}

func GetConn(db_name string, address string, creds ...string) *MConn {
	session := cachedConnections[db_name]
	if session == nil {
		var username, password string
		if len(creds) > 0 {
			username = creds[0]
			if len(creds) > 1 {
				username = creds[1]
			}
		}

		info := mgo.DialInfo{
			Addrs:    []string{address},
			Database: db_name,
			Direct:   true,
			Username: username,
			Password: password,
		}

		session, err := mgo.DialWithInfo(&info)
		if err != nil {
			panic(err)
		}

		cachedConnections[db_name] = session
		return &MConn{session.DB(db_name)}
	}
	return &MConn{session.DB(db_name)}
}
