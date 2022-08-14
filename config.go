package sqlite

import (
	"fmt"
	"regexp"
)

// AutoVacuumMode for the SQLite connection.
//
// See https://www.sqlite.org/pragma.html#pragma_auto_vacuum.
type AutoVacuumMode string

// Int will return the int value for "github.com/mattn/go-sqlite3".
func (s *AutoVacuumMode) Int() int {
	switch *s {
	case AutoVacuumNone:
		return 0
	case AutoVacuumFull:
		return 1
	case AutoVacuumIncremental:
		return 2
	}
	return -99
}

// Driver used to connect to SQLite.
type Driver string

// JournalMode for the SQLite Connection.
//
// See https://www.sqlite.org/pragma.html#pragma_journal_mode.
type JournalMode string

// SyncMode for the SQLite Connection.
//
// See https://www.sqlite.org/pragma.html#pragma_synchronous.
type SyncMode string

// Int will return the int value for "github.com/mattn/go-sqlite3".
func (s *SyncMode) Int() int {
	switch *s {
	case SyncOff:
		return 0
	case SyncNormal:
		return 1
	case SyncFull:
		return 2
	case SyncExtra:
		return 3
	}
	return -99
}

// The different available AutoVaccum modes for SQLite.
//
// See https://www.sqlite.org/pragma.html#pragma_auto_vacuum.
const (
	AutoVacuumDefault     AutoVacuumMode = ""
	AutoVacuumNone        AutoVacuumMode = "NONE"
	AutoVacuumFull        AutoVacuumMode = "FULL"
	AutoVacuumIncremental AutoVacuumMode = "INCREMENTAL"
)

// Based on the Driver the connection DSN to SQLite is build.
const (
	DriverMattn   Driver = "github.com/mattn/go-sqlite3"
	DriverModernc Driver = "modernc.org/sqlite"
)

// These are the default [sql.Open] driverNames for the SQLite drivers.
const (
	DriverNameMattn   = "sqlite3"
	DriverNameModernc = "sqlite"
)

// The different available Journal modes for SQLite.
//
// See https://www.sqlite.org/pragma.html#pragma_journal_mode.
const (
	JournalDefault  JournalMode = ""
	JournalDelete   JournalMode = "DELETE"
	JournalTruncate JournalMode = "TRUNCATE"
	JournalPersist  JournalMode = "PERSIST"
	JournalMemory   JournalMode = "MEMORY"
	JournalWAL      JournalMode = "WAL"
	JournalOff      JournalMode = "OFF"
)

// The different available Sync modes for SQLite.
//
// See https://www.sqlite.org/pragma.html#pragma_synchronous.
const (
	SyncDefault SyncMode = ""
	SyncOff     SyncMode = "OFF"
	SyncNormal  SyncMode = "NORMAL"
	SyncFull    SyncMode = "FULL"
	SyncExtra   SyncMode = "EXTRA"
)

var regexPath = regexp.MustCompile(`^file:.+\..+$`)

// Config for the SQLite connection.
type Config struct {
	DSN             string // DSN string for [sql.Open]
	Driver          Driver // [sqlite.DriverMattn] or [sqlite.DriverModernc]
	DriverName      string // DriverName used in [sql.Open]
	Path            string // Path to the SQLite database
	LimitConnection bool   // Should we set the default limits?

	AutoVacuumMode    AutoVacuumMode // https://www.sqlite.org/pragma.html#pragma_auto_vacuum
	BusyTimeout       int            // https://www.sqlite.org/pragma.html#pragma_busy_timeout
	CaseSensitiveLike bool           // https://www.sqlite.org/pragma.html#pragma_case_sensitive_like
	DeferForeignKeys  bool           // https://www.sqlite.org/pragma.html#pragma_defer_foreign_keys
	ForeignKey        bool           // https://www.sqlite.org/pragma.html#pragma_foreign_keys
	JournalMode       JournalMode    // https://www.sqlite.org/pragma.html#pragma_journal_mode
	JournalSizeLimit  int            // https://www.sqlite.org/pragma.html#pragma_journal_size_limit
	SyncMode          SyncMode       // https://www.sqlite.org/pragma.html#pragma_synchronous
}

func newConfig() *Config {
	return &Config{
		LimitConnection: true,

		BusyTimeout:      4000,
		ForeignKey:       true,
		JournalMode:      JournalWAL,
		JournalSizeLimit: 100000000,
		SyncMode:         SyncNormal,
	}
}

func validatePath(path string) error {
	if path == ":memory" {
		return nil
	}
	if regexPath.MatchString(path) {
		return nil
	}
	return fmt.Errorf("given '%s', %w", path, ErrInvalidPath)
}

func buildConfig(opts ...Option) (*Config, error) {
	config := newConfig()
	for _, opt := range opts {
		opt(config)
	}

	if config.Path == "" {
		config.Path = ":memory"
	}
	if err := validatePath(config.Path); err != nil {
		return nil, err
	}

	if config.Driver == "" {
		var err error
		if config.Driver, err = detectDriver(); err != nil {
			return nil, err
		}
	}

	if config.DriverName == "" {
		if config.Driver == DriverModernc {
			config.DriverName = DriverNameModernc
		} else if config.Driver == DriverMattn {
			config.DriverName = DriverNameMattn
		}
	}

	if config.Driver == DriverModernc {
		config.DSN = buildModerncDSN(config)
	} else if config.Driver == DriverMattn {
		config.DSN = buildMattnDSN(config)
	}

	return config, nil
}
