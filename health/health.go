package main

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var redissvc = os.Getenv("redissvc")

type healthreq struct {
	Code        string `json:"code"`
	SumInsured  string `json:"sumInsured"`
	DateOfBirth string `json:"dateOfBirth"`
}

type response struct {
	Premium string `json:"premium"`
}

func main() {
	log.Println("premium api starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/healths/premiums", premium)
	mux.HandleFunc("/api/v1/healths/premiums/loads", loadMatrix)
	mux.HandleFunc("/api/v1/healths/premiums/unloads", unloadMatrix)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func premium(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	h, err := marshallReq(string(body))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	s, err := calPremium(h)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		data, _ := json.Marshal(response{Premium: s})
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", data)
	}
}

func loadMatrix(w http.ResponseWriter, req *http.Request) {
	if err := load(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func unloadMatrix(w http.ResponseWriter, req *http.Request) {
	if err := unload(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func marshallReq(data string) (*healthreq, error) {
	var h healthreq
	err := json.Unmarshal([]byte(data), &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func calulateScore(age int) (int, error) {
	if age >= 18 && age <= 35 {
		return 1, nil
	} else if age >= 36 && age <= 45 {
		return 2, nil
	} else if age >= 46 && age <= 55 {
		return 3, nil
	} else if age >= 56 && age <= 60 {
		return 4, nil
	} else if age >= 61 && age <= 65 {
		return 5, nil
	} else if age >= 66 && age <= 70 {
		return 6, nil
	} else if age > 70 {
		return 7, nil
	}
	return 0, nil
}

func calculateAge(bdate string) int {
	const layoutISO = "2006-01-02"
	dob, _ := time.Parse(layoutISO, bdate)
	now := time.Now()
	years := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		years--
	}
	return years
}

func calPremium(h *healthreq) (string, error) {
	c, err := conn()
	if err != nil {
		return "", err
	}
	defer c.Close()
	age := calculateAge(h.DateOfBirth)
	score, err := calulateScore(age)
	key := h.Code + ":" + h.SumInsured

	members, err := redis.Strings(c.Do("ZRANGEBYSCORE", key, score, score))

	if err != nil {
		return "", fmt.Errorf("Cannot get premium for code %s error %v", key, err)
	}
	if len(members) != 1 {
		return "", redis.ErrNil
	}
	return members[0], nil
}

func load() error {
	xlsx, err := excelize.OpenFile("./premium_tables.xlsx")
	if err != nil {
		fmt.Println(err)
		return err
	}

	c, err := conn()
	if err != nil {
		return err
	}
	defer c.Close()

	rows, _ := xlsx.GetRows("matrix")
	var score = 0
	for _, row := range rows {
		score++
		var key string
		var premium int
		for ci, cellv := range row {
			if ci == 0 {
				key = cellv
			}
			if ci == 1 {
				key = key + ":" + cellv
			}
			if ci == 3 {
				premium, _ = strconv.Atoi(cellv)
			}
		}
		fmt.Printf("key %v premium %v score %v ", key, premium, score)
		fmt.Println("")
		_, err := c.Do("ZADD", key, score, premium)
		if err != nil {
			return nil
		}
		if score == 8 {
			score = 0
		}
	}
	return nil
}

func unload() error {
	c, err := conn()
	if err != nil {
		return err
	}
	defer c.Close()
	_, errFlush := c.Do("FLUSHALL")
	if errFlush != nil {
		return err
	}
	return nil
}

func conn() (redis.Conn, error) {
	c, err := redis.DialURL("redis://" + redissvc + ":6379/0")
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to redis %v ", err)
	}
	return c, nil
}
