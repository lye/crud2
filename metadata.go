package crud

// FieldBinder allows structs to select which fields, by name, they would like
// bound to their members. It is used instead of the reflect package to build
// a slice of interface{} values to pass to sql.Rows.Scan.
type FieldBinder interface {
	// BindFields takes a slice of SQL column names forced to lowercase and a
	// slice of values. Each name corresponds to the value with the same index.
	// For each name that the implementation believes belongs to it, the
	// corresponding element in values should be set to a pointer-to-member.
	BindFields(names []string, values []interface{})
}

// Cloner is an optional optimization to avoid multiple calls to BindFields.
// When extracting multiple rows, objects extracted after the first can be
// copied over and the bindings from the first can be re-used. This exchanges
// N string compares with N word copies.
type Cloner interface {
	// Copy should return a new instance of this FieldBinder with all members
	// set to the current values. The old primitive members can be left as-is,
	// since they will be overwritten by the next call to BindFields. Non-primitive
	// members need to be fully copied.
	Clone() FieldBinder
}

// FieldEnumerator provides structs with a method of emitting all their field
// names and corresponding values, for both insertion and updates.
type FieldEnumerator interface {
	// EnumerateFields should return a slice of all SQL column names the
	// object would like to modify, and a slice of all corresponding values.
	EnumerateFields() ([]string, []interface{})
}
