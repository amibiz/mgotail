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
	Ts bson.MongoTimestamp "ts"
	H  int64               "h"
	V  int                 "v"
	Op string              "op"
	Ns string              "ns"
	O  bson.M              "o"
	O2 bson.M              "o2"
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
	done <-true
}
