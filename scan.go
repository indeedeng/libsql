package libsql

import (
	"database/sql"
	"io"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i scanDoer -o scan_mock_test.go

type scanDoer interface {
	// Do scans row(s) returned by the query using rowScanner
	Do(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) error
}

func defaultScanDoer() scanDoer {
	return scanDoerFunc(scan)
}

type scanDoerFunc func(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) error

var _ scanDoer = (scanDoerFunc)(nil)

// Do implements scanDoer.Do
func (f scanDoerFunc) Do(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) error {
	return f(rowScanner, oneRow, query)
}

func ignoreClose(c io.Closer) {
	_ = c.Close()
}

func scan(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) error {
	rows, err := query()
	if err != nil {
		return err
	}

	defer ignoreClose(rows)

	rowsScanned := 0

	for (!oneRow || rowsScanned < 1) && rows.Next() {
		if err := rows.Scan(rowScanner.Into()...); err != nil {
			return err
		}
		if err := rowScanner.RowScanned(); err != nil {
			return err
		}
		rowsScanned++
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if oneRow && rowsScanned != 1 {
		return ErrNoRows
	}

	return nil
}

func feedScanner(scanner RowScanner, rows ...[]interface{}) error {
	for _, row := range rows {
		if err := feedRow(scanner, row); err != nil {
			return err
		}
	}
	return nil
}

func feedRow(scanner RowScanner, row []interface{}) error {
	for idx, v := range row {
		// this is a simplified version of convertAssign from database/sql/convert.go
		toPtr := scanner.Into()[idx]
		switch toPtr.(type) {
		case sql.Scanner:
			err := toPtr.(sql.Scanner).Scan(v)
			if err != nil {
				return errors.Wrap(err, "failed to scan into column "+strconv.Itoa(idx))
			}
		default:
			to := reflect.Indirect(reflect.ValueOf(toPtr))
			to.Set(reflect.ValueOf(v))
		}
	}
	return scanner.RowScanned()
}

type simpleScanner []interface{}

var _ RowScanner = (*simpleScanner)(nil)

// Into implements RowScanner.Into
func (s simpleScanner) Into() []interface{} {
	return s
}

// RowScanned implements RowScanner.RowScanned
func (s simpleScanner) RowScanned() error {
	return nil
}
