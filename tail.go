package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type OplogQuery struct {
	Session *mgo.Session
	Query   bson.M
	Timeout time.Duration
}

type Oplog struct {
	Timestamp    bson.MongoTimestamp "ts"
	HistoryId    int64               "h"
	MongoVersion int                 "v"
	Operation    string              "op"
	Namespace    string              "ns"
	Object       bson.M              "o"
	QueryObject  bson.M              "o2"
}

func LastTime(session *mgo.Session) bson.MongoTimestamp {
	// Get the Mongo Timestamp of the last member in the oplog
	var member Oplog
	session.DB("local").C("oplog.rs").Find(nil).Sort("-$natural").One(&member)
	return member.Timestamp
}

func (query *OplogQuery) Tail(logs chan Oplog, done chan bool) {
	// Add a tail to the oplog.rs collection and send any new logs to the Oplog channel.
	db := query.Session.DB("local")
	collection := db.C("oplog.rs")
	iter := collection.Find(query.Query).Tail(query.Timeout)
	var log Oplog
	for iter.Next(&log) {
		logs <- log
	}
	iter.Close()
	done <- true
}
