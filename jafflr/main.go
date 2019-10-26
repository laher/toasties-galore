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
		done          = make(chan bool)
	)
	h := &handler{
		client: chillybinClient{chillybinAddr}, // <- oh, teh arm // HL
	}
	router := newRouter(h)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: tpi.Middleware(router),
	}
	go func() {
		tpi.GracefulShutdown(server)
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
