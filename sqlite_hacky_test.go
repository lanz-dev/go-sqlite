package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	_ "unsafe"

	"github.com/DATA-DOG/go-sqlmock"
)

var (
	errUnitTest = errors.New("unittest")
)

var (
	_ driver.Driver = &unitTestDriver{}
)

type unitTestDriver struct {
	openFn func(name string) (driver.Conn, error)
}

func (u unitTestDriver) Open(name string) (driver.Conn, error) {
	if u.openFn != nil {
		return u.openFn(name)
	}
	return nil, nil
}

//go:linkname unregisterAllDrivers database/sql.unregisterAllDrivers
func unregisterAllDrivers()

//go:linkname drivers database/sql.drivers
var drivers map[string]driver.Driver

func TestDependsOnGlobalRegisterGroup(t *testing.T) {
	sqlMockDriver := drivers["sqlmock"]

	t.Run("Test_detectDriver", func(t *testing.T) {
		tests := []struct {
			name       string
			driverName string
			register   bool
			want       Driver
			wantErr    error
		}{
			{`Detecting "modernc.org/sqlite"`, DriverNameModernc, true, DriverModernc, nil},
			{`Detecting "github.com/mattn/go-sqlite3"`, DriverNameMattn, true, DriverMattn, nil},
			{`No registered driver`, "None", false, "", ErrCantDetectDriver},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				defer unregisterAllDrivers()

				if tc.register {
					sql.Register(tc.driverName, unitTestDriver{})
				}

				got, err := detectDriver()
				if err != nil && tc.wantErr == nil {
					t.Fatalf("did not expect error '%s'", err)
				}
				if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
					t.Fatalf("expect error to be '%s', got '%s'", tc.wantErr, err)
				}
				if tc.wantErr == nil && got != tc.want {
					t.Fatalf("expected to detect Driver '%s', got '%s'", tc.want, got)
				}
			})
		}
	})

	t.Run("Test_buildConfig_WithoutRegisteredDrivers", func(t *testing.T) {
		_, err := buildConfig()
		if !errors.Is(err, ErrCantDetectDriver) {
			t.Fatalf("expect error to be '%s', got '%s'", ErrCantDetectDriver, err)
		}
	})

	t.Run("Test_buildConfig_WithRegisteredDriverModernc", func(t *testing.T) {
		sql.Register(DriverNameModernc, unitTestDriver{})
		defer unregisterAllDrivers()

		config, err := buildConfig()
		if err != nil {
			t.Fatalf("did not expect error '%s'", err)
		}
		if config.Driver != DriverModernc {
			t.Fatalf("expected to detect Driver '%s', got '%s'", DriverModernc, config.Driver)
		}
	})

	t.Run("Test_buildConfig_WithRegisteredDriverMattn", func(t *testing.T) {
		sql.Register(DriverNameMattn, unitTestDriver{})
		defer unregisterAllDrivers()

		config, err := buildConfig()
		if err != nil {
			t.Fatalf("did not expect error '%s'", err)
		}
		if config.Driver != DriverMattn {
			t.Fatalf("expected to detect Driver '%s', got '%s'", DriverMattn, config.Driver)
		}
	})

	t.Run("Test_connect_WithoutRegisteredDrivers", func(t *testing.T) {
		_, err := connect(sql.Open)
		if !errors.Is(err, ErrCantDetectDriver) {
			t.Fatalf("expect error to be '%s', got '%s'", ErrCantDetectDriver, err)
		}
	})

	t.Run("Test_connect_ErrorWithOpen", func(t *testing.T) {
		driverName := "unittest"
		sql.Register(driverName, unitTestDriver{})
		defer unregisterAllDrivers()

		_, err := connect(
			func(_, _ string) (*sql.DB, error) {
				return nil, errUnitTest
			},
			WithDriverName(driverName),
			WithDriver(DriverModernc),
		)
		if !errors.Is(err, errUnitTest) {
			t.Fatalf("expect error to be '%s', got '%s'", errUnitTest, err)
		}
	})

	t.Run("Test_connect_ErrorWithPingAtDriver", func(t *testing.T) {
		driverName := "unittest"
		sql.Register(driverName, unitTestDriver{
			openFn: func(name string) (driver.Conn, error) {
				return nil, errUnitTest
			},
		})
		defer unregisterAllDrivers()

		_, err := connect(
			sql.Open,
			WithDriverName(driverName),
			WithDriver(DriverModernc),
		)
		if !errors.Is(err, errUnitTest) {
			t.Fatalf("expect error to be '%s', got '%s'", errUnitTest, err)
		}
	})

	// dirty hack level > 9000
	drivers["sqlmock"] = sqlMockDriver

	t.Run("Test_connect_ErrorWithPing", func(t *testing.T) {
		db, mock, err := sqlmock.New(
			sqlmock.MonitorPingsOption(true),
		)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectPing().WillReturnError(errUnitTest)

		_, err = connect(
			func(_, _ string) (*sql.DB, error) {
				return db, nil
			},
			WithDriverName("sqlmock"),
			WithDriver(DriverModernc),
		)
		if !errors.Is(err, errUnitTest) {
			t.Fatalf("expect error to be '%s', got '%s'", errUnitTest, err)
		}
	})

	t.Run("Test_connect_LimitConnections", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectExec("PRAGMA .*").WillReturnResult(sqlmock.NewResult(1, 1))

		db, err = connect(
			func(_, _ string) (*sql.DB, error) {
				return db, nil
			},
			WithDriverName("sqlmock"),
			WithDriver(DriverModernc),
		)
		if err != nil {
			t.Fatalf("did not expect error '%s'", err)
		}

		checkInt := func(db *sql.DB, field string, check int) bool {
			v := reflect.ValueOf(db).Elem().FieldByName(field)
			return int(v.Int()) == check
		}

		if !checkInt(db, "maxIdleCount", 1) {
			t.Fatal("expected SetMaxIdleConns 1")
		}
		if !checkInt(db, "maxOpen", 1) {
			t.Fatal("expected SetMaxOpenConns 1")
		}
		if !checkInt(db, "maxLifetime", 0) {
			t.Fatal("expected SetConnMaxIdleTime 0")
		}
		if !checkInt(db, "maxIdleTime", 0) {
			t.Fatal("expected SetConnMaxIdleTime 0")
		}
	})

	t.Run("Test_connect_JournalSizeLimitError", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectExec("PRAGMA .*").WillReturnError(errUnitTest)

		_, err = connect(
			func(_, _ string) (*sql.DB, error) {
				return db, nil
			},
			WithDriverName("sqlmock"),
			WithDriver(DriverModernc),
		)
		if !errors.Is(err, errUnitTest) {
			t.Fatalf("expect error to be '%s', got '%s'", errUnitTest, err)
		}
	})

	t.Run("Test_Connect", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectPing()
		mock.ExpectExec("PRAGMA").
			WillReturnResult(sqlmock.NewResult(1, 1))

		oldOpenFunc := openFunc
		openFunc = func(_, _ string) (*sql.DB, error) {
			return db, nil
		}
		defer func() {
			openFunc = oldOpenFunc
		}()

		db, err = Connect(
			WithDriverName("sqlmock"),
			WithDriver(DriverModernc),
		)
		if err != nil {
			t.Fatalf("did not expected error '%s'", err)
		}
	})

}
