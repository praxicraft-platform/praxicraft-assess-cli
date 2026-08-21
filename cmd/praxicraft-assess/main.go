package main

import (
	"os"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdroot"
)

// Set via -ldflags "-X github.com/praxicraft-platform/praxicraft-assess-cli/internal/cmdroot.Version=..."
var version = "0.1.4"

func main() {
	if version != "" {
		cmdroot.Version = version
	}
	os.Exit(cmdroot.Execute())
}
