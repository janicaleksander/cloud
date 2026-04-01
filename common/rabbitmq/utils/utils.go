package utils

import "reflect"

func NameOfType(msg any) string {
	t := reflect.TypeOf(msg)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
