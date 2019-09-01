package main

import (
	"log"
	"net/http"
	"os"

	"github.com/laher/toasties-galore/tpi"
)

func main() {
	var (
		listenAddr    = os.Getenv("ADDR")
		done          = make(chan bool)
		chillybinAddr = os.Getenv("CHILLYBIN_ADDR")
	)
	if listenAddr == "" {
		listenAddr = ":7000"
	}
	if chillybinAddr == "" {
		chillybinAddr = "http://127.0.0.1:7001"
	}
	h := &handler{
		client: chillybinClient{chillybinAddr},
	}
	router := newRouter(h)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
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

func newRouter(h *handler) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("done"))
	})
	router.HandleFunc("/toastie", h.makeToastie)
	router.HandleFunc("/", h.status)
	return router
}
