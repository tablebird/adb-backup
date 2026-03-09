package validator

import (
	"adb-backup/internal/web/base"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("deviceId", deviceId)
		v.RegisterValidation("deviceIdConnect", deviceIdConnect)
	}
}

func getContext(fl validator.FieldLevel) Context {
	parent := fl.Parent()
	req, ok := parent.Interface().(base.ContextReq)
	if ok {
		return &req
	}
	ctxReq := parent.FieldByName("ContextReq")
	if ctxReq.IsValid() && ctxReq.Kind() == reflect.Struct {
		if req, ok := ctxReq.Interface().(base.ContextReq); ok {
			return &req
		}
	}
	return nil
}
