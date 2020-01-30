package main

import (
	"github.com/davidpenn/sensu-plugins/sdk/checks"
	"github.com/davidpenn/sensu-plugins/sdk/sensu"
	"github.com/spf13/cobra"
)

var commands []func() *cobra.Command = []func() *cobra.Command{
	checks.NewMySQLReplicationStatusCommand,
}

func main() {
	root := &cobra.Command{
		Use: "checK",
		Run: func(*cobra.Command, []string) {
			sensu.Exit(sensu.ConfigError, "command required")
		},
	}
	for _, cmd := range commands {
		root.AddCommand(cmd())
	}
	if err := root.Execute(); err != nil {
		sensu.Exit(sensu.RuntimeError)
	}
}
