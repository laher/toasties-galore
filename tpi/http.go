package tpi

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Middleware(in http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		log.Printf("request: %s to %v", req.Method, req.URL)
		in.ServeHTTP(resp, req)
		log.Printf("response: %s to %v", req.Method, req.URL)
	})
}

func GracefulShutdownOSInterrupt(server *http.Server) {
	osInterrupt := make(chan os.Signal)
	signal.Notify(osInterrupt, os.Interrupt)
	<-osInterrupt
	log.Println("Received OS interrupt - shutting down...")
	duration := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
