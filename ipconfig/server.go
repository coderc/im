package ipconfig

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/coderc/im/ipconfig/domain"
	"github.com/coderc/im/ipconfig/source"
)

func Run() {
	// 初始化数据源
	source.Init()
	// 初始化调度层
	domain.Init()
	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", GetIpInfoList)
	s.Spin()
}
