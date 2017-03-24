package libsql

import "context"

func newPreparerMixin(preparer sqlPreparer) Preparer {
	return preparerMixin{preparer: preparer, newStatement: newStatement}
}

type preparerMixin struct {
	preparer     sqlPreparer
	newStatement func(sqlStmt) Statement
}

var _ Preparer = (*preparerMixin)(nil)

// Prepared implements Preparer.Prepared
func (m preparerMixin) Prepared(ctx context.Context, sql string, work func(Statement) error) error {
	stmt, err := m.preparer.Prepare(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return work(m.newStatement(stmt))
}
