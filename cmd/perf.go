package cmd

import (
	"github.com/coderc/im/common/config"
	"github.com/coderc/im/perf"
	"github.com/spf13/cobra"
)

func init() {
	imCmd.AddCommand(perfCmd)
	perf.TcpConnNum = config.GetEnvInt32("tcp_conn_num", 100000, "tcp 最大连接数")
}

var (
	perfCmd = &cobra.Command{
		Use: "perf",
		Run: perfRun,
	}
)

func perfRun(cmd *cobra.Command, args []string) {
	perf.Run()
}
