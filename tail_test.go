package mgotail

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/facebookgo/mgotest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func printlog(buffer io.Writer, logs chan Oplog) {
	// Print logs from an oplog channel to a buffer
	for log := range logs {
		id := log.Object["_id"].(bson.ObjectId).Hex()
		fmt.Fprintf(buffer, "%s|%s|%s\n", log.Namespace, log.Operation, id)
	}
}

func Test_Tail(t *testing.T) {
	replset := mgotest.NewReplicaSet(2, t)
	session := replset.Session()
	session.EnsureSafe(&mgo.Safe{WMode: "majority"})

	var results bytes.Buffer
	var buffer bytes.Buffer

	logs := make(chan Oplog)
	done := make(chan bool)
	last := LastTime(session)

	q := OplogQuery{session, bson.M{"ts": bson.M{"$gt": last}, "ns": "TailTest.test"}, time.Second * 3}
	go q.Tail(logs, done)
	go printlog(&results, logs)

	db := session.DB("TailTest")
	coll := db.C("test")
	for i := 0; i < 5; i++ {
		id := bson.NewObjectId()
		if err := coll.Insert(bson.M{"name": "test_0", "_id": id}); err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&buffer, "TailTest.test|i|%s\n", id.Hex())
	}

	<-done
	close(logs)

	resultsString := results.String()
	bufferString := buffer.String()
	if resultsString != bufferString {
		fmt.Printf("Got:\n %s\n\n Should have gotten: \n%s", resultsString, bufferString)
		t.Fail()
	}

	db.DropDatabase()
}
