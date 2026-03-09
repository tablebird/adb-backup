package base

import (
	"github.com/gin-gonic/gin"
)

type ContextReq struct {
	Context *gin.Context
}

func (req *ContextReq) ShouldBindJSON(c *gin.Context, target interface{}) error {
	req.Context = c
	return c.ShouldBindJSON(target)
}

func (req *ContextReq) Set(key any, value any) {
	ctx := req.Context
	if ctx != nil {
		ctx.Set(key, value)
	}
}

func (req *ContextReq) SetTypeKey(value any) {
	SetContextTypeKey(req.Context, value)
}
