package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/lib/pq"
)

type handler struct {
	db *sql.DB
}

func (h *handler) checkStock(w http.ResponseWriter, r *http.Request) {
	var (
		query       = "SELECT name, quantity FROM INGREDIENTS"
		err         error
		ingredients = map[string]int{}
	)
	rows, err := h.db.Query(query)
	if err != nil {
		log.Printf("Error fetching stock: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching stock"))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			n string
			q int
		)
		rows.Scan(&n, &q)
		ingredients[n] = q
	}
	b, err := json.Marshal(ingredients)
	if err != nil {
		log.Printf("Error serializing stock: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error serializing stock"))
		return
	}
	w.Write(b)
}

func (h *handler) pick(w http.ResponseWriter, r *http.Request) {
	var (
		query    = "UPDATE ingredients SET quantity= quantity - $2 WHERE name=$1 and quantity - $2 >= 0"
		values   = r.URL.Query()
		n        = values.Get("name")
		q        = values.Get("quantity")
		quantity = 0
		err      error
	)
	if quantity, err = strconv.Atoi(q); err != nil {
		log.Printf("Error fetching ingredient %s with quantity '%s': %v", n, q, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("input error - quantity must be integer"))
		return
	}
	res, err := h.db.Exec(query, n, quantity)
	if err != nil {
		switch e := err.(type) {
		case *pq.Error:
			log.Printf("pq error: %s", e.Code.Name())
		}
		log.Printf("Error fetching ingredient %s with quantity '%s': %v", n, q, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching ingredient"))
		return
	}
	affected, err := res.RowsAffected()
	if err != nil || affected > 1 {
		log.Printf("Error fetching ingredient %s with quantity '%s': %v", n, q, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching ingredient"))
		return
	}
	if affected < 1 {
		log.Printf("Error fetching ingredient %s with quantity '%s': 0 rows affected", n, q)
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("could not find requested ingredient"))
		return
	}
	w.Write([]byte("ok"))
}

// restock all ingredients
func (h *handler) restock(w http.ResponseWriter, r *http.Request) {
	var (
		query       = "INSERT INTO INGREDIENTS (name, quantity) VALUES ($1,$2) ON CONFLICT(name) DO UPDATE SET quantity=$2"
		err         error
		ingredients = map[string]int{
			"cheese":   10,
			"bread":    10,
			"vegemite": 10,
		}
	)
	for ingredient, q := range ingredients {
		if _, err = h.db.Exec(query, ingredient, q); err != nil {
			log.Printf("Error restocking ingredient %s: %v", ingredient, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error restocking ingredient: " + ingredient))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ingredients restocked"))
}

func routes(h *handler, version string) http.Handler {
	router := http.NewServeMux()

	// check health
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(version))
	})

	// check stock
	router.HandleFunc("/", h.checkStock)

	// pick ingredient by name and quantity
	router.HandleFunc("/pick", h.pick)

	// restock all ingredients
	router.HandleFunc("/restock", h.restock)

	return router
}
