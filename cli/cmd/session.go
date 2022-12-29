package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/unweave/unweave-v2/api"
	"github.com/unweave/unweave-v2/cli/config"
	"github.com/unweave/unweave-v2/session/model"
)

func SessionCreate(cmd *cobra.Command, args []string) error {
	uwc := InitUnweaveClient()

	if config.SSHKeyPath == "" {
		log.Fatal("SSH key path not set")
	}
	// load the ssh key

	params := api.SessionCreateParams{
		Runtime: model.LambdaLabsProvider,
		SSHKey: model.SSHKey{
			Name:       "",
			PrivateKey: "",
			PublicKey:  "",
		},
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
