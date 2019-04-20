package main

import (
	"fmt"
	//"log"
	"strings"
	"testing"
)

const data string = `{"code":"1A","sumInsured":"100000","dateOfBirth":"1990-07-06"}`

func TestMarshal(t *testing.T) {
	h, err := marshallReq(data)
	if err != nil && h.Code != "2A" {
		t.Errorf("%s, wanted %q", h.Code, "2A")
	}
}

func TestCalcPremium(t *testing.T) {
	h, err := marshallReq(data)
	premium, err := calPremium(h)
	fmt.Println(premium)
	if err != nil {
		fmt.Println(err)
		if strings.Compare(premium, "1300") != 0 {
			t.Errorf("got %s, wanted %q", premium, "1300")
		}
	}
}
func TestLoad(t *testing.T) {
	err := load()
	if err != nil {
		t.Errorf("load failed")
	}
}

func TestCalAge(t *testing.T) {
	age := calculateAge("1990-07-06")
	fmt.Println(age)
}
