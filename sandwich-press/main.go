package main

import (
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
	return &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
}
