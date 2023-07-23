package main

import (
	"github.com/zer0go/netguard-client/cmd"
)

var Version = "development"

func main() {
	cmd.Execute(Version)
}
