package main

import (
	"fmt"

	"github.com/Aykutfgoktas/orc/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
