// Package must contains generic functions for ensuring a function is called
// that returns a nil error.
package must

// Do panics if calling the function results in an error.
func Do(fn func() error) {
	if err := fn(); err != nil {
		panic(err)
	}
}

// Compute returns fn()'s first return value and panics if err is non-nil.
func Compute[T any](fn func() (T, error)) T {
	t, err := fn()
	if err != nil {
		panic(err)
	}
	return t
}

// Value returns the first value if the err is non-nil; otherwise, it panics.
func Value[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
