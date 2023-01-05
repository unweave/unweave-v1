package main

import (
	"fmt"
	"os"

	"github.com/muesli/reflow/wordwrap"
	"github.com/spf13/cobra"
	"github.com/unweave/unweave/cli/cmd"
	"github.com/unweave/unweave/cli/config"
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

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&config.AuthToken, "token", "t", "", "Use a specific token to authenticate - overrides login token")
	flags.StringVarP(&config.ProjectPath, "path", "p", "", "ProjectPath to an Unweave project to run")

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
		Use:   "create <node-type-id> [region]",
		Short: "Create a new Unweave session.",
		Long: wordwrap.String("Create a new Unweave session. If no region is provided,"+
			"the first available one will be selected.", 100),
		Args: cobra.RangeArgs(1, 2),
		RunE: cmd.SessionCreate,
	}
	createCmd.Flags().StringVarP(&config.SSHKeyName, "ssh-key", "k", "", "Name of the SSH key to use for the session")
	createCmd.Flags().StringVar(&config.SSHKeyPath, "ssh-key-path", "", "Absolute Path to the SSH public key to use")
	sessionCmd.AddCommand(createCmd)

	lsCmd := &cobra.Command{
		Use:   "ls",
		Short: "List active Unweave sessions",
		Long:  "List active Unweave sessions. To list all sessions, use the --all flag.",
		Args:  cobra.NoArgs,
		RunE:  cmd.SessionList,
	}
	lsCmd.Flags().BoolVarP(&config.All, "all", "a", false, "List all sessions")
	sessionCmd.AddCommand(lsCmd)

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
