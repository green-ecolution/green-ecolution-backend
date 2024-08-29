package utils

import (
	"path"
	"path/filepath"
	"runtime"
)

func P[T any](v T) *T {
	return &v
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), "../")
	return filepath.Dir(d)
}

func CompareAndUpdate[T comparable](old T, new *T) T {
  if new == nil {
    return old
  }
  if old != *new {
    return *new
  }
  return old
}
