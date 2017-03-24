package libsql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func Test_queryerMixinSuite(t *testing.T) {
	suite.Run(t, &QueryerMixinSuite{})
}

type QueryerMixinSuite struct {
	suite.Suite

	queryer *SqlDBMock
	scan    *ScanDoerMock

	mixin queryerMixin
}

var (
	_ suite.SetupTestSuite    = (*QueryerMixinSuite)(nil)
	_ suite.TearDownTestSuite = (*QueryerMixinSuite)(nil)
)

// SetupTest implements SetupTestSuite.SetupTest
func (s *QueryerMixinSuite) SetupTest() {
	s.queryer = NewSqlDBMock(s.T())
	s.scan = NewScanDoerMock(s.T())
	s.mixin = queryerMixin{q: s.queryer, scan: s.scan}
}

// TearDownTest implements TearDownTestSuite.TearDownTest
func (s *QueryerMixinSuite) TearDownTest() {
	s.queryer.MinimockFinish()
	s.scan.MinimockFinish()
}

func (s *QueryerMixinSuite) TestScan() {
	s.doTestScan(s.mixin.Scan, false)
}

func (s *QueryerMixinSuite) TestScanOne() {
	s.doTestScan(s.mixin.ScanOne, true)
}

func (s *QueryerMixinSuite) TestUpdate() {
	sqlResultMock := NewSqlResultMock(s.T())
	defer sqlResultMock.MinimockFinish()

	expCtx := context.Background()
	expQuery := "UPDATE aTable SET x = ? * ?"
	expArgs := []interface{}{1, 94}

	s.queryer.ExecMock.When(
		expCtx,
		expQuery,
		expArgs...,
	).Then(sqlResultMock, (error)(nil))

	actualResult, err := s.mixin.Update(
		expCtx,
		expQuery,
		expArgs...,
	)
	s.Require().NoError(err)
	s.Require().Equal(sqlResultMock, actualResult)
}

func (s *QueryerMixinSuite) TestUpdateAndGetRowsAffected() {
	sqlResultMock := NewSqlResultMock(s.T())

	expCtx := context.Background()
	expQuery := "UPDATE aTable SET x = ? * ?"
	expArgs := []interface{}{1, 94}
	expRowsAffected := int64(11)

	s.queryer.ExecMock.When(
		expCtx,
		expQuery,
		expArgs...,
	).Then(sqlResultMock, (error)(nil))

	sqlResultMock.RowsAffectedMock.Return(expRowsAffected, (error)(nil))

	actualRowsAffected, err := s.mixin.UpdateAndGetRowsAffected(
		expCtx,
		expQuery,
		expArgs...,
	)
	s.Require().NoError(err)
	s.Require().Equal(expRowsAffected, actualRowsAffected)

}

func (s *QueryerMixinSuite) TestUpdateAndGetLastInsertID() {
	sqlResultMock := NewSqlResultMock(s.T())

	expectedCtx := context.Background()
	expectedQuery := "UPDATE aTable SET x = ? * ?"
	expectedArgs := []interface{}{1, 94}
	expectedLastInsertID := int64(11)

	s.queryer.ExecMock.When(
		expectedCtx,
		expectedQuery,
		expectedArgs...,
	).Then(sqlResultMock, (error)(nil))

	sqlResultMock.LastInsertIdMock.Return(expectedLastInsertID, (error)(nil))

	actualLastInsertID, err := s.mixin.UpdateAndGetLastInsertID(expectedCtx, expectedQuery, expectedArgs...)
	s.Require().NoError(err)
	s.Require().Equal(expectedLastInsertID, actualLastInsertID)
}

func (s *QueryerMixinSuite) doTestScan(
	scan func(context.Context, RowScanner, string, ...interface{}) error,
	oneRow bool,
) {

	expRowScanner := NewRowScannerMock(s.T())
	defer expRowScanner.MinimockFinish()

	expSqlRows := NewSqlRowsMock(s.T())
	defer expSqlRows.MinimockFinish()

	expQuery := "SELECT something FROM somewhere WHERE x >= ? AND x < ?"
	expArgs := []interface{}{1, 94}
	expCtx := context.Background()

	s.queryer.QueryMock.When(
		expCtx, expQuery, expArgs...,
	).Then(expSqlRows, (error)(nil))

	s.scan.DoMock.Set(func(rowScanner RowScanner, oneRow bool, query func() (sqlRows, error)) (err error) {
		// execute query and and the assertion is that the call is
		// delegated to queryer
		_, _ = query()
		return (error)(nil)
	})

	err := scan(expCtx, expRowScanner, expQuery, expArgs...)
	s.Require().NoError(err)
}
