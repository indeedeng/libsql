package libsql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func Test_databaseImpl_Transaction(t *testing.T) {
	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	expSQLTx := NewSqlTxMock(t)
	defer expSQLTx.MinimockFinish()

	expCtx := context.Background()

	sqlDB.BeginMock.When(expCtx).Then(expSQLTx, (error)(nil))

	expSQLTx.CommitMock.Return((error)(nil))

	expSQLTx.RollbackMock.Return(sql.ErrTxDone)

	expTx := newTransaction(expSQLTx)

	newTXFuncCalls := 0
	newTXFunc := func(actualSQLTX sqlTx) Transaction {
		require.Equal(t, expSQLTx, actualSQLTX)
		newTXFuncCalls++
		return expTx
	}

	workFuncCalls := 0
	workFunc := func(actualTX Transaction) error {
		require.Equal(t, expTx, actualTX)
		workFuncCalls++
		return nil
	}

	database := &databaseImpl{
		db:    sqlDB,
		newTX: newTXFunc,
	}
	err := database.Transaction(expCtx, workFunc)
	require.NoError(t, err)

	require.Equal(t, 1, newTXFuncCalls)
	require.Equal(t, 1, workFuncCalls)
}

func Test_databaseImpl_TransactionBeginErrorIsReturned(t *testing.T) {
	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	sqlTx := NewSqlTxMock(t)
	defer sqlTx.MinimockFinish()

	expCtx := context.Background()
	expErr := errors.New("a-test-error")

	sqlDB.BeginMock.When(expCtx).Then(sqlTx, expErr)

	actualError := newDatabase(sqlDB).Transaction(expCtx, nil)
	require.Equal(t, expErr, actualError)
}

func Test_databaseImpl_TransactionWorkErrorIsReturned(t *testing.T) {
	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	expSQLTx := NewSqlTxMock(t)
	defer expSQLTx.MinimockFinish()

	expCtx := context.Background()

	sqlDB.BeginMock.When(expCtx).Then(expSQLTx, (error)(nil))

	expSQLTx.RollbackMock.Return(sql.ErrTxDone)

	expectedTX := newTransaction(expSQLTx)

	newTXFuncCalls := 0
	newTXFunc := func(actualSQLTX sqlTx) Transaction {
		require.Equal(t, expSQLTx, actualSQLTX)
		newTXFuncCalls++
		return expectedTX
	}

	expErr := errors.New("a-test-error")

	workFuncCalls := 0
	workFunc := func(actualTX Transaction) error {
		require.Equal(t, expectedTX, actualTX)
		workFuncCalls++
		return expErr
	}

	database := &databaseImpl{
		db:    sqlDB,
		newTX: newTXFunc,
	}
	actualError := database.Transaction(expCtx, workFunc)
	require.Equal(t, expErr, actualError)

	require.Equal(t, 1, newTXFuncCalls)
	require.Equal(t, 1, workFuncCalls)
}

func Test_databaseImpl_TransactionIsRolledBackOnWorkPanic(t *testing.T) {
	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	expSQLTx := NewSqlTxMock(t)
	defer expSQLTx.MinimockFinish()

	expCtx := context.Background()

	sqlDB.BeginMock.When(expCtx).Then(expSQLTx, error(nil))

	expSQLTx.RollbackMock.Return(sql.ErrTxDone)

	expTx := newTransaction(expSQLTx)

	newTXFuncCalls := 0
	newTXFunc := func(actualSQLTX sqlTx) Transaction {
		require.Equal(t, expSQLTx, actualSQLTX)
		newTXFuncCalls++
		return expTx
	}

	workFuncCalls := 0
	workFunc := func(actualTX Transaction) error {
		require.Equal(t, expTx, actualTX)
		workFuncCalls++
		panic("an-expected-panic")
	}

	database := &databaseImpl{
		db:    sqlDB,
		newTX: newTXFunc,
	}

	require.Panics(t, func() {
		_ = database.Transaction(expCtx, workFunc)
	})

	require.Equal(t, 1, newTXFuncCalls)
	require.Equal(t, 1, workFuncCalls)
}

func Test_databaseImpl_PrepareStatement(t *testing.T) {
	ctx := context.WithValue(context.Background(), "a-key-to-make-a-unique-context", "a-value")
	const expectedQuery = "SELECT 1 FROM DUAL"

	sqlStmt := NewSqlStmtMock(t)
	defer sqlStmt.MinimockFinish()

	expectedExecError := errors.New("an-expected-exec-error")
	sqlStmt.
		ExecMock.Expect(ctx).Return(nil, expectedExecError).
		CloseMock.Expect().Return(nil)

	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	sqlDB.PrepareMock.
		Expect(ctx, expectedQuery).
		Return(sqlStmt, nil)

	s, err := newDatabase(sqlDB).PrepareStatement(ctx, expectedQuery)
	require.NoError(t, err)
	_, err = s.Update(ctx)
	require.Error(t, err)
	require.Equal(t, expectedExecError, err)
	err = s.Close()
	require.NoError(t, err)
}

func Test_databaseImpl_PrepareStatement_propagatesErrors(t *testing.T) {
	ctx := context.WithValue(context.Background(), "a-key-to-make-a-unique-context", "a-value")
	const expectedQuery = "SELECT 1 FROM DUAL"

	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	expectedError := errors.New("a-test-error")
	sqlDB.PrepareMock.
		Expect(ctx, expectedQuery).
		Return(nil, expectedError)

	_, err := newDatabase(sqlDB).PrepareStatement(ctx, expectedQuery)
	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func Test_databaseImpl_CloseIsPropagated(t *testing.T) {
	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	expErr := errors.New("a-test-error")

	sqlDB.CloseMock.Return(expErr)

	actualError := newDatabase(sqlDB).Close()
	require.Equal(t, expErr, actualError)
}
