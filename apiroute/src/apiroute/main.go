package main

import (
	_ "apiroute/routers"
	"os"
)

func main() {
	cli := NewCLI(os.Stdout, os.Stderr)
	os.Exit(cli.Run(os.Args))
}
