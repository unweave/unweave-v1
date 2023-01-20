package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unweave/unweave/cli/config"
	"github.com/unweave/unweave/cli/ui"
)

func Logout(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	config.Config.Unweave.User.Token = ""
	if err := config.Config.Unweave.Save(); err != nil {
		ui.Errorf("Failed to logout. Failed to save config: %s", err)
	}

	ui.Successf("Logged out")
	return nil
}
