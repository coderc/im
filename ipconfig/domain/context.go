package domain

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type IpConfigContext struct {
	Ctx       *context.Context
	AppCtx    *app.RequestContext
	ClientCtx *ClientContext
}

type ClientContext struct {
	IP string `json:"ip"`
}

func BuildIpConfigContext(c *context.Context, ctx *app.RequestContext) *IpConfigContext {
	ipConfigContext := &IpConfigContext{
		Ctx:       c,
		AppCtx:    ctx,
		ClientCtx: &ClientContext{},
	}
	return ipConfigContext
}
