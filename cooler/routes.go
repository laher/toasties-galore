package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func routes(db *sql.DB, version string) http.Handler {
	router := http.NewServeMux()

	// check health
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(version))
	})

	// check stock
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var (
			query       = "SELECT name, quantity FROM INGREDIENTS"
			err         error
			ingredients = map[string]int{}
		)
		rows, err := db.Query(query)
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
	})

	// pick ingredient by name and quantity
	router.HandleFunc("/pick", func(w http.ResponseWriter, r *http.Request) {
		var (
			query  = "UPDATE ingredients SET quantity= quantity - $2 WHERE name=$1 and quantity - $2 >= 0"
			values = r.URL.Query()
			n      = values.Get("name")
			q      = values.Get("quantity")
		)
		quantity, err := strconv.Atoi(q)
		if err != nil {
			log.Printf("Error fetching ingredient %s with quantity '%s': %v", n, q, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("input error - quantity must be integer"))
			return
		}
		res, err := db.Exec(query, n, quantity)
		if err != nil {
			log.Printf("Error fetching ingredient %s with quantity '%s': %v", n, q, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error fetching ingredient"))
			return
		}
		affected, err := res.RowsAffected()
		if err != nil {
			log.Printf("Error fetching ingredient %s with quantity '%s': %v", n, q, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error fetching ingredient"))
			return
		}
		if affected != 1 {
			log.Printf("Error fetching ingredient %s with quantity '%s': %v rows affected", n, q, affected)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error fetching ingredient"))
			return
		}
		w.Write([]byte("ok"))
	})

	// restock all ingredients
	router.HandleFunc("/restock", func(w http.ResponseWriter, r *http.Request) {
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
			if _, err = db.Exec(query, ingredient, q); err != nil {
				log.Printf("Error restocking ingredient %s: %v", ingredient, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error restocking ingredient: " + ingredient))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ingredients restocked"))
	})

	return router
}
