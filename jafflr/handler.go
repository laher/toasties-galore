package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type handler struct {
	client     chillybinClient
	toasting   bool
	statusLock sync.RWMutex
}

func (h *handler) status(w http.ResponseWriter, r *http.Request) {
	h.statusLock.RLock()
	defer h.statusLock.RUnlock()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("toasting: %v\n", h.toasting)))
}

func (h *handler) setStatus(val bool) {
	h.statusLock.Lock()
	defer func() {
		h.statusLock.Unlock()
	}()
	h.toasting = val
}

func (h *handler) makeToastie(w http.ResponseWriter, r *http.Request) {
	var (
		values      = r.URL.Query()
		ingredients = values["i"]
	)
	if err := validate(ingredients); err != nil {
		log.Printf("Error toasting toastie: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("input error - bad toastie"))
		return
	}
	ingredients = append(ingredients, "bread", "bread")
	for _, ingredient := range ingredients {
		if err := h.client.pick(ingredient, 1); err != nil {
			log.Printf("Error fetching ingredient: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ingredient error"))
			return
		}
	}
	h.cook(ingredients)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done"))
}

func (h *handler) cook(ingredients []string) {
	h.setStatus(true)
	defer h.setStatus(false)
	time.Sleep(time.Second * 10)
}

func validate(ingredients []string) error {
	if len(ingredients) < 1 {
		return errors.New("Not enough ingredients")
	}
	for _, c := range ingredients {
		if c == "cheese" {
			return nil
		}
	}
	return errors.New("No cheese on this toastie ")
}
