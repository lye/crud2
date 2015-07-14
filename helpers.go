package crud

import (
	"database/sql"
	"fmt"
	"reflect"
)

// DefaultDialect is the Dialect used by the package-level Scan/Insert/Update.
// It is provided as a convenience as most applications will likely only use
// one dialect at a time.
var DefaultDialect Dialect = SQLite3Dialect{}

// Scan is shorthand for DefaultDialect.Scan.
func Scan(rows *sql.Rows, args ...FieldBinder) error {
	return DefaultDialect.Scan(rows, args...)
}

func ScanAll(rows *sql.Rows, slicePtr interface{}) error {
	defer rows.Close()

	sliceVal := reflect.ValueOf(slicePtr).Elem()

	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("Argument to crud.ScanAll is not a slice")
	}

	elemType := sliceVal.Type().Elem()

	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("Argument to crud.ScanAll must be a slice of structs")
	}

	for rows.Next() {
		newVal := reflect.New(elemType)

		if er := Scan(rows, newVal.Interface().(FieldBinder)); er != nil {
			return er
		}

		sliceVal.Set(reflect.Append(sliceVal, newVal.Elem()))
	}

	return nil
}

// Insert is shorthand for DefaultDialect.Insert.
func Insert(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) (int64, error) {
	return DefaultDialect.Insert(db, table, sqlIdFieldName, obj)
}

// Update is shorthand for DefaultDialect.Update.
func Update(db DbIsh, table, sqlIdFieldName string, obj FieldEnumerator) error {
	return DefaultDialect.Update(db, table, sqlIdFieldName, obj)
}

func inflate(val interface{}) (er error) {
	if inflater, ok := val.(Inflater); ok {
		er = inflater.CrudInflate()
	}

	return
}

func deflate(val interface{}) (er error) {
	if deflater, ok := val.(Deflater); ok {
		er = deflater.CrudDeflate()
	}

	return
}
