package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/laher/toasties-galore/tpi"
)

var (
	chillybinAddr = tpi.Getenv("CHILLYBIN_ADDR", "http://localhost:7011")
	jafflrAddr    = tpi.Getenv("JAFFLR_ADDR", "http://localhost:7010")
)

func TestHappyPath(t *testing.T) {
	// reset environment
	resp, err := http.Get(fmt.Sprintf("%s/restock", chillybinAddr))
	if err != nil {
		t.Fatalf("error restocking chillybin: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error restocking chillybin (%s): %s", resp.Status, body)
	}
	m := getChillybinStats(t)
	if m["cheese"] != float64(10) {
		t.Fatalf("wrong amount of cheese after restocking: %v, %T", m, m["cheese"])
	}
	resp, err = http.Get(fmt.Sprintf("%s/toastie?i=cheese&i=vegemite", jafflrAddr))
	if err != nil {
		t.Fatalf("error fetching toastie: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error fetching toastie (%s): %s", resp.Status, body)
	}
	m = getChillybinStats(t)
	if m["cheese"] != float64(9) {
		t.Fatalf("wrong amount of cheese after grilling toastie: %v", m)
	}
}

func TestPickV2(t *testing.T) {
	p := []byte(`{"ingredient": "cheese", "quantity": 2, "customer": "gita"}`)
	// reset environment
	resp, err := http.Post(fmt.Sprintf("%s/v2/pick", chillybinAddr), "application/json", bytes.NewBuffer(p))
	if err != nil {
		t.Fatalf("error running /v2/pick: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error picking v2 (%s): %s", resp.Status, body)
	}
}

func getChillybinStats(t *testing.T) map[string]interface{} {
	resp, err := http.Get(fmt.Sprintf("%s/", chillybinAddr))
	if err != nil {
		t.Fatalf("error fetching status: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error fetching status (%s): %s", resp.Status, body)
	}

	var m = map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		t.Fatalf("error decoding body: %v", err)
	}
	t.Logf("Amounts: %v", m)
	return m
}
