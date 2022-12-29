package main

import (
	"fmt"
	"os"

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
	rootCmd.Version = ""
	rootCmd.Flags().BoolP("version", "v", false, "Get the version of current Unweave CLI")
	rootCmd.PersistentFlags().StringVarP(&config.UnweaveConfig.User.Token, "token", "t", "", "Use a specific token to authenticate - overrides login token")
	rootCmd.PersistentFlags().StringVarP(&config.ProjectPath, "path", "p", "", "ProjectPath to an Unweave project to run")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "config",
		Short: "Show the current config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.UnweaveConfig.String())
		},
	})

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
	createCmd.Flags().StringVarP(&config.SSHKeyPath, "ssh-key", "k", "", "ProjectPath to an SSH key to use for the session")
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

	// Auth
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication tokens: create-user-token|get-user-tokens",
		Args:  cobra.NoArgs,
	}

	authCmd.AddCommand(&cobra.Command{
		Use:   "create-user-token",
		Short: "Create a new token",
		RunE:  cmd.CreateUserToken,
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "get-user-tokens",
		Short: "Get all tokens for the current user",
		RunE:  cmd.GetUserTokens,
	})
	rootCmd.AddCommand(authCmd)
}

func main() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
