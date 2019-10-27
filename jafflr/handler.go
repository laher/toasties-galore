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
		doneness    = values.Get("doneness")
		customer    = values.Get("customer")
	)
	if err := validate(ingredients); err != nil {
		log.Printf("Error toasting toastie: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("input error - bad toastie"))
		return
	}
	ingredients = append(ingredients, "bread", "bread")
	for _, ingredient := range ingredients {
		if isChillybinV2(customer) {
			if err := h.client.pickV2(ingredient, 1, customer); err != nil {
				log.Printf("Error fetching ingredient: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Ingredient error"))
				return
			}
		} else {
			if err := h.client.pick(ingredient, 1); err != nil {
				log.Printf("Error fetching ingredient: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Ingredient error"))
				return
			}
		}
	}
	if err := h.cook(ingredients, doneness); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Doneness error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done"))
}

func (h *handler) cook(ingredients []string, doneness string) error {
	h.setStatus(true)
	defer h.setStatus(false)
	var duration int
	switch doneness {
	case "light", "": // light is default
		duration = 1000
	case "medium":
		duration = 2000
	case "well-done":
		duration = 5000
	case "burnt":
		duration = 10000
	default:
		return errors.New("Invalid doneness")
	}
	time.Sleep(time.Millisecond * time.Duration(duration))
	return nil
}

var (
	notEnoughIngredients = errors.New("Not enough ingredients")
	noCheese             = errors.New("No cheese on this toastie")
)

func validate(ingredients []string) error {
	if len(ingredients) < 1 {
		return notEnoughIngredients
	}
	for _, c := range ingredients {
		if c == "cheese" {
			return nil
		}
	}
	return noCheese
}
