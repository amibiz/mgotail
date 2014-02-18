package mgotail

import (
	"bytes"
	"fmt"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"testing"
	"time"
)

func printlog(buffer io.Writer, logs chan Oplog) {
	// Print logs from an oplog channel to a buffer
	for log := range logs {
		id := log.Object["_id"].(bson.ObjectId).Hex()
		fmt.Fprintf(buffer, "%s|%s|%s\n", log.Namespace, log.Operation, id)
	}
}
func Test_Tail(t *testing.T) {
	// Test the `Tail` on the Oplog
	fmt.Println("Testing `Tail`...")

	session, err := mgo.Dial(os.Getenv("MONGODB_PORT_27017_TCP_ADDR"))
	if err != nil {
		fmt.Printf("Cannot connect to Mongodb: %s.\n %s", os.Getenv("MONGO_URL"), err)
		t.Fail()
	}
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
		err = coll.Insert(bson.M{"name": "test_0", "_id": id})
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
	fmt.Println("..done.\n")
}
