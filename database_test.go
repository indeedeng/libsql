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

func Test_databaseImpl_CloseIsPropagated(t *testing.T) {
	sqlDB := NewSqlDBMock(t)
	defer sqlDB.MinimockFinish()

	expErr := errors.New("a-test-error")

	sqlDB.CloseMock.Return(expErr)

	actualError := newDatabase(sqlDB).Close()
	require.Equal(t, expErr, actualError)
}
