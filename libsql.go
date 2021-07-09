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

	// PrepareStatement prepares a statement for later queries.
	//
	// In addition to preparing a statement on a single connection, the returned
	// PreparedStatement can use other connections from the Database's connection
	// pool re-preparing the statement as needed.
	// PreparedStatement is safe for concurrent use.
	// Refer to sql.Stmt's godoc for details.
	//
	// It is best to use this API for queries that are executed frequently from
	// different contexts. For example, using a PreparedStatement to fetch data
	// served via a service's HTTP API endpoint by many goroutines can reduce
	// the latency considerably.
	// In other scenarios, prefer using the Prepared method: it takes care of
	// closing the Statement and is thus less error-prone.
	//
	// The provided context is used for the preparation of the statement, not
	// for its execution.
	// The caller must call Close on the returned PreparedStatement when it is
	// no longer needed.
	PrepareStatement(ctx context.Context, sql string) (PreparedStatement, error)
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

	// ScanOne executes sql and scans the first result row with RowScanner.
	// Returns ErrNoRows if no rows were returned. Remaining rows are discarded
	ScanOne(ctx context.Context, scanner RowScanner, sql string, args ...interface{}) error

	// Update executes sql insert, update, or delete
	Update(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)

	// UpdateAndGetRowsAffected sql insert, update, or delete and returns affected row count
	// Shorthand for Update(...) followed by UpdateResult.RowsAffected
	UpdateAndGetRowsAffected(ctx context.Context, sql string, args ...interface{}) (int64, error)

	// UpdateAndGetLastInsertID sql insert, update, or delete and returns last generated row id.
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

	// ScanOne executes the prepared statement and scans the first result row with RowScanner.
	// Returns ErrNoRows if no rows were returned. Remaining rows are discarded
	ScanOne(ctx context.Context, scanner RowScanner, args ...interface{}) error

	// Update executes the prepared insert, update, or delete
	Update(ctx context.Context, args ...interface{}) (sql.Result, error)

	// UpdateAndGetRowsAffected the prepared insert, update, or delete and returns affected row count
	// Shorthand for Update(...) followed by UpdateResult.RowsAffected
	UpdateAndGetRowsAffected(ctx context.Context, args ...interface{}) (int64, error)

	// UpdateAndGetLastInsertID the prepared insert, update, or delete and returns last generated row id.
	// Shorthand for Update(...) followed by UpdateResult.LastInsertId
	UpdateAndGetLastInsertID(ctx context.Context, args ...interface{}) (int64, error)
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i PreparedStatement -o libsqltest/ -s _mock.go

// PreparedStatement is a Statement that must be closed by the caller
type PreparedStatement interface {
	io.Closer
	Statement
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
