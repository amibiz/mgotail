package main

import (
	"fmt"
	"log"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
	"os"
)

func Now() bson.MongoTimestamp {
	return bson.MongoTimestamp(time.Now().Unix() << 32)
}

func SendSFTPJobs(gmlogs chan Oplog) {
	// Post the SFTP jobs that come throught the channel
	gmHost := os.Getenv("GEARMAN_HOST")
	gmPort := os.Getenv("GEARMAN_PORT")
	gmUser := os.Getenv("GEARMAN_USER")
	gmPw := os.Getenv("GEARMAN_PASSWORD")
	gm := Gearman{"https", gmHost, gmPort, gmUser, gmPw}

	for logdoc := range gmlogs {
		job_data := map[string]string{"system": logdoc.Object["_id"].(bson.ObjectId).Hex()}
		id, err := gm.SendJob("sftp", job_data)
		if err != nil {
			log.Println("Error posting job:", err)
		} else {
			log.Println("Posted job:", id)
		}
	}
}

func PostJobs(timeout time.Duration, finished chan bool) {
	// Tail the oplog for new SFTP systems and post the job when it comes through.
	mongo_url := os.Getenv("MONGO_URL")
	session,err := mgo.Dial(mongo_url)

	if err != nil {
		log.Println("Error: Could not connect to db", mongo_url)
	}

	ns := fmt.Sprintf("%s.systems", os.Getenv("MONGO_DB"))
	q := bson.M{"ts": bson.M{"$gt": Now()}, "o.type": "sftp", "ns": ns, "op": "i"}
	query := OplogQuery{session, q, timeout}

	done := make(chan bool)
	logs := make(chan Oplog)

	go query.Tail(logs, done)
	go SendSFTPJobs(logs)

	<-done
	close(logs)
	finished <-true
}
