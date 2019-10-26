package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/laher/toasties-galore/tpi"
)

var (
	chillybinAddr = tpi.Getenv("CHILLYBIN_ADDR", "http://localhost:7011")
	jafflrAddr    = tpi.Getenv("JAFFLR_ADDR", "http://localhost:7010")
)

func TestConnectivity(t *testing.T) {
	d := 10
	for i := 0; i < d; i++ {
		_, err := http.Get(fmt.Sprintf("%s/", chillybinAddr))
		if err != nil {
			t.Logf("could not connect - wait 2s: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		t.Logf("connected - continue")
		return
	}
	t.Errorf("Could NOT connect after %d attempts. Fail", d)
}

func TestHappyPath(t *testing.T) {
	// reset environment // HL
	resp, err := http.Post(fmt.Sprintf("%s/restock", chillybinAddr), "text/plain", nil)
	if err != nil {
		t.Fatalf("error restocking chillybin: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error restocking chillybin (%s): %s", resp.Status, body)
	}
	m := getChillybinStats(t)
	if m["cheese"] != float64(10) { // assert initial state // HL
		t.Fatalf("wrong amount of cheese after restocking: %v, %T", m, m["cheese"])
	}
	resp, err = http.Post(fmt.Sprintf("%s/toastie?i=cheese&i=vegemite", jafflrAddr), "text/plain", nil)
	if err != nil {
		t.Fatalf("error fetching toastie: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error fetching toastie (%s): %s", resp.Status, body)
	}
	m = getChillybinStats(t)
	if m["cheese"] != float64(9) { // assert resulting state // HL
		t.Fatalf("wrong amount of cheese after grilling toastie: %v", m)
	}
}

func TestRestock(t *testing.T) {
	// reset environment
	resp, err := http.Post(fmt.Sprintf("%s/restock", chillybinAddr), "text/plain", nil)
	if err != nil {
		t.Fatalf("error restocking chillybin: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error restocking chillybin (%s): %s", resp.Status, body)
	}
}

func TestBurnt(t *testing.T) {
	doneness := os.Getenv("DONENESS")
	if doneness == "" {
		t.Skip()
	}
	url := fmt.Sprintf("%s/toastie?i=cheese&i=vegemite&doneness=%s", jafflrAddr, doneness)
	t.Logf("GETting %s", url)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("error fetching toastie: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error fetching toastie (%s): %s", resp.Status, body)
	}
}

func getChillybinStats(t *testing.T) map[string]interface{} {
	resp, err := http.Get(fmt.Sprintf("%s/", chillybinAddr))
	if err != nil {
		t.Fatalf("error fetching status: %v", err)
	} else if resp.StatusCode != http.StatusOK {
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
