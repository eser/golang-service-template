package main

import (
	"github.com/eser/go-service/cmd/manage/subcommands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "esimcli",
		Short: "eSIM CLI for eSIM management",
		Long:  `eSIM CLI provides various functionalities for eSIM management including reporting and administration.`,
	}

	rootCmd.AddCommand(subcommands.CmdHealthCheck())

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
