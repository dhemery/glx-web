# Test Archive with Archive Errors

This archive passes `glx validate` structural validation,
but produces archive-level errors.
Archive errors cause deserialization to fail,
so `glx-web` need not handle them.

Each file and ID that yields an error starts with `error-`.
