package ipconfig

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coderc/im/ipconfig/domain"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func GetIpInfoList(c context.Context, ctx *app.RequestContext) {
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()

	// 构建客户请求信息
	ipConfigCtx := domain.BuildIpConfigContext(&c, ctx)
	// 进行IP调度
	eds := domain.Dispatch(ipConfigCtx)
	// 根据得分取top5返回

	ipConfigCtx.AppCtx.JSON(consts.StatusOK, packRes(top5EdnPoints(eds)))
}
