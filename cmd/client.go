package cmd

import (
	"github.com/coderc/im/client"
	"github.com/spf13/cobra"
)

func init() {
	// 将client命令执行入口注册到根执行入口中
	// 当命令执行为 ./im client时将clientRun作为程序入口
	imCmd.AddCommand(clientCmd)
}

var (
	clientCmd = &cobra.Command{
		Use: "client",
		Run: clientRun,
	}
)

func clientRun(cmd *cobra.Command, args []string) {
	client.Run()
}
