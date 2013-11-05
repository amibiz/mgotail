package main

import (
	"bytes"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"time"
	"testing"
	"io"
)

func Print_log(buffer io.Writer, logs chan Oplog) {
	for log := range logs {
		fmt.Fprintf(buffer, "%s|%s|%s\n", log.Ns, log.Op, log.O["_id"].(bson.ObjectId).Hex())
	}
}
func Test_Tail(t *testing.T) {
	// Test the `Tail` on the Oplog
	fmt.Println("Testing `Tail`...")
	session, err := mgo.Dial(os.Getenv("MONGO_URL"))
	if err != nil {
		fmt.Fprintf(os.Stdout,"Cannot connect to Mongodb: %s.\n %s", os.Getenv("MONGO_URL"), err)
		t.Fail()
	}

	session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	
	var results bytes.Buffer
	var buffer bytes.Buffer

	db := session.DB("TailTest")
	coll := db.C("test")

	logs := make(chan Oplog)
	done := make(chan bool)
	now := bson.MongoTimestamp(time.Now().Unix() << 32)

	q := OplogQuery{session, bson.M{"ts": bson.M{"$gt": now}, "ns": "TailTest.test"}, time.Second * 3}
	go q.Tail(logs, done)
	go Print_log(&results, logs)
	for i := 0; i < 5; i++ {
		id := bson.NewObjectId()
		err = coll.Insert(bson.M{"name": "test_0", "_id": id})
		fmt.Fprintf(&buffer,"TailTest.test|i|%s\n", id.Hex())
	}
	<-done
	
	close(logs)

	resultsString := results.String()
	bufferString := buffer.String()
	if resultsString != bufferString {
		fmt.Fprintf(os.Stdout, "Got:\n %s\n\n Should have gotten: \n%s", resultsString, bufferString)
		t.Fail()
	}

	db.DropDatabase()
	fmt.Println("..done.\n")
}

func Test_PostJobs(t *testing.T) {
	// Test that PostJobs posts SFTP jobs on tail. Gearman not working yet.
	fmt.Println("Testing `PostJobs`...")
	finished := make(chan bool)
	go PostJobs(time.Second * 3, finished)
	
	session, err := mgo.Dial(os.Getenv("MONGO_URL"))
	session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	if err != nil {
		fmt.Fprintf(os.Stdout,"Cannot connect to Mongodb: %s\n %s", os.Getenv("MONGO_URL"), err)
		t.Fail()
	}

	db := session.DB(os.Getenv("MONGO_DB"))
	coll := db.C("systems")


	for i := 0; i < 5; i++ {
		id := bson.NewObjectId()
		err = coll.Insert(bson.M{"type": "sftp", "_id": id})
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		fmt.Printf("%s|i|%s.systems\n", id.Hex(), os.Getenv("MONGO_DB"))
	}

	<-finished
	
	fmt.Println("..done.\n")
}