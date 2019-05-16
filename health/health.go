package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

//Errcodes Error Codes
type Errcodes int

//Error Code Enum
const (
	SystemErr = iota
	InputJSONInvalid
	AgeRangeInvalid
	RiskDetailsInvalid
	InvalidRestMethod
	InvalidContentType
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

type erroresponse struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

func main() {
	/* file, err := os.OpenFile("premium.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(file)
		defer file.Close()
	} */
	log.Info("premium api starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", healthz)
	mux.HandleFunc("/api/v1/healths/premiums", premium)
	mux.HandleFunc("/api/v1/healths/premiums/loads", loadMatrix)
	mux.HandleFunc("/api/v1/healths/premiums/unloads", unloadMatrix)
	srv := http.Server{Addr: ":8000", Handler: mux}
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Info("shutting down health premium server...")
			srv.Shutdown(ctx)
			<-ctx.Done()
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	data := (time.Now()).String()
	log.Debug("health ok")
	w.Write([]byte(data))
}

func validateReq(w http.ResponseWriter, req *http.Request) (*healthreq, *erroresponse) {
	if req.Method != http.MethodPost {
		return nil, &erroresponse{Code: InvalidRestMethod, Message: fmt.Sprintf("Invalid method %s", req.Method)}
	}

	if req.Header.Get("Content-Type") != "application/json" {
		msg := fmt.Sprintf("Invalid content-type %s require %s", req.Header.Get("Content-Type"), "application/json")
		return nil, &erroresponse{Code: InvalidContentType, Message: msg}
	}

	body, _ := ioutil.ReadAll(req.Body)
	h, err := marshallReq(string(body))

	if err != nil {
		return nil, err
	}

	return h, nil
}

func premium(w http.ResponseWriter, req *http.Request) {
	h, err := validateReq(w, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(err)
		fmt.Fprintf(w, "%s", data)
	} else {
		premium, calErr := calPremium(h)
		if calErr != nil {
			if calErr.Code == SystemErr {
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				data, _ := json.Marshal(calErr)
				fmt.Fprintf(w, "%s", data)
			}
		} else {
			data, _ := json.Marshal(response{Premium: premium})
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s", data)
		}
	}
}

func loadMatrix(w http.ResponseWriter, req *http.Request) {
	if err := load(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func unloadMatrix(w http.ResponseWriter, req *http.Request) {
	if err := unload(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func marshallReq(data string) (*healthreq, *erroresponse) {
	var h healthreq
	err := json.Unmarshal([]byte(data), &h)
	if err != nil {
		log.Errorf("err %v during unmarshalling data %s ", err, data)
		return nil, &erroresponse{Code: SystemErr, Message: "input invalid"}
	}

	_, errDob := calculateAge(h.DateOfBirth)

	if errDob != nil {
		return nil, &erroresponse{Code: InputJSONInvalid, Message: "Invalid Date of birth enter for yyyy-mm-dd"}
	}

	if len(h.Code) == 0 || err != nil || len(h.SumInsured) == 0 {
		return nil, &erroresponse{Code: InputJSONInvalid, Message: "Invalid Input"}
	}
	return &h, nil
}

func calulateScore(age int) int {
	if age >= 18 && age <= 35 {
		return 1
	} else if age >= 36 && age <= 45 {
		return 2
	} else if age >= 46 && age <= 55 {
		return 3
	} else if age >= 56 && age <= 60 {
		return 4
	} else if age >= 61 && age <= 65 {
		return 5
	} else if age >= 66 && age <= 70 {
		return 6
	} else if age > 70 {
		return 7
	}
	return 0
}

func calculateAge(bdate string) (int, error) {
	const layoutISO = "2006-01-02"
	dob, err := time.Parse(layoutISO, bdate)
	if err != nil {
		return 0, err
	}
	now := time.Now()
	years := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		years--
	}
	return years, nil
}

func calPremium(h *healthreq) (string, *erroresponse) {
	c, err := connRead()
	if err != nil {
		log.Errorf(err.Error())
		return "", &erroresponse{Code: SystemErr, Message: "system err"}
	}
	defer c.Close()

	age, _ := calculateAge(h.DateOfBirth)
	if age > 70 {
		log.Errorf("age %v not in range of 18 to 70", age)
		msg := fmt.Sprintf("Age should be between 18 and 70")
		return "", &erroresponse{Code: AgeRangeInvalid, Message: msg}
	}
	score := calulateScore(age)
	key := h.Code + ":" + h.SumInsured

	members, err := redis.Strings(c.Do("ZRANGEBYSCORE", key, score, score))

	if err != nil {
		log.Errorf("Cannot get premium for code %s error %v", key, err)
		msg := fmt.Sprintf("Premium cannot be calculated risk details")
		return "", &erroresponse{Code: RiskDetailsInvalid, Message: msg}
	}
	if len(members) != 1 {
		log.Errorf("code %s dob %s sum assured %s combination not found ", h.Code, h.DateOfBirth, h.SumInsured)
		msg := fmt.Sprintf("Premium cannot be calculated for risk details")
		return "", &erroresponse{Code: RiskDetailsInvalid, Message: msg}
	}
	return members[0], nil
}

func load() error {
	xlsx, err := excelize.OpenFile("./premium_tables.xlsx")
	if err != nil {
		return fmt.Errorf("cannot load matrix file %v", err)
	}

	c, err := connWrite()
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
		log.Debugf("key %v premium %v score %v ", key, premium, score)
		_, err := c.Do("ZADD", key, score, premium)
		if err != nil {
			return fmt.Errorf("err adding key %v score %v premium %v to redis", key, score, premium)
		}
		if score == 8 {
			score = 0
		}
	}
	return nil
}

func unload() error {
	c, err := connWrite()
	if err != nil {
		return err
	}
	defer c.Close()
	_, errFlush := c.Do("FLUSHALL")
	if errFlush != nil {
		return fmt.Errorf("err flusing all keys %v", errFlush)
	}
	return nil
}

func connRead() (redis.Conn, error) {
	c, err := redis.DialURL("redis://" + redissvc + ":6379/0")
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to redis %v ", err)
	}
	return c, nil
}

func connWrite() (redis.Conn, error) {
	sc, err := redis.DialURL("redis://" + redissvc + ":26379/0")
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to redis sentinel %v ", err)
	}
	defer sc.Close()

	minfo, err := redis.Strings(sc.Do("sentinel", "get-master-addr-by-name", "redis-premium-master"))
	log.Println(minfo)
	if err != nil {
		return nil, fmt.Errorf("Cannot find redis master %v ", err)
	}

	mc, err := redis.DialURL("redis://" + minfo[0] + ":6379/0")
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to redis master %v ", err)
	}
	sc.Close()
	return mc, nil
}
