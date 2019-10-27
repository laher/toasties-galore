package main

import (
	"os"
	"testing"
)

func TestMigrate(t *testing.T) {

	if os.Getenv("DB_DSN") == "" {
		t.Skip("no db - skip")
	}
	db, err := connect(os.Getenv("DB_DSN"))
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
		t.FailNow()
	}
	err = runMigrationsSource(db)
	if err != nil {
		t.Errorf("error running migrations: %v", err)
	}

}
