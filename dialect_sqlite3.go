package crud

import (
	"database/sql"
)

type SQLite3Dialect struct{}

func (SQLite3Dialect) Scan(rows *sql.Rows, args ...FieldBinder) error {
	return genericScan(rows, args...)
}

func (SQLite3Dialect) Insert(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) (int64, error) {
	return genericInsert(db, table, sqlIdFieldName, obj)
}

func (SQLite3Dialect) Update(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) error {
	return genericUpdate(db, table, sqlIdFieldName, obj)
}
