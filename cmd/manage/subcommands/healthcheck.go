package subcommands

import (
	"context"
	"fmt"

	"github.com/eser/go-service/pkg/sample/adapters/appcontext"
	"github.com/spf13/cobra"
)

func CmdHealthCheck() *cobra.Command {
	healthCheckCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "healthcheck",
		Short: "Check the health of the database",
		Long:  `Check the health of the database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return execHealthCheck(cmd.Context())
		},
	}

	return healthCheckCmd
}

func execHealthCheck(ctx context.Context) error {
	_, err := appcontext.NewAppContext(ctx)
	if err != nil {
		return err //nolint:wrapcheck
	}

	fmt.Println("Health check passed") //nolint:forbidigo

	return nil
}
