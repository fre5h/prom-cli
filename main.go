package main

import (
	"fmt"
	"os"

	"github.com/fre5h/prom-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("error on command execution: ", err)
		os.Exit(1)
	}
}
