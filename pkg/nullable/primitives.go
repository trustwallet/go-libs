package nullable

import (
	"fmt"
	"time"
)

func String(s string) *string {
	return &s
}

func Stringf(s string, args ...interface{}) *string {
	s = fmt.Sprintf(s, args...)
	return &s
}

func Int(i int) *int {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func Bool(b bool) *bool {
	return &b
}

func Time(t time.Time) *time.Time {
	return &t
}
