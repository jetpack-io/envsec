package envcli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.jetpack.io/envsec"
)

const environmentFlagName = "environment"

func ListCmd(cmdCfg *CmdConfig) *cobra.Command {
	type envListCmdFlags struct {
		ShowValues bool
	}

	flags := envListCmdFlags{}

	command := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List all stored environment variables",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Populate the valid Environments
			envNames := []string{"DEV", "PROD"}
			// If a specific environment was set by the user, then just use that one.
			if cmd.Flags().Changed(environmentFlagName) {
				envNames = []string{cmdCfg.EnvId.EnvName}
			}

			// TODO: parallelize
			for _, envName := range envNames {
				envId := envsec.EnvId{
					OrgId:     cmdCfg.EnvId.OrgId,
					ProjectId: cmdCfg.EnvId.ProjectId,
					EnvName:   envName,
				}
				envVars, err := cmdCfg.Store.List(cmd.Context(), envId)
				if err != nil {
					return errors.WithStack(err)
				}

				err = printEnv(cmd, envId, envVars, flags.ShowValues)
				if err != nil {
					return errors.WithStack(err)
				}
			}
			return nil
		},
	}

	command.Flags().BoolVarP(
		&flags.ShowValues,
		"show",
		"s",
		false,
		"Display the value of each environment variable (secrets included)",
	)

	return command
}
