package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

var rootCmd = &cobra.Command{
	Use:   "im",
	Short: "this is a IM project",
	Run:   IM,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panicln(err)
	}
}

func IM(cmd *cobra.Command, args []string) {
	log.Println("IM handle, do nothing...")
}

func initConfig() {

}
