package domain

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type Context struct {
	Ctx       *context.Context
	AppCtx    *app.RequestContext
	ClientCtx *ClientContext
}

type ClientContext struct {
	IP string `json:"ip"`
}

func BuildIpConfContext(c *context.Context, ctx *app.RequestContext) *Context {
	return &Context{
		Ctx:       c,
		AppCtx:    ctx,
		ClientCtx: &ClientContext{},
	}
}
