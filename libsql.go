package libsql // import "oss.indeed.com/go/libsql"

import (
	"context"
	"database/sql"
	"errors"
	"io"
)

// Wrap returns a Database wrapping a given *sql.DB.
func Wrap(db *sql.DB) Database {
	return newDatabase(newSQLDB(db))
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Database -o libsqltest/ -s _mock.go

// Database provides fluent database API.
type Database interface {
	io.Closer

	Queryer
	Preparer

	// Transaction performs work in transaction.
	// Transaction is committed if work returns nil, and rolled back otherwise.
	Transaction(ctx context.Context, work func(Transaction) error) error
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Transaction -o libsqltest/ -s _mock.go

// Transaction represents an open transaction
type Transaction interface {
	Queryer
	Preparer
}

// ErrNoRows is returned by ScanOne when a query returns no rows
var ErrNoRows = errors.New("no rows, expected 1")

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Queryer -o libsqltest/ -s _mock.go

// Queryer performs scans and updates
type Queryer interface {
	// Scan executes sql and scans result rows with RowScanner
	Scan(ctx context.Context, scanner RowScanner, sql string, args ...interface{}) error

	// Scans executes sql and scans the first result row with RowScanner.
	// Returns ErrNoRows if no rows were returned. Remaining rows are discarded
	ScanOne(ctx context.Context, scanner RowScanner, sql string, args ...interface{}) error

	// Update executes sql insert, update, or delete
	Update(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)

	// UpdateAndGetRowsAffected sql insert, update, or delete and returns affected row count
	// Shorthand for Update(...) followed by UpdateResult.RowsAffected
	UpdateAndGetRowsAffected(ctx context.Context, sql string, args ...interface{}) (int64, error)

	// InsertAndGetLastInsertID sql insert, update, or delete and returns last generated row id.
	// Shorthand for Update(...) followed by UpdateResult.LastInsertId
	UpdateAndGetLastInsertID(ctx context.Context, sql string, args ...interface{}) (int64, error)
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Preparer -o libsqltest/ -s _mock.go

// Preparer prepares a sql statement
type Preparer interface {
	// Prepared prepares a sql statement and runs work.
	// The prepared statement is always closed when this method returns.
	Prepared(ctx context.Context, sql string, work func(Statement) error) error
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Statement -o libsqltest/ -s _mock.go

// Statement is a prepared statement
type Statement interface {
	// Scan executes the prepared statement and scans result rows with RowScanner
	Scan(ctx context.Context, scanner RowScanner, args ...interface{}) error

	// Scans executes the prepared statement and scans the first result row with RowScanner.
	// Returns ErrNoRows if no rows were returned. Remaining rows are discarded
	ScanOne(ctx context.Context, scanner RowScanner, args ...interface{}) error

	// Update executes the prepared insert, update, or delete
	Update(ctx context.Context, args ...interface{}) (sql.Result, error)

	// UpdateAndGetRowsAffected the prepared insert, update, or delete and returns affected row count
	// Shorthand for Update(...) followed by UpdateResult.RowsAffected
	UpdateAndGetRowsAffected(ctx context.Context, args ...interface{}) (int64, error)

	// InsertAndGetLastInsertID the prepared insert, update, or delete and returns last generated row id.
	// Shorthand for Update(...) followed by UpdateResult.LastInsertId
	UpdateAndGetLastInsertID(ctx context.Context, args ...interface{}) (int64, error)
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i RowScanner -o libsqltest/ -s _mock.go

// RowScanner scans database rows into arbitrary data structures
type RowScanner interface {
	// Into returns a slice of pointers to scan into.
	// Should always return a slice of same pointers.
	Into() []interface{}

	// RowScanned notifies scanner that a row has been scanned into slice returned by Into()
	RowScanned() error
}

// Into creates a RowScanner that scans values into the passed pointers.
// Only suitable for scanning single or last row.
func Into(valuePointers ...interface{}) RowScanner {
	return simpleScanner(valuePointers)
}

// FeedScanner feeds the rows to scanner.
// NOTE: Only use this func in tests. Value conversion may differ from the one used in actual SQL execution.
func FeedScanner(scanner RowScanner, rows ...[]interface{}) error {
	return feedScanner(scanner, rows...)
}
