package libsql

// Implementations of interfaces from sqliface.go

import (
	"context"
	"database/sql"
)

func newSQLDB(db *sql.DB) sqlDB {
	return sqlDBImpl{db}
}

func newSQLRows(rows *sql.Rows) sqlRows {
	return rows
}

func newSQLStmt(stmt *sql.Stmt) sqlStmt {
	return sqlStmtImpl{stmt}
}

func newSQLTx(tx *sql.Tx) sqlTx {
	return sqlTxImpl{tx}
}

type sqlDBImpl struct {
	*sql.DB
}

// Query implements sqlQueryer.Query
func (s sqlDBImpl) Query(ctx context.Context, query string, args ...interface{}) (sqlRows, error) {
	rows, err := s.DB.QueryContext(ctx, query, args...)
	return newSQLRows(rows), err
}

// Exec implements sqlQueryer.Exec
func (s sqlDBImpl) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.DB.ExecContext(ctx, query, args...)
}

// Prepare implements sqlPreparer.Prepare
func (s sqlDBImpl) Prepare(ctx context.Context, query string) (sqlStmt, error) {
	stmt, err := s.DB.PrepareContext(ctx, query)
	return newSQLStmt(stmt), err
}

// Begin implements sqlDB.Begin
func (s sqlDBImpl) Begin(ctx context.Context) (sqlTx, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	return newSQLTx(tx), err
}

type sqlStmtImpl struct {
	*sql.Stmt
}

// Query imolements sqlStmt.Query
func (s sqlStmtImpl) Query(ctx context.Context, args ...interface{}) (sqlRows, error) {
	rows, err := s.Stmt.QueryContext(ctx, args...)
	return newSQLRows(rows), err
}

// Exec imolements sqlStmt.Exec
func (s sqlStmtImpl) Exec(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return s.Stmt.ExecContext(ctx, args...)
}

type sqlTxImpl struct {
	*sql.Tx
}

// Query implements sqlQueryer.Query
func (s sqlTxImpl) Query(ctx context.Context, query string, args ...interface{}) (sqlRows, error) {
	rows, err := s.Tx.QueryContext(ctx, query, args...)
	return newSQLRows(rows), err
}

// Exec implements sqlQueryer.Exec
func (s sqlTxImpl) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.Tx.ExecContext(ctx, query, args...)
}

// Prepare implements sqlPreparer.Prepare
func (s sqlTxImpl) Prepare(ctx context.Context, query string) (sqlStmt, error) {
	stmt, err := s.Tx.PrepareContext(ctx, query)
	return newSQLStmt(stmt), err
}
