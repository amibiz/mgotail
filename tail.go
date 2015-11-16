package mgotail

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OplogQuery describes a query you'd like to perform on the oplog.
type OplogQuery struct {
	Session *mgo.Session
	Query   bson.M
	Timeout time.Duration
}

// Oplog is a deserialization of the fields present in an oplog entry.
type Oplog struct {
	Timestamp    bson.MongoTimestamp `bson:"ts"`
	HistoryID    int64               `bson:"h"`
	MongoVersion int                 `bson:"v"`
	Operation    string              `bson:"op"`
	Namespace    string              `bson:"ns"`
	Object       bson.M              `bson:"o"`
	QueryObject  bson.M              `bson:"o2"`
}

// LastTime gets the timestamp of the last operation in the oplog.
// The return value can be used for construting queries on the "ts" oplog field.
func LastTime(session *mgo.Session) bson.MongoTimestamp {
	var member Oplog
	session.DB("local").C("oplog.rs").Find(nil).Sort("-$natural").One(&member)
	return member.Timestamp
}

// Tail starts the tailing of the oplog.
// Entries matching the configured query are published to channel passed to the function.
func (query *OplogQuery) Tail(logs chan Oplog, done chan bool) {
	db := query.Session.DB("local")
	collection := db.C("oplog.rs")
	iter := collection.Find(query.Query).LogReplay().Tail(query.Timeout)
	var log Oplog
	for iter.Next(&log) {
		logs <- log
	}
	iter.Close()
	done <- true
}
