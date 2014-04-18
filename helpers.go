package crud

import (
	"database/sql"
)

// DefaultDialect is the Dialect used by the package-level Scan/Insert/Update.
// It is provided as a convenience as most applications will likely only use
// one dialect at a time.
var DefaultDialect Dialect = SQLite3Dialect{}

// Scan is shorthand for DefaultDialect.Scan.
func Scan(rows *sql.Rows, args ...FieldBinder) error {
	return DefaultDialect.Scan(rows, args...)
}

// Insert is shorthand for DefaultDialect.Insert.
func Insert(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) (int64, error) {
	return DefaultDialect.Insert(db, table, sqlIdFieldName, obj)
}

// Update is shorthand for DefaultDialect.Update.
func Update(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) error {
	return DefaultDialect.Update(db, table, sqlIdFieldName, obj)
}
