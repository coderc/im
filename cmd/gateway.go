package cmd

import (
	"github.com/coderc/im/gateway"
	"github.com/spf13/cobra"
)

func init() {
	imCmd.AddCommand(gatewayCmd)
}

var (
	gatewayCmd = &cobra.Command{
		Use: "gateway",
		Run: gatewayRun,
	}
)

func gatewayRun(cmd *cobra.Command, args []string) {
	gateway.Run()
}
