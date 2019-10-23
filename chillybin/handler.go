package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

type PickRequest struct {
	Ingredient string
	Quantity   int
	Customer   string
	// TODO some other fields here
}

func (h *handler) pickV2(w http.ResponseWriter, r *http.Request) {
	var (
		m = PickRequest{}
	)
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tx, err := h.db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching ingredient"))
		return
	}
	var (
		query = "INSERT INTO orders (customer, ingredient, quantity, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id"
		id    int64
	)

	if err := tx.QueryRow(query, m.Customer, m.Ingredient, m.Quantity).Scan(&id); err != nil {
		switch e := err.(type) {
		case *pq.Error:
			log.Printf("pq error: %s", e.Code.Name())
		}
		log.Printf("Error creating order: '%+v', %v", m, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error creating order"))
		return
	}
	if err := h.doPick(w, r, tx, m.Ingredient, m.Quantity); err != nil {
		tx.Rollback()
		log.Printf("Error picking ingredient: '%+v', %v", m, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching ingredient"))
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: '%+v', %v", m, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching ingredient"))
		return
	}
	b, err := json.Marshal(map[string]int64{"id": id})
	if err != nil {
		log.Printf("Error marshaling response with id: '%+v', %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error fetching ingredient"))
		return
	}
	w.Write(b)

}

type execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func (h *handler) doPick(w http.ResponseWriter, r *http.Request, db execer, n string, quantity int) error {
	var (
		query = "UPDATE ingredients SET quantity= quantity - $2 WHERE name=$1 and quantity - $2 >= 0"
	)
	res, err := db.Exec(query, n, quantity)
	if err != nil {
		switch e := err.(type) {
		case *pq.Error:
			log.Printf("pq error: %s", e.Code.Name())
		}
		log.Printf("Error fetching ingredient %s with quantity '%d': %v", n, quantity, err)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil || affected > 1 {
		log.Printf("Error fetching ingredient %s with quantity '%d': %v", n, quantity, err)
		return err
	}
	if affected < 1 {
		log.Printf("Error fetching ingredient %s with quantity '%d': 0 rows affected", n, quantity)
		return err
	}
	return nil
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
	router.HandleFunc("/v2/pick", h.pickV2) // Updated route

	// restock all ingredients
	router.HandleFunc("/restock", h.restock)

	return router
}
