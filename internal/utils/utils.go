package utils

import (
	"encoding/json"
	"net/http"
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
	//nolint:dogsled
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), "../")
	return filepath.Dir(d)
}

// Helper function to decode JSON response
func ParseJSONResponse(body *http.Response, target any) error {
	defer body.Body.Close()
	return json.NewDecoder(body.Body).Decode(target)
}

func StringPtrToString(source *string) string {
	if source == nil {
		return ""
	}
	return *source
}
