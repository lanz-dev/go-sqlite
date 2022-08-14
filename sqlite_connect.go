package sqlite

import (
	"database/sql"
	"fmt"
	_ "unsafe" // For go:linkname
)

func buildDSN(path string, params []string) string {
	dsn := path
	if len(params) == 0 {
		return dsn
	}

	dsn += "?" + params[0]
	for _, p := range params[1:] {
		dsn += "&" + p
	}

	return dsn
}

// buildMattnDSN works with "github.com/mattn/go-sqlite3"
func buildMattnDSN(config *Config) string {
	var params []string

	if config.AutoVacuumMode != AutoVacuumDefault {
		mode := config.AutoVacuumMode.Int()
		if mode != -99 {
			params = append(params, fmt.Sprintf("_auto_vacuum=%d", mode))
		}
	}
	if config.BusyTimeout > 0 {
		params = append(params, fmt.Sprintf("_timeout=%d", config.BusyTimeout))
	}
	if config.CaseSensitiveLike {
		params = append(params, "_case_sensitive_like=true")
	}
	if config.ForeignKey {
		params = append(params, "_fk=true")
		if config.DeferForeignKeys {
			params = append(params, "defer_foreign_keys=true")
		}
	}
	if config.JournalMode != JournalDefault {
		params = append(params, fmt.Sprintf("_journal=%s", config.JournalMode))
	}
	if config.SyncMode != SyncDefault {
		mode := config.SyncMode.Int()
		if mode != -99 {
			params = append(params, fmt.Sprintf("_sync=%d", mode))
		}
	}

	return buildDSN(config.Path, params)
}

// buildModerncDSN works with "modernc.org/sqlite"
func buildModerncDSN(config *Config) string {
	var params []string

	if config.AutoVacuumMode != AutoVacuumDefault {
		params = append(params, fmt.Sprintf("_pragma=auto_vacuum(%s)", config.AutoVacuumMode))
	}
	if config.BusyTimeout > 0 {
		params = append(params, fmt.Sprintf("_pragma=busy_timeout(%d)", config.BusyTimeout))
	}
	if config.CaseSensitiveLike {
		params = append(params, "_pragma=case_sensitive_like(1)")
	}
	if config.ForeignKey {
		params = append(params, "_pragma=foreign_keys(1)")
		if config.DeferForeignKeys {
			params = append(params, "_pragma=defer_foreign_keys(1)")
		}
	}
	if config.JournalMode != JournalDefault {
		params = append(params, fmt.Sprintf("_pragma=journal_mode(%s)", config.JournalMode))
	}
	if config.SyncMode != SyncDefault {
		params = append(params, fmt.Sprintf("_pragma=synchronous(%s)", config.SyncMode))
	}

	return buildDSN(config.Path, params)
}

func detectDriver() (Driver, error) {
	for _, d := range sql.Drivers() {
		if d == DriverNameModernc {
			return DriverModernc, nil
		}
		if d == DriverNameMattn {
			return DriverMattn, nil
		}
	}
	return "", ErrCantDetectDriver
}

type sqlOpenFunc func(driverName, dataSourceName string) (*sql.DB, error)

func connect(openFunc sqlOpenFunc, opts ...Option) (*sql.DB, error) {
	config, err := buildConfig(opts...)
	if err != nil {
		return nil, err
	}

	db, err := openFunc(config.DriverName, config.DSN)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	if config.LimitConnection {
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(1)
		db.SetConnMaxLifetime(0)
		db.SetConnMaxIdleTime(0)
	}

	if config.JournalSizeLimit > 0 {
		_, err = db.Exec(fmt.Sprintf("PRAGMA journal_size_limit = %d;", config.JournalSizeLimit))
		if err != nil {
			return nil, err
		}
	}

	return db, err
}
