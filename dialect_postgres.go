package crud

import (
	"database/sql"
	"fmt"
	"strings"
)

type PostgresDialect struct{}

func (PostgresDialect) Scan(rows *sql.Rows, args ...FieldBinder) error {
	return genericScan(rows, args...)
}

func (PostgresDialect) Insert(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) (id int64, er error) {
	if er := deflate(obj); er != nil {
		return 0, er
	}

	objFields, objValues := obj.EnumerateFields()

	if len(objFields) != len(objValues) {
		panic("crud2: FieldEnumerator.EnumerateFields' return values must have same length")
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

	var q string

	if sqlIdFieldName != "" {
		q = `
			INSERT INTO %s 
			(%s)
			VALUES (%s)
			RETURNING %s
		`
		q = fmt.Sprintf(q, table, strings.Join(sqlFields, ", "), strings.Join(placeholders, ", "), sqlIdFieldName)

		rows, er := db.Query(q, sqlValues...)
		if er != nil {
			return 0, er
		}
		defer rows.Close()

		rows.Next()
		er = rows.Scan(&id)

	} else {
		q = `
			INSERT INTO %s 
			(%s)
			VALUES (%s)
		`
		q = fmt.Sprintf(q, table, strings.Join(sqlFields, ", "), strings.Join(placeholders, ", "))

		_, er = db.Exec(q, sqlValues...)
	}

	return
}

func (PostgresDialect) Update(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) error {
	return genericUpdate(db, table, sqlIdFieldName, obj)
}
