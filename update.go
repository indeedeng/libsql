package libsql

import "database/sql"

func rowsAffected(r sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func lastInsertID(r sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}
