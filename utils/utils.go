package utils

import (
	"math/rand"
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

func RandomDelay(delayRange int) func() time.Duration {
	return func() time.Duration {
		return time.Second * time.Duration(rand.Intn(delayRange)+1)
	}
}
