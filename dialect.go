package crud

import (
	"database/sql"
	"fmt"
	"strings"
)

// Dialect provides an SQL dialect abstraction, as pretty much every RDBS
// has a slightly different syntax for doing certain operations. The main
// pain point is that PostgreSQL requires a `RETURNING $field` clause after
// an INSERT statement in order to retrieve the last value of the primary
// key sequence (whereas MySQL/SQLite3 will return it by default).
//
// If more functionality is added, the requirements of the Dialect interface
// will likely grow.
type Dialect interface {
	Scan(rows *sql.Rows, args ...FieldBinder) error
	Insert(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) (int64, error)
	Update(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) error
}

func genericScan(rows *sql.Rows, args ...FieldBinder) error {
	columns, er := rows.Columns()
	if er != nil {
		return er
	}

	// Force the column names to all be lower-case. This pre-emptively works
	// around some databases that have case-insensitive column names.
	for i, column := range columns {
		columns[i] = strings.ToLower(column)
	}

	values := make([]interface{}, len(columns))

	for _, arg := range args {
		arg.BindFields(columns, values)
	}

	for i, value := range values {
		if value == nil {
			values[i] = new(interface{})
		}
	}

	if er := rows.Scan(values...); er != nil {
		return er
	}

	for _, arg := range args {
		if er := inflate(arg); er != nil {
			return er
		}
	}

	return nil
}

func genericInsert(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) (int64, error) {
	if er := deflate(obj); er != nil {
		return 0, er
	}

	objFields, objValues := obj.EnumerateFields()

	if len(objFields) != len(objValues) {
		return 0, ErrLengthMismatch
	}

	sqlFields := make([]string, 0, len(objFields))
	sqlValues := make([]interface{}, 0, len(objFields))
	placeholders := make([]string, 0, len(objFields))

	for i, field := range objFields {
		// If there's an id field, skip it so it can be automatically assigned.
		if field != sqlIdFieldName {
			sqlValues = append(sqlValues, objValues[i])
			sqlFields = append(sqlFields, field)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(sqlValues)))
		}
	}

	q := `
		INSERT INTO %s 
		(%s)
		VALUES (%s)
	`
	q = fmt.Sprintf(q, table, strings.Join(sqlFields, ", "), strings.Join(placeholders, ", "))

	res, er := db.Exec(q, sqlValues...)
	if er != nil {
		return 0, er
	}

	return res.LastInsertId()
}

func genericUpdate(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) error {
	if er := deflate(obj); er != nil {
		return er
	}

	objFields, objValues := obj.EnumerateFields()

	if len(objFields) != len(objValues) {
		return ErrLengthMismatch
	}

	sqlFields := make([]string, 0, len(objFields))
	sqlValues := make([]interface{}, 0, len(objFields))

	var idValue interface{} = nil

	for i, field := range objFields {
		// Yank the id field out of the SET expression so it can be used as a WHERE constraint.
		if field == sqlIdFieldName {
			idValue = objValues[i]

		} else {
			sqlValues = append(sqlValues, objValues[i])
			sqlFields = append(sqlFields, fmt.Sprintf("%s = $%d", field, len(sqlValues)))
		}
	}

	if idValue == nil {
		return ErrUnsetPKey
	}

	sqlValues = append(sqlValues, idValue)

	q := `
		UPDATE %s
		SET %s
		WHERE %s = $%d
	`
	q = fmt.Sprintf(q, table, strings.Join(sqlFields, ", "), sqlIdFieldName, len(sqlValues))

	_, er := db.Exec(q, sqlValues...)
	return er
}
