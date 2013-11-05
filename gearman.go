package main

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"time"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
)

type Job struct {
	Job_id string "job_id"
	Time_created time.Time "time_created"
	Time_updated time.Time "time_updaded"
	Time_finished time.Time "time_finished"
	State string "state"
}

type Gearman struct {
	Protocol string
	Host string
	Port string
	User string
	Password string
}

func (gm Gearman) GetUrl(worker string) string {
	protocol := gm.Protocol
	if protocol == "" {
		protocol = "http"
	}
	base := gm.Host
	if gm.User != "" {
		base = fmt.Sprintf("%s:%s@%s", gm.User, gm.Password, base)
	}
	if gm.Port != "" {
		base = fmt.Sprintf("%s:%s", base, gm.Port)
	}
	return fmt.Sprintf("%s://%s/job/%s", protocol, base, worker)
}
func (gm Gearman) SendJob(worker string, job_data map[string]string) (string, error) {
	gmURL := gm.GetUrl(worker)
	log.Println("Posting ", job_data, " to ", worker, " worker")
	form := url.Values{}
	for key,val := range job_data {
		form.Add(key, val)
	}
	resp, err := http.PostForm(gmURL, form)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var job Job
	err = bson.Unmarshal(body, &job)
	return job.Job_id, err
}
