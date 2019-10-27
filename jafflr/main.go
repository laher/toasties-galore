package main

import (
	"log"
	"net/http"

	"github.com/laher/toasties-galore/tpi"
)

func main() {
	var (
		listenAddr    = tpi.Getenv("ADDR", ":7000")
		chillybinAddr = tpi.Getenv("CHILLYBIN_ADDR", "http://127.0.0.1:7001")
		h             = &handler{
			client: chillybinClient{chillybinAddr}, // <- teh mechanical arm // HL
		}
		server = &http.Server{
			Addr:    listenAddr,
			Handler: tpi.Middleware(newRouter(h)), // <- middleware for observability // HL
		}
	)
	go tpi.GracefulShutdown(server) // <- for downtimeless deploy // HL
	log.Println("Server is about to listen for requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}
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
