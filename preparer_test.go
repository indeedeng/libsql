package libsql

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func Test_preparerMixin_Prepare(t *testing.T) {
	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	expSQLStmt := NewSqlStmtMock(t)
	defer expSQLStmt.MinimockFinish()

	expStatement := newStatement(expSQLStmt)

	expQuery := "SELECT something From somewhere"

	expCtx := context.Background()

	sqlDB.PrepareMock.When(expCtx, expQuery).Then(expSQLStmt, (error)(nil))

	expSQLStmt.CloseMock.Return((error)(nil))

	newStatementFuncCalls := 0
	newStatementFunc := func(actualSQLStatement sqlStmt) Statement {
		require.Equal(t, expSQLStmt, actualSQLStatement)
		newStatementFuncCalls++
		return expStatement
	}

	mixin := preparerMixin{
		preparer:     sqlDB,
		newStatement: newStatementFunc,
	}

	workFuncCalls := 0
	workFunc := func(actualStatement Statement) error {
		require.Equal(t, expStatement, actualStatement)
		workFuncCalls++
		return nil
	}

	err := mixin.Prepared(expCtx, expQuery, workFunc)
	require.NoError(t, err)

	require.Equal(t, 1, newStatementFuncCalls)
	require.Equal(t, 1, workFuncCalls)
}

func Test_preparerMixin_PrepareErrorPropagated(t *testing.T) {

	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	sqlStmt := NewSqlStmtMock(t)
	defer sqlStmt.MinimockFinish()

	expQuery := "SELECT something From somewhere"
	expErr := errors.New("a-test-errror")
	expCtx := context.Background()

	sqlDB.PrepareMock.When(expCtx, expQuery).Then(sqlStmt, expErr)

	mixin := preparerMixin{
		preparer: sqlDB,
	}

	actualError := mixin.Prepared(expCtx, expQuery, nil)
	require.Equal(t, expErr, actualError)
}
