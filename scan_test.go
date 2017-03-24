package libsql

import (
	"database/sql"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func Test_scan(t *testing.T) {
	rowScannerMock := NewRowScannerMock(t)
	defer rowScannerMock.MinimockFinish()

	sqlRowsMock := NewSqlRowsMock(t)
	defer sqlRowsMock.MinimockFinish()

	var column1 int
	var column2 int
	scannerTargets := []interface{}{&column1, &column2}

	// what we return for sqlRows.Next changes (true, true, false)
	mockRowsNextCounter := 0
	sqlRowsMock.NextMock.Set(func() bool {
		mockRowsNextCounter++
		switch mockRowsNextCounter {
		case 1, 2:
			t.Logf("sqlRowsMock.Next => true")
			return true
		case 3:
			t.Logf("sqlRowsMock.Next => false")
			return false
		default:
			panic("too many calls to Next")
		}
	})

	rowScannerMock.IntoMock.Return(scannerTargets)

	sqlRowsMock.ScanMock.Expect(scannerTargets...).Return((error)(nil))
	sqlRowsMock.CloseMock.Return((error)(nil))

	sqlRowsMock.ScanMock.Expect(scannerTargets...).Return((error)(nil))
	sqlRowsMock.CloseMock.Return((error)(nil))

	rowScannerMock.RowScannedMock.Return((error)(nil))
	sqlRowsMock.ErrMock.Return((error)(nil))

	err := scan(rowScannerMock, false, func() (sqlRows, error) {
		return sqlRowsMock, nil
	})
	require.NoError(t, err)
}

func Test_scan_queryError(t *testing.T) {
	rowScannerMock := NewRowScannerMock(t)
	defer rowScannerMock.MinimockFinish()

	sqlRowsMock := NewSqlRowsMock(t)
	defer sqlRowsMock.MinimockFinish()

	expErr := errors.New("a-test-error")
	err := scan(rowScannerMock, false, func() (sqlRows, error) {
		return sqlRowsMock, expErr
	})
	require.Error(t, err)
	require.Equal(t, expErr, err)
}

func Test_scan_scanError(t *testing.T) {
	rowScannerMock := NewRowScannerMock(t)
	defer rowScannerMock.MinimockFinish()

	sqlRowsMock := NewSqlRowsMock(t)
	defer sqlRowsMock.MinimockFinish()

	var column1 int
	var column2 int
	scannerTargets := []interface{}{&column1, &column2}

	rowScannerMock.IntoMock.Return(scannerTargets)

	expErr := errors.New("a-test-error")
	sqlRowsMock.NextMock.Return(true)
	sqlRowsMock.ScanMock.Expect(scannerTargets...).Return(expErr)
	sqlRowsMock.CloseMock.Return((error)(nil))

	err := scan(rowScannerMock, false, func() (sqlRows, error) {
		return sqlRowsMock, nil
	})
	require.Error(t, err)
	require.Equal(t, expErr, err)
}

func Test_scan_scannerError(t *testing.T) {
	rowScannerMock := NewRowScannerMock(t)
	defer rowScannerMock.MinimockFinish()

	sqlRowsMock := NewSqlRowsMock(t)
	defer sqlRowsMock.MinimockFinish()

	var column1 int
	var column2 int
	scannerTargets := []interface{}{&column1, &column2}

	rowScannerMock.IntoMock.Return(scannerTargets)

	expErr := errors.New("a-test-error")
	rowScannerMock.RowScannedMock.Return(expErr)

	sqlRowsMock.NextMock.Return(true)
	sqlRowsMock.ScanMock.Expect(scannerTargets...).Return((error)(nil))
	sqlRowsMock.CloseMock.Return((error)(nil))

	err := scan(rowScannerMock, false, func() (sqlRows, error) {
		return sqlRowsMock, nil
	})
	require.Error(t, err)
	require.Equal(t, expErr, err)

}

func Test_scan_iterationError(t *testing.T) {
	rowScannerMock := NewRowScannerMock(t)
	defer rowScannerMock.MinimockFinish()

	sqlRowsMock := NewSqlRowsMock(t)
	defer sqlRowsMock.MinimockFinish()

	sqlRowsMock.NextMock.Return(false)
	expErr := errors.New("a-test-error")
	sqlRowsMock.ErrMock.Return(expErr)
	sqlRowsMock.CloseMock.Return((error)(nil))

	err := scan(rowScannerMock, false, func() (sqlRows, error) {
		return sqlRowsMock, nil
	})
	require.Error(t, err)
	require.Equal(t, expErr, err)
}

func Test_scan_oneRowNoResults(t *testing.T) {
	rowScannerMock := NewRowScannerMock(t)
	defer rowScannerMock.MinimockFinish()

	sqlRowsMock := NewSqlRowsMock(t)
	defer sqlRowsMock.MinimockFinish()

	sqlRowsMock.NextMock.Return(false)
	sqlRowsMock.ErrMock.Return((error)(nil))
	sqlRowsMock.CloseMock.Return((error)(nil))

	err := scan(rowScannerMock, true, func() (sqlRows, error) {
		return sqlRowsMock, nil
	})
	require.Error(t, err)
	require.Equal(t, ErrNoRows, err)
}

func Test_FeedScanner(t *testing.T) {
	rowScannerMock := NewRowScannerMock(t)
	defer rowScannerMock.MinimockFinish()

	sqlRowsMock := NewSqlRowsMock(t)
	defer sqlRowsMock.MinimockFinish()

	var column1 int
	var column2 int
	scannerTargets := []interface{}{&column1, &column2}

	row1Column1 := 11
	row1Column2 := 12
	row2Column1 := 21
	row2Column2 := 22

	rowScannerMock.IntoMock.Return(scannerTargets)

	rowScannedCount := 0
	rowScannerMock.RowScannedMock.Set(func() error {
		rowScannedCount++
		// make assertions about what was passed in, given the count

		switch rowScannedCount {
		case 1:
			// assert r1c1 == c1
			// assert r1c2 == c2
			require.Equal(t, row1Column1, column1)
			require.Equal(t, row1Column2, column2)
		case 2:
			// assert r2c1 == c1
			// assert r2c2 == c2
			require.Equal(t, row2Column1, column1)
			require.Equal(t, row2Column2, column2)
		default:
			panic("RowScanned called too many times")
		}
		return (error)(nil)
	})

	err := FeedScanner(rowScannerMock,
		[]interface{}{row1Column1, row1Column2},
		[]interface{}{row2Column1, row2Column2})
	require.NoError(t, err)
}

func Test_FeedScanner_RowScannerErrorIsPropagated(t *testing.T) {
	rowScannerMock := NewRowScannerMock(t)
	defer rowScannerMock.MinimockFinish()

	expErr := errors.New("a-test-error")
	rowScannerMock.RowScannedMock.Return(expErr)

	actualError := FeedScanner(rowScannerMock, []interface{}{})
	require.Error(t, actualError)
	require.Equal(t, expErr, actualError)
}

func Test_FeedScanner_ValueScannerErrorIsPropagated(t *testing.T) {
	valueScanner := sql.Scanner(&sql.NullInt64{})
	err := FeedScanner(Into(valueScanner), []interface{}{
		"not-a-number",
	})
	require.Error(t, err)
}

func Test_Into(t *testing.T) {
	var column1 int
	var column2 int

	scanner := Into(&column1, &column2)
	require.Equal(t, []interface{}{&column1, &column2}, scanner.Into())
	require.NoError(t, scanner.RowScanned())
}
