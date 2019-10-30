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

// This is here because I want to run this locally while my services run in docker-compose
// docker-compose can't easily tell us when services are 'up' - just that they've started.
// So, retry in this preliminary test, and keep tests clean
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
	// reset environment: // HL
	if resp, err := http.Post(fmt.Sprintf("%s/restock", chillybinAddr), "text/plain", nil); err != nil {
		t.Fatalf("error restocking chillybin: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error restocking chillybin (%s): %s", resp.Status, body)
	}
	if m := getChillybinStats(t); m["cheese"] != 10 { // assert initial state // HL
		t.Fatalf("wrong amount of cheese after restocking: %v, %T", m, m["cheese"])
	}
	// Method under test: // HL
	if resp, err := http.Post(fmt.Sprintf("%s/toastie?i=cheese&i=vegemite", jafflrAddr), "text/plain", nil); err != nil {
		t.Fatalf("error fetching toastie: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error fetching toastie (%s): %s", resp.Status, body)
	}
	if m := getChillybinStats(t); m["cheese"] != 9 { // assert resulting state // HL
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

func getChillybinStats(t *testing.T) map[string]int {
	resp, err := http.Get(fmt.Sprintf("%s/", chillybinAddr))
	if err != nil {
		t.Fatalf("error fetching status: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error fetching status (%s): %s", resp.Status, body)
	}

	var m = map[string]int{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		t.Fatalf("error decoding body: %v", err)
	}
	t.Logf("Amounts: %v", m)
	return m
}
