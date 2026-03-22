package utils

import (
	"reflect"
)

func NameOfType(msg any) string {
	t := reflect.TypeOf(msg)
	return t.Name()
}
