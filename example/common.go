package main

import (
	// _ "github.com/mattn/go-sqlite3" // commented out because of go.mod
	// or
	// _ "modernc.org/sqlite" // commented out because of go.mod

	"github.com/lanz-dev/go-sqlite"
)

func main() {
	db, err := sqlite.Connect(
		sqlite.WithPath("file:data.db"),
	)
	if err != nil {
		panic(err)
	}
	db.Exec("SELECT 1")
}
