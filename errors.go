package crud

import (
	"errors"
)

var (
	ErrLengthMismatch = errors.New("crud2: FieldEnumerator.EnumerateFields' return values must have same length")
	ErrUnsetPKey      = errors.New("crud2: FieldEnumerator.EnumerateFields did not return a field that matched sqlIdFieldName")
)
