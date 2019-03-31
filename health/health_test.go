package main

import (
	"strings"
	"testing"
)

const data string = `{"code":"2A","sumInsured":"100000","dateOfBirth":"07/06/1990"}`

func TestMarshal(t *testing.T) {
	h := marshallReq(data)
	if h.Code != "2A" {
		t.Errorf("%s, wanted %q", h.Code, "2A")
	}
}

func TestCalcPremium(t *testing.T) {
	premium, err := calPremium(marshallReq(data))
	//log.Println(premium)
	if err != nil {
		if strings.Compare(premium, "30000") != 0 {
			t.Errorf("got %s, wanted %q", "3000", premium)
		}
	}
}
