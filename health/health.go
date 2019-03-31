package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var redissvc = os.Getenv("redissvc")

type healthreq struct {
	Code        string `json:"code"`
	SumInsured  string `json:"sumInsured"`
	DateOfBirth string `json:"dateOfBirth"`
}

func main() {
	log.Println("premium api starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health/premiums", premium)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func premium(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	h := marshallReq(string(body))
	type response struct {
		Premium string `json:"premium"`
	}
	s, err := calPremium(h)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		data, _ := json.Marshal(response{Premium: s})
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", data)
	}
}

func marshallReq(data string) healthreq {
	var h healthreq
	err := json.Unmarshal([]byte(data), &h)
	if err != nil {
		log.Println("there was an error in marshalling request", err.Error())
	}
	return h
}

func calPremium(h healthreq) (string, error) {
	log.Println(redissvc)
	c, err := redis.DialURL("redis://" + redissvc + ":6379/0")
	if err != nil {
		return "", fmt.Errorf("Cannot connect to redis %v", err)
	}
	defer c.Close()
	c.Do("SET", "2A", "3000")
	s, err := redis.String(c.Do("GET", h.Code))
	if err != nil {
		return "", fmt.Errorf("Cannot get premium for code %s from redis %v", h.Code, err)
	}
	return s, nil
}
