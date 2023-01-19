package cmd

import (
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/unweave/unweave/cli/config"
	"github.com/unweave/unweave/cli/ui"
)

func Link(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	projectID, err := uuid.Parse(args[0])
	if err != nil {
		ui.Errorf("Invalid project ID: %s", args[0])
		os.Exit(1)
	}

	uwc := InitUnweaveClient()
	project, err := uwc.Account.ProjectGet(cmd.Context(), projectID)
	if err != nil {
		return ui.HandleError(err)
	}

	if _, err = os.Stat(config.ProjectConfigPath); err == nil {
		ui.Errorf("Project already linked. Delete the '.unweave' directory to unlink.")
		os.Exit(1)
	}

	if err = config.WriteProjectConfig(project.ID); err != nil {
		ui.Errorf("Failed to write project config: %s", err)
		os.Exit(1)
	}
	if err = config.WriteEnvConfig(); err != nil {
		ui.Errorf("Failed to write environment config: %s", err)
		os.Exit(1)
	}
	if err = config.WriteGitignore(); err != nil {
		ui.Errorf("Failed to write .gitignore: %s", err)
		os.Exit(1)
	}

	ui.Successf("Project linked: %s", project.Name)
	return nil
}
