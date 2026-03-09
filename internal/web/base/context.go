package base

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	ContextDeviceIdKey = "deviceId"

	ContextDeviceKey = "device"
)

func TypeKey[T any]() any {
	return reflect.TypeFor[T]()
}

func typeKeyOf(t any) any {
	return reflect.TypeOf(t)
}

func SetContextTypeKey(c *gin.Context, value any) {
	if c == nil {
		return
	}
	c.Set(typeKeyOf(value), value)
}
