package main

import (
	"github.com/myoperator/k8saggregatorapiserver/cmd/server/app"
	"os"
)

func main() {
	cmd := app.NewServerCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
