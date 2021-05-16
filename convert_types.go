package octopusenergy

import (
	"time"
)

// String returns a pointer to the string value passed in.
// This is a helper function when you need to provide pointers to
// optional fields in the input options object.
func String(v string) *string {
	return &v
}

// Int returns a pointer to the int value passed in.
// This is a helper function when you need to provide pointers to
// optional fields in the input options object.
func Int(v int) *int {
	return &v
}

// Bool returns a pointer to the bool value passed in.
// This is a helper function when you need to provide pointers to
// optional fields in the input options object.
func Bool(v bool) *bool {
	return &v
}

// Time returns a pointer to the time.Time value passed in.
// This is a helper function when you need to provide pointers to
// optional fields in the input options object.
func Time(v time.Time) *time.Time {
	return &v
}
