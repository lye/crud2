### crudgen

`crudgen` is a utility for `crud2` that parses all Go files in the current directory and emits a `z_crud.go` file which extends all `crud:`-tagged structs to implement both `FieldEnumerator` and `FieldBinder`.
