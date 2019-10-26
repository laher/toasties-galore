package tpi

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	version = Getenv("VERSION", "v0.?")
)

func Middleware(in http.Handler) http.Handler {
	//	return TracingMiddleware(LoggingMiddleware(in))
	return LoggingMiddleware(in)
}

type responseStatusRecorder struct {
	http.ResponseWriter // embedded member already satisfies the interface
	status              int
	size                int
}

func (w *responseStatusRecorder) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (rw *responseStatusRecorder) Write(b []byte) (int, error) {
	if !rw.written() {
		// always trigger this
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseStatusRecorder) written() bool {
	return rw.status != 0
}

func someTracingObserve(duration time.Duration, labelsMap map[string]string) {
	// TODO insert statsd/prometheus/honeycomb thing here
}

func TracingMiddleware(in http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now().UTC()
		wsr := &responseStatusRecorder{ResponseWriter: w}
		in.ServeHTTP(wsr, r) // the wrapped handler is still invoked here // HL
		labelsMap := map[string]string{
			"method":  r.Method,
			"status":  strconv.Itoa(wsr.status),
			"version": version, // <- we can use this to compare error rates // HL
		}
		someTracingObserve(time.Since(t), labelsMap) // maybe honeycomb/prometheus/statsd ... // HL
	})
}

func LoggingMiddleware(in http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now().UTC()
		log.Printf("[%s] request: [%s:%v]", version, r.Method, r.URL)
		wsr := &responseStatusRecorder{ResponseWriter: w}
		in.ServeHTTP(wsr, r)
		log.Printf("[%s] response: [%d, %db in %s] for [%s:%v]", version, wsr.status, wsr.size, time.Since(t), r.Method, r.URL)
	})
}

func GracefulShutdown(server *http.Server) {
	osInterrupt := make(chan os.Signal)
	signal.Notify(osInterrupt, os.Interrupt)
	<-osInterrupt // subscribe to Ctrl+C or kill/SIGTERM // HL
	log.Println("Received OS interrupt - stop listening but avoid interrupting existing connections...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil { // keeps existing connections open until ctx.Done() // HL
		log.Fatalf("Server shutdown error: %v", err)
	}
}
