package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/cli/ui"
)

func SSHKeyAdd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	publicKeyPath := args[0]
	name := filepath.Base(publicKeyPath)

	if len(args) == 2 {
		name = args[1]
	}

	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed reading public key file: %v", err)
	}

	ctx := cmd.Context()
	uwc := InitUnweaveClient()
	params := api.SSHKeyAddParams{
		Name:      &name,
		PublicKey: string(publicKey),
	}

	if err = uwc.SSHKey.Add(ctx, params); err != nil {
		var e *api.HTTPError
		if errors.As(err, &e) {
			uie := &ui.Error{HTTPError: e}
			fmt.Println(uie.Verbose())
			os.Exit(1)
		}
		return err
	}
	return nil
}
