// mgotail is a simple package to tail mongodb oplogs: http://docs.mongodb.org/manual/core/replica-set-oplog/.
//
//
// Usage
//
// Here's an example program that tails all operations starting from the time the program is launched:
//
//         package main
//
//         import (
//         	"fmt"
//         	"labix.org/v2/mgo"
//         	"labix.org/v2/mgo/bson"
//         	"os"
//         	"github.com/Clever/mgotail"
//         )
//
//         func main() {
//         	session, err := mgo.Dial(os.Getenv("MONGO_URL"))
//         	if err != nil {
//         		panic(err)
//         	}
//
//         	q := mgotail.OplogQuery{
//         		Session: session,
//         		Query:   bson.M{"ts": bson.M{"$gt": mgotail.LastOplogTime(session)}},
//         		Timeout: -1, // See http://godoc.org/labix.org/v2/mgo#Query.Tail
//         	}
//
//         	logs := make(chan Oplog)
//         	done := make(chan bool)
//         	go q.Tail(logs, done)
//         	go func() {
//         		for log := range logs {
//         			fmt.Printf("%s|%s|%s\n", log.Namespace, log.Operation, log.Object)
//         		}
//         	}()
//         	<-done
//         }
//
package mgotail
