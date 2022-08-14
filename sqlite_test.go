package sqlite_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/lanz-dev/go-sqlite"
)

func buildMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestOptimize(t *testing.T) {
	db, mock := buildMockDB(t)
	defer db.Close()

	mock.ExpectExec("PRAGMA optimize;").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := sqlite.Optimize(db); err != nil {
		t.Fatalf("did not expected error '%s'", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOptimize_WithError(t *testing.T) {
	db, mock := buildMockDB(t)
	defer db.Close()

	errUnitTest := errors.New("unittest")

	mock.ExpectExec("PRAGMA optimize;").
		WillReturnError(errUnitTest)

	err := sqlite.Optimize(db)
	if !errors.Is(err, errUnitTest) {
		t.Fatalf("expect error to be '%s', got '%s'", errUnitTest, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShutdown(t *testing.T) {
	db, mock := buildMockDB(t)
	defer db.Close()

	mock.ExpectExec("PRAGMA optimize;").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectClose()

	if err := sqlite.Shutdown(db); err != nil {
		t.Fatalf("did not expected error '%s'", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShutdown_WithErrorInOptimize(t *testing.T) {
	db, mock := buildMockDB(t)
	defer db.Close()

	errUnitTest := errors.New("unittest")

	mock.ExpectExec("PRAGMA optimize;").
		WillReturnError(errUnitTest)

	err := sqlite.Shutdown(db)
	if !errors.Is(err, errUnitTest) {
		t.Fatalf("expect error to be '%s', got '%s'", errUnitTest, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShutdown_WithErrorOnClose(t *testing.T) {
	db, mock := buildMockDB(t)
	defer db.Close()

	errUnitTest := errors.New("unittest")

	mock.ExpectExec("PRAGMA optimize;").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectClose().WillReturnError(errUnitTest)

	err := sqlite.Shutdown(db)
	if !errors.Is(err, errUnitTest) {
		t.Fatalf("expect error to be '%s', got '%s'", errUnitTest, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestVacuum(t *testing.T) {
	db, mock := buildMockDB(t)
	defer db.Close()

	mock.ExpectExec("VACUUM;").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := sqlite.Vacuum(db); err != nil {
		t.Fatalf("did not expected error '%s'", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestVacuum_WithError(t *testing.T) {
	db, mock := buildMockDB(t)
	defer db.Close()

	errUnitTest := errors.New("unittest")

	mock.ExpectExec("VACUUM;").
		WillReturnError(errUnitTest)

	err := sqlite.Vacuum(db)
	if !errors.Is(err, errUnitTest) {
		t.Fatalf("expect error to be '%s', got '%s'", errUnitTest, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
