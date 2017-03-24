package libsql

import (
	"context"
	"database/sql"
)

func newStatement(sqlStatement sqlStmt) Statement {
	return &statementImpl{statement: sqlStatement, scan: defaultScanDoer()}
}

type statementImpl struct {
	statement sqlStmt
	scan      scanDoer
}

var _ Statement = (*statementImpl)(nil)

// Scan implements Statement.Scan
func (s statementImpl) Scan(ctx context.Context, scanner RowScanner, args ...interface{}) error {
	return s.scan.Do(scanner, false, s.queryFunc(ctx, args...))
}

// ScanOne implements Statement.ScanOne
func (s statementImpl) ScanOne(ctx context.Context, scanner RowScanner, args ...interface{}) error {
	return s.scan.Do(scanner, true, s.queryFunc(ctx, args...))
}

// Update implements Statement.Update
func (s statementImpl) Update(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return s.statement.Exec(ctx, args...)
}

// UpdateAndGetRowsAffected implements Statement.UpdateAndGetRowsAffected
func (s statementImpl) UpdateAndGetRowsAffected(ctx context.Context, args ...interface{}) (int64, error) {
	return rowsAffected(s.Update(ctx, args...))
}

// UpdateAndGetLastInsertID implements Statement.UpdateAndGetLastInsertID
func (s statementImpl) UpdateAndGetLastInsertID(ctx context.Context, args ...interface{}) (int64, error) {
	return lastInsertID(s.Update(ctx, args...))
}

func (s statementImpl) queryFunc(ctx context.Context, args ...interface{}) func() (sqlRows, error) {
	return func() (sqlRows, error) {
		return s.statement.Query(ctx, args...)
	}
}
