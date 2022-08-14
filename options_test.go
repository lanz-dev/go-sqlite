package sqlite

import (
	"testing"
)

func optionRunner(config *Config, opts ...Option) {
	for _, opt := range opts {
		opt(config)
	}
}

func TestWithAutoVacuumMode(t *testing.T) {
	t.Parallel()

	expected := AutoVacuumFull

	config := newConfig()
	optionRunner(
		config,
		WithAutoVacuumMode(expected),
	)

	got := config.AutoVacuumMode
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestWithBusyTimeout(t *testing.T) {
	t.Parallel()

	expected := 100

	config := newConfig()
	optionRunner(
		config,
		WithBusyTimeout(expected),
	)

	got := config.BusyTimeout
	if got != expected {
		t.Errorf("expected '%d', got '%d'", expected, got)
	}
}

func TestWithCaseSensitiveLike(t *testing.T) {
	t.Parallel()

	config := newConfig()
	optionRunner(
		config,
		WithCaseSensitiveLike(true),
	)

	got := config.CaseSensitiveLike
	if !got {
		t.Errorf("expected '%v', got '%v'", true, got)
	}
}

func TestWithDeferredForeignKeys(t *testing.T) {
	t.Parallel()

	config := newConfig()
	optionRunner(
		config,
		WithDeferredForeignKeys(true),
	)

	got := config.DeferForeignKeys
	if !got {
		t.Errorf("expected '%v', got '%v'", true, got)
	}
}

func TestWithDisabledLimits(t *testing.T) {
	t.Parallel()

	config := newConfig()
	optionRunner(
		config,
		WithDisabledLimits(),
	)

	got := config.LimitConnection
	if got {
		t.Errorf("expected '%v', got '%v'", false, got)
	}
}

func TestWithDriver(t *testing.T) {
	t.Parallel()

	expected := DriverModernc

	config := newConfig()
	optionRunner(
		config,
		WithDriver(expected),
	)

	got := config.Driver
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestWithDriverName(t *testing.T) {
	t.Parallel()

	expected := DriverNameModernc

	config := newConfig()
	optionRunner(
		config,
		WithDriverName(expected),
	)

	got := config.DriverName
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestWithForeignKeySupport(t *testing.T) {
	t.Parallel()

	config := newConfig()
	optionRunner(
		config,
		WithForeignKeySupport(false),
	)

	got := config.ForeignKey
	if got {
		t.Errorf("expected '%v', got '%v'", false, got)
	}
}

func TestWithJournalMode(t *testing.T) {
	t.Parallel()

	expected := JournalPersist

	config := newConfig()
	optionRunner(
		config,
		WithJournalMode(expected),
	)

	got := config.JournalMode
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestWithJournalSizeLimit(t *testing.T) {
	t.Parallel()

	expected := 100

	config := newConfig()
	optionRunner(
		config,
		WithJournalSizeLimit(expected),
	)

	got := config.JournalSizeLimit
	if got != expected {
		t.Errorf("expected '%d', got '%d'", expected, got)
	}
}

func TestWithPath(t *testing.T) {
	t.Parallel()

	expected := "file:/data.db"

	config := newConfig()
	optionRunner(
		config,
		WithPath(expected),
	)

	got := config.Path
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestWithSyncMode(t *testing.T) {
	t.Parallel()

	expected := SyncExtra

	config := newConfig()
	optionRunner(
		config,
		WithSyncMode(expected),
	)

	got := config.SyncMode
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}
