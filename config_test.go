package sqlite

import (
	"errors"
	"testing"
)

func Test_newConfig_Defaults(t *testing.T) {
	t.Parallel()

	config := newConfig()
	if config.LimitConnection != true {
		t.Errorf("expected '%v', got '%v'", true, config.LimitConnection)
	}
	if config.BusyTimeout != 4000 {
		t.Errorf("expected '%d', got '%d'", 4000, config.BusyTimeout)
	}
	if config.ForeignKey != true {
		t.Errorf("expected '%v', got '%v'", true, config.ForeignKey)
	}
	if config.JournalMode != JournalWAL {
		t.Errorf("expected '%s', got '%s'", JournalWAL, config.JournalMode)
	}
	if config.JournalSizeLimit != 100000000 {
		t.Errorf("expected '%d', got '%d'", 100000000, config.JournalSizeLimit)
	}
	if config.SyncMode != SyncNormal {
		t.Errorf("expected '%s', got '%s'", SyncNormal, config.SyncMode)
	}
}

func Test_buildConfig_DefaultPath(t *testing.T) {
	t.Parallel()

	config, err := buildConfig(
		WithDriver(DriverModernc),
	)
	if err != nil {
		t.Fatalf("did not expect error '%s'", err)
	}
	expected := ":memory"
	if config.Path != expected {
		t.Fatalf("expected '%s', got '%s'", expected, config.Path)
	}
}

func Test_buildConfig_InvalidPath(t *testing.T) {
	t.Parallel()

	_, err := buildConfig(
		WithPath("invalid"),
		WithDriver(DriverModernc),
	)
	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("expected to receive error '%s', got '%s'", ErrInvalidPath, err)
	}
}

func Test_buildConfig_DriverNameModernc(t *testing.T) {
	t.Parallel()

	config, err := buildConfig(
		WithDriver(DriverModernc),
	)
	if err != nil {
		t.Fatalf("did not expect error '%s'", err)
	}
	expected := DriverNameModernc
	if config.DriverName != expected {
		t.Fatalf("expected '%s', got '%s'", expected, config.DriverName)
	}
}

func Test_buildConfig_DriverNameMattn(t *testing.T) {
	t.Parallel()

	config, err := buildConfig(
		WithDriver(DriverMattn),
	)
	if err != nil {
		t.Fatalf("did not expect error '%s'", err)
	}
	expected := DriverNameMattn
	if config.DriverName != expected {
		t.Fatalf("expected '%s', got '%s'", expected, config.DriverName)
	}
}

func Test_buildConfig_DSNModernc(t *testing.T) {
	t.Parallel()

	config, err := buildConfig(
		WithDriver(DriverModernc),
	)
	if err != nil {
		t.Fatalf("did not expect error '%s'", err)
	}
	expected := ":memory?_pragma=busy_timeout(4000)&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
	if config.DSN != expected {
		t.Fatalf("expected '%s', got '%s'", expected, config.DSN)
	}
}

func Test_buildConfig_DSNMattn(t *testing.T) {
	t.Parallel()

	config, err := buildConfig(
		WithDriver(DriverMattn),
	)
	if err != nil {
		t.Fatalf("did not expect error '%s'", err)
	}
	expected := ":memory?_timeout=4000&_fk=true&_journal=WAL&_sync=1"
	if config.DSN != expected {
		t.Fatalf("expected '%s', got '%s'", expected, config.DSN)
	}
}

func TestAutoVacuumMode_Int(t *testing.T) {
	tests := []struct {
		name string
		s    AutoVacuumMode
		want int
	}{
		{"AutoVaccumNone", AutoVacuumNone, 0},
		{"AutoVacuumFull", AutoVacuumFull, 1},
		{"AutoVacuumIncremental", AutoVacuumIncremental, 2},
		{"AutoVacuumInvalid", AutoVacuumMode("invalid"), -99},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.s.Int(); got != tc.want {
				t.Errorf("Int() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSyncMode_Int(t *testing.T) {
	tests := []struct {
		name string
		s    SyncMode
		want int
	}{
		{"SyncOff", SyncOff, 0},
		{"SyncNormal", SyncNormal, 1},
		{"SyncFull", SyncFull, 2},
		{"SyncExtra", SyncExtra, 3},
		{"SyncInvalid", SyncMode("invalid"), -99},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.s.Int(); got != tc.want {
				t.Errorf("Int() = %v, want %v", got, tc.want)
			}
		})
	}
}

func Test_validatePath(t *testing.T) {
	tests := []struct {
		path    string
		wantErr error
	}{
		{"invalid:/path", ErrInvalidPath},
		{"invalid", ErrInvalidPath},
		{"file", ErrInvalidPath},
		{"file:", ErrInvalidPath},
		{"file:/", ErrInvalidPath},
		{"file:/fsafd", ErrInvalidPath},
		{"file:/data", ErrInvalidPath},
		{"file:/data.", ErrInvalidPath},
		{"file:/data.db", nil},
		{"file:data.db", nil},
		{"file:../../../../../../data.db", nil},
		{":memory", nil},
	}
	for _, tc := range tests {
		tc := tc

		t.Run("With path '"+tc.path+"'", func(t *testing.T) {
			t.Parallel()

			err := validatePath(tc.path)
			if err != nil && tc.wantErr == nil {
				t.Fatalf("did not expect error '%s'", err)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("expect error to be '%s', got '%s'", tc.wantErr, err)
			}
			if tc.wantErr == nil && err != nil {
				t.Fatalf("expected '%s' to be valid", tc.path)
			}
		})
	}
}
