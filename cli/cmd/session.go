package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/cli/config"
	"github.com/unweave/unweave/cli/ui"
	"github.com/unweave/unweave/types"
)

const defaultProjectID = "00000000-0000-0000-0000-000000000002"

func SessionCreate(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	uwc := InitUnweaveClient()
	sshKey := &types.SSHKey{}

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
	_, err := uwc.Session.Create(cmd.Context(), uuid.MustParse(defaultProjectID), params)
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
	cmd.SilenceUsage = true

	sessionID, err := uuid.Parse(args[0])
	if err != nil {
		fmt.Println("Invalid session ID")
		return nil
	}

	confirm := ui.Confirm(fmt.Sprintf("Are you sure you want to terminate session %q", sessionID))
	if !confirm {
		return nil
	}

	uwc := InitUnweaveClient()
	err = uwc.Session.Terminate(cmd.Context(), uuid.MustParse(defaultProjectID), sessionID)
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
