package main

import (
	"github.com/spf13/cobra"
	"github.com/unweave/unweave-v2/cli/cmd"
	"github.com/unweave/unweave-v2/cli/config"
)

var rootCmd = &cobra.Command{
	Use:   "unweave <command>",
	Short: "Create serverless sessions to train your ML models",
	Example: "unweave session create\n" +
		"unweave ssh --sync-fs <session-id>\n" +
		"unweave exec python train.py\n",
	Args:          cobra.MinimumNArgs(0),
	SilenceUsage:  false,
	SilenceErrors: false,
}

func init() {
	rootCmd.Version = Version
	rootCmd.Flags().BoolP("version", "v", false, "Get the version of current Unweave CLI")
	rootCmd.PersistentFlags().StringVarP(&config.UnweaveConfig.User.Token, "token", "t", "", "Use a specific token to authenticate - overrides login token")
	rootCmd.PersistentFlags().StringVarP(&Path, "path", "p", "", "Path to an Unweave project to run")

	// Session commands
	sessionCmd := &cobra.Command{
		Use:   "session",
		Short: "Manage Unweave sessions: create|ls|terminate",
		Args:  cobra.NoArgs,
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Unweave session",
		Args:  cobra.NoArgs,
		RunE:  cmd.SessionCreate,
	}
	createCmd.Flags().StringVarP(&SSHKeyPath, "ssh-key", "k", "", "Path to an SSH key to use for the session")
	sessionCmd.AddCommand(createCmd)

	sessionCmd.AddCommand(&cobra.Command{
		Use:   "ls",
		Short: "List all active Unweave sessions",
		Args:  cobra.NoArgs,
		RunE:  cmd.SessionList,
	})
	sessionCmd.AddCommand(&cobra.Command{
		Use:   "terminate <session-id>",
		Short: "Terminate an Unweave session",
		Args:  cobra.ExactArgs(1),
		RunE:  cmd.SessionTerminate,
	})
	rootCmd.AddCommand(sessionCmd)

	// Token
	tokenCmd := &cobra.Command{
		Use:   "token",
		Short: "Configure authentication tokens for the current user",
		Args:  cobra.NoArgs,
	}

	tokenCmd.AddCommand(&cobra.Command{
		Use:   "create-user-token",
		Short: "Create a new token",
		RunE:  cmd.CreateUserToken,
	})

	tokenCmd.AddCommand(&cobra.Command{
		Use:   "get-user-tokens",
		Short: "Get all tokens for the current user",
		RunE:  cmd.GetUserTokens,
	})
	rootCmd.AddCommand(tokenCmd)

}
