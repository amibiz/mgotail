package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// OplogQuery describes a query you'd like to perform on the oplog.
type OplogQuery struct {
	Session *mgo.Session
	Query   bson.M
	Timeout time.Duration
}

// Oplog is a deserialization of the fields present in an oplog entry.
type Oplog struct {
	Timestamp    bson.MongoTimestamp "ts"
	HistoryID    int64               "h"
	MongoVersion int                 "v"
	Operation    string              "op"
	Namespace    string              "ns"
	Object       bson.M              "o"
	QueryObject  bson.M              "o2"
}

// Helper function to get the timestamp of the last operation in the oplog.
// The return value can be used for construting queries on the "ts" oplog field.
func LastTime(session *mgo.Session) bson.MongoTimestamp {
	var member Oplog
	session.DB("local").C("oplog.rs").Find(nil).Sort("-$natural").One(&member)
	return member.Timestamp
}

// Begin tailing for oplog entries. Publishes Oplog objects to a channel.
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
