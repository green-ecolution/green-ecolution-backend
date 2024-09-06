package utils

import (
	"path"
	"path/filepath"
	"runtime"
)

// P returns a pointer to the value passed as an argument.
func P[T any](v T) *T {
	return &v
}

// RootDir returns the root directory of the project.
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), "../")
	return filepath.Dir(d)
}

// CompareAndUpdate compares two values and updates the new value if it is different from the old value. If the new value is nil, the old value is returned. If the old value is different from the new value, the new value is returned. Otherwise, the old value is returned.
func CompareAndUpdate[T comparable](old T, new *T) T {
	if new == nil {
		return old
	}
	if old != *new {
		return *new
	}
	return old
}
