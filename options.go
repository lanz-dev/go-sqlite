package sqlite

// Option is a func to set configuration options for SQLite.
type Option func(c *Config)

// WithAutoVacuumMode will set the auto vacuum mode.
//
// Setting the value [sqlite.AutoVacuumDefault] will not set the pragma at all and uses the driver default behaviour.
//
// See https://www.sqlite.org/pragma.html#pragma_auto_vacuum.
func WithAutoVacuumMode(mode AutoVacuumMode) Option {
	return func(c *Config) {
		c.AutoVacuumMode = mode
	}
}

// WithBusyTimeout will set the busy timeout.
//
// Setting a value of 0 will not set the pragma at all and uses the driver default behaviour.
//
// See https://www.sqlite.org/pragma.html#pragma_busy_timeout.
func WithBusyTimeout(timout int) Option {
	return func(c *Config) {
		c.BusyTimeout = timout
	}
}

// WithCaseSensitiveLike will enable or disable the case-sensitive like.
//
// See https://www.sqlite.org/pragma.html#pragma_case_sensitive_like.
func WithCaseSensitiveLike(enabled bool) Option {
	return func(c *Config) {
		c.CaseSensitiveLike = enabled
	}
}

// WithDeferredForeignKeys will enable or disable deferred foreign keys.
//
// See https://www.sqlite.org/pragma.html#pragma_defer_foreign_keys.
func WithDeferredForeignKeys(enabled bool) Option {
	return func(c *Config) {
		c.DeferForeignKeys = enabled
	}
}

// WithDisabledLimits will disable the default limits for the connection.
//
// These are the default limits:
//   - db.SetMaxIdleConns(1)
//   - db.SetMaxOpenConns(1)
//   - db.SetConnMaxLifetime(0)
//   - db.SetConnMaxIdleTime(0)
func WithDisabledLimits() Option {
	return func(c *Config) {
		c.LimitConnection = false
	}
}

// WithDriver will enforce the usage of a specific [driver.Driver].
//
// Available options are:
//   - [sqlite.DriverMattn] for "github.com/mattn/go-sqlite3"
//   - [sqlite.DriverModernc] for "modernc.org/sqlite"
func WithDriver(driver Driver) Option {
	return func(c *Config) {
		c.Driver = driver
	}
}

// WithDriverName will set the name used for [sql.Open].
//
// Normally you just need this, if you modified the sql [driver.Driver]
// This happens normally if you want to add tracing, profiling, etc.
func WithDriverName(name string) Option {
	return func(c *Config) {
		c.DriverName = name
	}
}

// WithForeignKeySupport will enable or disable the foreign key support.
//
// See https://www.sqlite.org/pragma.html#pragma_foreign_keys.
func WithForeignKeySupport(enabled bool) Option {
	return func(c *Config) {
		c.ForeignKey = enabled
	}
}

// WithJournalMode will set the journal mode for the connection.
//
// Setting the value [sqlite.JournalDefault] will not set the pragma at all and uses the driver default behaviour.
//
// See https://www.sqlite.org/pragma.html#pragma_journal_mode.
func WithJournalMode(mode JournalMode) Option {
	return func(c *Config) {
		c.JournalMode = mode
	}
}

// WithJournalSizeLimit will set the journal size limit.
//
// Setting a value of 0 will not set the pragma at all and uses the driver default behaviour.
//
// See https://www.sqlite.org/pragma.html#pragma_journal_size_limit.
func WithJournalSizeLimit(limit int) Option {
	return func(c *Config) {
		c.JournalSizeLimit = limit
	}
}

// WithPath will set the db path for sqlite
//
// dbPath should be in format "file:your/path/to/data.db" or ":memory" for an in-memory sqlite connection.
// The format will be checked per regex `^file\:.+\..+$` on [sqlite.Connect].
func WithPath(dbPath string) Option {
	return func(c *Config) {
		c.Path = dbPath
	}
}

// WithSyncMode will set the sync mode for the connection.
//
// Setting the value [sqlite.SyncDefault] will not set the pragma at all and uses the driver default behaviour.
//
// See https://www.sqlite.org/pragma.html#pragma_synchronous.
func WithSyncMode(sync SyncMode) Option {
	return func(c *Config) {
		c.SyncMode = sync
	}
}
