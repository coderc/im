package main

import (
	"log"
	"os"

	"github.com/coderc/im/cmd"
)

func main() {
	log.Println(os.Args)
	cmd.CmdExecute()
}
