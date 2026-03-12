package utils

import (
	"reflect"
	"time"
)

func GetTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	return t.Name()
}

func Delay(d time.Duration) func() time.Duration {
	return func() time.Duration { return d }
}
