package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

var (
	imCmd = &cobra.Command{
		Use:   "im",
		Short: "这是一个im系统",
		Run:   imRun,
	}
)

// CmdExecute 命令行根入口
func CmdExecute() {
	if err := imCmd.Execute(); err != nil {
		log.Fatalf("imCmd Execute failed, err: %v", err)
	}
}

// imRun im系统的执行入口, 当命令执行为./im时之行此函数（入口函数）
func imRun(cmd *cobra.Command, args []string) {

}

func initConfig() {

}
