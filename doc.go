/*
Package sqlite is a collection of tools to work with a SQLite database in go.

# Example

This is a typical usage example.

	package main

	import (
		_ "github.com/mattn/go-sqlite3"
		// or
		_ "modernc.org/sqlite"

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

# Changelog

v0.1.0
  - Initial Release
*/
package sqlite
