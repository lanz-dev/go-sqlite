package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

// ErrCantDetectDriver will be returned if no specific [sqlite.Driver] was set per [sqlite.WithDriver] and
// a registered driver could not be found automatically.
//
// Did you import your driver?
//
//	  import (
//		  _ "github.com/mattn/go-sqlite3"
//		  // or
//		  _ "modernc.org/sqlite"
//	  )
var ErrCantDetectDriver = errors.New("cannot detect the correct driver")

// ErrInvalidPath will be returned if the provided path is invalid.
var ErrInvalidPath = errors.New("invalid path provided")

var openFunc sqlOpenFunc = sql.Open

// Connect will connect to a SQLite database with some typical performance settings and foreign key support.
//
// You should call Shutdown / ShutdownContext on defer or within a shutdown hook.
// You should call after some hours Optimize / OptimizeContext on a regular basis.
//
// If no [sqlite.Driver] is manually set per [sqlite.WithDriver], [sqlite.Connect] tries to detect a [sqlite.Driver].
// The detected [sqlite.Driver] is based on the registered sql drivers:
//   - "sqlite" for "modernc.org/sqlite"
//   - "sqlite3" for "github.com/mattn/go-sqlite3"
//
// These are the default Settings]:
//   - Path is ":memory" for an in-memory sqlite connection
//   - JournalMode WAL
//   - SyncMode NORMAL
//   - Sets on sql.DB the connection limit to 1 and lifetime to 0
//   - Foreign key support enabled
//   - BusyTimeout set to 4000
//   - JournalSizeLimit set to 100000000
func Connect(opts ...Option) (*sql.DB, error) {
	return connect(openFunc, opts...)
}

// Optimize will run `PRAGMA optimize` which should be run on connection close and every few hours.
//
// See https://www.sqlite.org/pragma.html#pragma_optimize
func Optimize(db *sql.DB) error {
	return OptimizeContext(context.Background(), db)
}

// OptimizeContext will run `PRAGMA optimize` which should be run on connection close and every few hours.
//
// See https://www.sqlite.org/pragma.html#pragma_optimize
func OptimizeContext(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, "PRAGMA optimize;"); err != nil {
		return err
	}
	return nil
}

// Shutdown should be called before the application exits.
func Shutdown(db *sql.DB) error {
	return ShutdownContext(context.Background(), db)
}

// ShutdownContext should be called before the application exits.
func ShutdownContext(ctx context.Context, db *sql.DB) error {
	if err := OptimizeContext(ctx, db); err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

// Vacuum will call the VACUUM statement.
//
// This is usefully if you don't have auto-vacuum enabled and want to shrink the sqlite db file.
//
// See https://www.sqlite.org/lang_vacuum.html.
func Vacuum(db *sql.DB) error {
	return VacuumContext(context.Background(), db)
}

// VacuumContext will call the VACUUM statement.
//
// This is usefully if you don't have auto-vacuum enabled and want to shrink the sqlite db file.
//
// See https://www.sqlite.org/lang_vacuum.html.
func VacuumContext(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, "VACUUM;"); err != nil {
		return err
	}
	return nil
}
