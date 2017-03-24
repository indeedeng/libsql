package libsql

import (
	"context"
	"database/sql"
)

func newQueryerMixin(q sqlQueryer) Queryer {
	return queryerMixin{q: q, scan: defaultScanDoer()}
}

type queryerMixin struct {
	q    sqlQueryer
	scan scanDoer
}

var _ Queryer = (*queryerMixin)(nil)

// Scan implements Queryer.Scan
func (m queryerMixin) Scan(ctx context.Context, scanner RowScanner, sql string, args ...interface{}) error {
	return m.scan.Do(scanner, false, m.queryFunc(ctx, sql, args...))
}

// ScanOne implements Queryer.ScanOne
func (m queryerMixin) ScanOne(ctx context.Context, scanner RowScanner, sql string, args ...interface{}) error {
	return m.scan.Do(scanner, true, m.queryFunc(ctx, sql, args...))
}

// Update implements Queryer.Update
func (m queryerMixin) Update(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	return m.q.Exec(ctx, sql, args...)
}

// UpdateAndGetRowsAffected implements Queryer.UpdateAndGetRowsAffected
func (m queryerMixin) UpdateAndGetRowsAffected(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return rowsAffected(m.Update(ctx, sql, args...))
}

// UpdateAndGetLastInsertID implements Queryer.UpdateAndGetLastInsertID
func (m queryerMixin) UpdateAndGetLastInsertID(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return lastInsertID(m.Update(ctx, sql, args...))
}

func (m queryerMixin) queryFunc(ctx context.Context, sql string, args ...interface{}) func() (sqlRows, error) {
	return func() (sqlRows, error) {
		return m.q.Query(ctx, sql, args...)
	}
}
