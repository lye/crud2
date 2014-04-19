### crud2

`crud2` is a drop-in replacement for [`crud`](https://github.com/lye/crud) that uses code generation instead of reflection to retrieve type metadata to make interaction with SQL databases easier. The `crudgen` utility parses all Go files in the current directory and extends their functionality to implement the `FieldEnumerator` and `FieldBinder` interfaces.

As an added bonus, `crud2` also supports an extensible layer for supporting different SQL markups (whereas the original `crud` didn't work on PostgreSQL).

Some of the original features are currently missing:

 * `ScanAll` is not implemented yet.
 * `,unix` times are neither handled nor parsed correctly.
 * Pointer values may not work correctly.
 * The documentation needs work.

These features will be forthcoming as they become required.
