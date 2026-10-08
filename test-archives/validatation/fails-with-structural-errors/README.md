# Test Archive with Structural Errors

This archive fails `glx validate` structural validation.
Structural errors cause deserialization to fail,
so `glx-web` need not handle them.

Each file and ID that yields an error starts with `error-`.
