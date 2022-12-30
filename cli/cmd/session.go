package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/unweave/unweave-v2/api"
	"github.com/unweave/unweave-v2/cli/config"
	"github.com/unweave/unweave-v2/cli/ui"
	"github.com/unweave/unweave-v2/types"
)

func SessionCreate(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	uwc := InitUnweaveClient()
	sshKey := types.SSHKey{}

	if config.SSHKeyName != "" {
		sshKey.Name = &config.SSHKeyName
	} else if config.SSHKeyPath != "" {
		f, err := os.ReadFile(config.SSHKeyPath)
		if err != nil {
			return err
		}
		s := string(f)
		sshKey.PublicKey = &s
	} else {
		newKey := ui.Confirm("No SSH key path provided. Do you want to generate a new SSH key")
		if !newKey {
			fmt.Println("No SSH key path provided")
			return nil
		}
		// Leave the sshKey fields empty to generate a new key
	}

	params := api.SessionCreateParams{
		Runtime: types.LambdaLabsProvider,
		SSHKey:  sshKey,
	}
	_, err := uwc.Session.Create(cmd.Context(), params)
	if err != nil {
		return err
	}
	return nil
}

func SessionList(cmd *cobra.Command, args []string) error {
	return nil
}

func SessionTerminate(cmd *cobra.Command, args []string) error {
	return nil
}
