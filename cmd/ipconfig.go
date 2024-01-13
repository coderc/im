package cmd

import (
	"github.com/coderc/im/ipconfig"
	"github.com/spf13/cobra"
)

func init() {
	imCmd.AddCommand(ipConfigCmd)
}

var (
	ipConfigCmd = &cobra.Command{
		Use: "ipConfig",
		Run: ipConfigRun,
	}
)

func ipConfigRun(cmd *cobra.Command, args []string) {
	ipconfig.Run()
}
