package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/laher/toasties-galore/tpi"
)

func main() {
	var (
		listenAddr = tpi.Getenv("ADDR", ":7001")
		dsn        = tpi.Getenv("DB_DSN", "postgres://root:secure@localhost:5432/postgres?sslmode=disable")
		version    = os.Getenv("VERSION")
		done       = make(chan bool)
	)
	db, err := connect(dsn)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	if err := runMigrationsSource(db); err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}
	log.Println("done migrations")

	h := &handler{db}
	router := routes(h, version)

	server := newServer(listenAddr, router)
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

func newServer(listenAddr string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
}

func connect(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func runMigrationsSource(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://resources/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}
	return nil
}
