package libsql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func Test_statementSuite(t *testing.T) {
	suite.Run(t, &StatementSuite{})
}

type StatementSuite struct {
	suite.Suite

	sqlStatement *SqlStmtMock
	scan         *ScanDoerMock

	statement statementImpl
}

var (
	_ suite.SetupTestSuite    = (*StatementSuite)(nil)
	_ suite.TearDownTestSuite = (*StatementSuite)(nil)
)

// SetupTest implements SetupTestSuite.SetupTest
func (s *StatementSuite) SetupTest() {
	s.sqlStatement = NewSqlStmtMock(s.T())
	s.scan = NewScanDoerMock(s.T())

	s.statement = statementImpl{statement: s.sqlStatement, scan: s.scan}
}

// TearDownTest implements TearDownTestSuite.TearDownTest
func (s *StatementSuite) TearDownTest() {
	s.scan.MinimockFinish()
	s.sqlStatement.MinimockFinish()
}

func (s *StatementSuite) TestScan() {
	s.doTestScan(s.statement.Scan, false)
}

func (s *StatementSuite) TestScanOne() {
	s.doTestScan(s.statement.ScanOne, true)
}

func (s *StatementSuite) TestUpdate() {

	sqlResultMock := NewSqlResultMock(s.T())
	defer sqlResultMock.MinimockFinish()

	expCtx := context.Background()
	expArgs := []interface{}{1, 94}

	s.sqlStatement.ExecMock.When(
		expCtx,
		expArgs...,
	).Then(sqlResultMock, (error)(nil))

	actualResult, err := s.statement.Update(expCtx, expArgs...)
	s.Require().NoError(err)
	s.Require().Equal(sqlResultMock, actualResult)
}

func (s *StatementSuite) TestUpdateAndGetRowsAffected() {
	sqlResultMock := NewSqlResultMock(s.T())
	defer sqlResultMock.MinimockFinish()

	expCtx := context.Background()
	expArgs := []interface{}{1, 94}
	expRowsAffected := int64(11)

	s.sqlStatement.ExecMock.When(
		expCtx,
		expArgs...,
	).Then(sqlResultMock, (error)(nil))

	sqlResultMock.RowsAffectedMock.Return(expRowsAffected, (error)(nil))

	actualRowsAffected, err := s.statement.UpdateAndGetRowsAffected(expCtx, expArgs...)
	s.Require().NoError(err)
	s.Require().Equal(expRowsAffected, actualRowsAffected)
}

func (s *StatementSuite) TestUpdateAndGetLastInsertID() {
	sqlResultMock := NewSqlResultMock(s.T())
	defer sqlResultMock.MinimockFinish()

	expCtx := context.Background()
	expArgs := []interface{}{1, 94}
	expLastInsertID := int64(11)

	s.sqlStatement.ExecMock.When(
		expCtx,
		expArgs...,
	).Then(sqlResultMock, (error)(nil))

	sqlResultMock.LastInsertIdMock.Return(expLastInsertID, (error)(nil))

	actualLastInsertID, err := s.statement.UpdateAndGetLastInsertID(expCtx, expArgs...)
	s.Require().NoError(err)
	s.Require().Equal(expLastInsertID, actualLastInsertID)
}

func (s *StatementSuite) doTestScan(
	scan func(context.Context, RowScanner, ...interface{}) error,
	oneRow bool,
) {
	expRowScanner := NewRowScannerMock(s.T())
	defer expRowScanner.MinimockFinish()

	expSqlRows := NewSqlRowsMock(s.T())
	defer expSqlRows.MinimockFinish()

	expCtx := context.Background()
	expArgs := []interface{}{1, 94}

	s.sqlStatement.QueryMock.When(
		expCtx,
		expArgs...,
	).Then(expSqlRows, (error)(nil))

	s.scan.DoMock.Set(func(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) (err error) {
		// execute query and and the assertion is that the call is
		// delegated to queryer
		_, _ = query()
		return (error)(nil)
	})

	err := scan(expCtx, expRowScanner, expArgs...)
	s.Require().NoError(err)
}
