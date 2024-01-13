package ipconfig

import (
	"github.com/coderc/im/ipconfig/domain"
	"github.com/coderc/im/ipconfig/source"
)

func Run() {
	source.Init()
	domain.Init()
}
