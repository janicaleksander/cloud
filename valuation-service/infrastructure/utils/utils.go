package utils

import "reflect"

func NameOfType(x any) string {
	t := reflect.TypeOf(x)
	return t.Name()
}
