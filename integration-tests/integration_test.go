package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/laher/toasties-galore/tpi"
)

var (
	chillybinAddr  = tpi.Getenv("CHILLYBIN_ADDR", "http://localhost:7011")
	jafflotronAddr = tpi.Getenv("JAFFLOTRON_ADDR", "http://localhost:7010")
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

	resp, err = http.Get(fmt.Sprintf("%s/", chillybinAddr))
	if err != nil {
		t.Fatalf("error fetching status: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error fetching status (%s): %s", resp.Status, body)
	}

	resp, err = http.Get(fmt.Sprintf("%s/toastie?i=cheese&i=vegemite", jafflotronAddr))
	if err != nil {
		t.Fatalf("error fetching toastie: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("error fetching toastie (%s): %s", resp.Status, body)
	}
}
