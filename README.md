# go-sqlite

[![Go Reference](https://pkg.go.dev/badge/github.com/lanz-dev/go-sqlite.svg)](https://pkg.go.dev/github.com/lanz-dev/go-sqlite)
[![Coverage Status](https://coveralls.io/repos/github/lanz-dev/go-sqlite/badge.svg?branch=main)](https://coveralls.io/github/lanz-dev/go-sqlite?branch=main)
[![Github Action](https://github.com/lanz-dev/go-sqlite/actions/workflows/main.yml/badge.svg)](https://github.com/lanz-dev/go-sqlite/actions/workflows/main.yml)

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
