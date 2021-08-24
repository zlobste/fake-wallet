package main

import (
	"github.com/zlobste/fake-wallet/internal/cli"
	"os"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}