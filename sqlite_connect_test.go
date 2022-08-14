package sqlite

import (
	"testing"
)

func Test_buildDSN_EmptyParams(t *testing.T) {
	t.Parallel()

	dsn := buildDSN(":memory", []string{})

	expected := ":memory"
	if dsn != expected {
		t.Fatalf("expected dsn '%s', got '%s'", expected, dsn)
	}
}

func Test_buildMattnDSN(t *testing.T) {
	t.Parallel()

	c := newConfig()
	c.AutoVacuumMode = AutoVacuumIncremental
	c.BusyTimeout = 100
	c.CaseSensitiveLike = true
	c.ForeignKey = true
	c.DeferForeignKeys = true
	c.JournalMode = JournalPersist
	c.SyncMode = SyncExtra

	dsn := buildMattnDSN(c)
	expected := "?_auto_vacuum=2&_timeout=100&_case_sensitive_like=true&_fk=true&defer_foreign_keys=true&_journal=PERSIST&_sync=3"
	if dsn != expected {
		t.Fatalf("expected dsn '%s', got '%s'", expected, dsn)
	}
}

func Test_buildMattnDSN_DeferForeignKeysJustWithForeignKey(t *testing.T) {
	t.Parallel()

	c := newConfig()
	c.ForeignKey = false
	c.DeferForeignKeys = true

	dsn := buildMattnDSN(c)
	expected := "?_timeout=4000&_journal=WAL&_sync=1"
	if dsn != expected {
		t.Fatalf("expected dsn '%s', got '%s'", expected, dsn)
	}
}

func Test_buildMattnDSN_InvalidVacuumMode(t *testing.T) {
	t.Parallel()

	c := newConfig()
	c.AutoVacuumMode = "invalid"

	dsn := buildMattnDSN(c)
	expected := "?_timeout=4000&_fk=true&_journal=WAL&_sync=1"
	if dsn != expected {
		t.Fatalf("expected dsn '%s', got '%s'", expected, dsn)
	}
}

func Test_buildMattnDSN_InvalidSyncMode(t *testing.T) {
	t.Parallel()

	c := newConfig()
	c.SyncMode = "invalid"

	dsn := buildMattnDSN(c)
	expected := "?_timeout=4000&_fk=true&_journal=WAL"
	if dsn != expected {
		t.Fatalf("expected dsn '%s', got '%s'", expected, dsn)
	}
}

func Test_buildModerncDSN(t *testing.T) {
	t.Parallel()

	c := newConfig()
	c.AutoVacuumMode = AutoVacuumIncremental
	c.BusyTimeout = 100
	c.CaseSensitiveLike = true
	c.ForeignKey = true
	c.DeferForeignKeys = true
	c.JournalMode = JournalPersist
	c.SyncMode = SyncExtra

	dsn := buildModerncDSN(c)
	expected := "?_pragma=auto_vacuum(INCREMENTAL)&_pragma=busy_timeout(100)&_pragma=case_sensitive_like(1)&_pragma=foreign_keys(1)&_pragma=defer_foreign_keys(1)&_pragma=journal_mode(PERSIST)&_pragma=synchronous(EXTRA)"
	if dsn != expected {
		t.Fatalf("expected dsn '%s', got '%s'", expected, dsn)
	}
}

func Test_buildModerncDSN_DeferForeignKeysJustWithForeignKey(t *testing.T) {
	t.Parallel()

	c := newConfig()
	c.ForeignKey = false
	c.DeferForeignKeys = true

	dsn := buildModerncDSN(c)
	expected := "?_pragma=busy_timeout(4000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
	if dsn != expected {
		t.Fatalf("expected dsn '%s', got '%s'", expected, dsn)
	}
}

func Test_buildModerncDSN_InvalidVacuumMode(t *testing.T) {
	t.Parallel()

	c := newConfig()
	c.AutoVacuumMode = "invalid"

	dsn := buildModerncDSN(c)
	expected := "?_pragma=auto_vacuum(invalid)&_pragma=busy_timeout(4000)&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
	if dsn != expected {
		t.Fatalf("expected dsn '%s', got '%s'", expected, dsn)
	}
}

func Test_buildModerncDSN_InvalidSyncMode(t *testing.T) {
	t.Parallel()

	c := newConfig()
	c.SyncMode = "invalid"

	dsn := buildModerncDSN(c)
	expected := "?_pragma=busy_timeout(4000)&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(invalid)"
	if dsn != expected {
		t.Fatalf("expected dsn '%s', got '%s'", expected, dsn)
	}
}
