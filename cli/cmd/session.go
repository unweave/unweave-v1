package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unweave/unweave-v2/api"
	"github.com/unweave/unweave-v2/client"
	"github.com/unweave/unweave-v2/session/model"
)

func SessionCreate(cmd *cobra.Command, args []string) error {
	uwc := client.NewClient(client.Config{
		ApiUrl: "",
		Token:  "",
	})

	params := api.SessionCreateParams{
		Runtime: "",
		SSHKey:  model.SSHKey{},
	}
	_, err := uwc.Session.Create(cmd.Context(), params)
	if err != nil {

	}
	return nil
}

func SessionList(cmd *cobra.Command, args []string) error {
	return nil
}

func SessionTerminate(cmd *cobra.Command, args []string) error {
	return nil
}
