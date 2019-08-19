package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/laher/toasties-galore/tpi"
)

func main() {
	var (
		listenAddr = os.Getenv("ADDR")
		done       = make(chan bool)
	)
	if listenAddr == "" {
		listenAddr = ":7000"
	}
	server := newServer(listenAddr)
	go func() {
		tpi.GracefulShutdownOSInterrupt(server)
		close(done)
	}()
	log.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}
	<-done // await done
	log.Println("Shutdown complete - stop")
}

func newServer(listenAddr string) *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("done"))
	})
	router.HandleFunc("/toastit", func(w http.ResponseWriter, r *http.Request) {

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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("done"))
	})
	return &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
}

func validate(ingredients []string) error {
	if len(ingredients) < 3 {
		return errors.New("Not enough ingredients")
	}
	if ingredients[0] != "bread" || ingredients[len(ingredients)-1] != "bread" {
		return errors.New("Wot no bread")
	}
	if ingredients[1] != "cheese" {
		return errors.New("Cheese comes after first bread")
	}
	return nil
}
