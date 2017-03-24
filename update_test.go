package libsql

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func Test_rowsAffected_ReturnsError(t *testing.T) {
	expErr := errors.New("a-test-error")
	_, actualError := rowsAffected(nil, expErr)
	require.Equal(t, expErr, actualError)
}

func Test_rowsAffected_CallsRowsAffected(t *testing.T) {
	sqlResult := NewSqlResultMock(t)
	defer sqlResult.MinimockFinish()

	expRowsAffected := int64(941)

	sqlResult.RowsAffectedMock.Return(expRowsAffected, (error)(nil))

	actualRowsAffected, err := rowsAffected(sqlResult, nil)
	require.NoError(t, err)
	require.Equal(t, expRowsAffected, actualRowsAffected)
}

func Test_lastInsertID_ReturnsError(t *testing.T) {
	expErr := errors.New("a-test-error")
	_, actualError := lastInsertID(nil, expErr)
	require.Equal(t, expErr, actualError)
}

func Test_lastInsertID_CallsLastInsertID(t *testing.T) {

	sqlResult := NewSqlResultMock(t)
	defer sqlResult.MinimockFinish()

	expectedLastInsertID := int64(941)

	sqlResult.LastInsertIdMock.Return(expectedLastInsertID, (error)(nil))

	actualLastInsertID, err := lastInsertID(sqlResult, nil)
	require.NoError(t, err)
	require.Equal(t, expectedLastInsertID, actualLastInsertID)
}
