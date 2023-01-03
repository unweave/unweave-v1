package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/cli/config"
	"github.com/unweave/unweave/cli/ui"
	"github.com/unweave/unweave/types"
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
		var e *api.HTTPError
		if errors.As(err, &e) {
			fmt.Println(e.Verbose())
			return nil
		}
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
