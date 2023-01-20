package main

import (
	"fmt"
	"os"

	"github.com/muesli/reflow/wordwrap"
	"github.com/spf13/cobra"
	"github.com/unweave/unweave/cli/cmd"
	"github.com/unweave/unweave/cli/config"
	"github.com/unweave/unweave/cli/ui"
)

var (
	groupDev        = "dev"
	groupManagement = "management"

	rootCmd = &cobra.Command{
		Use:   "unweave <command>",
		Short: "Create serverless sessions to train your ML models",
		Example: "unweave session create\n" +
			"unweave ssh --sync-fs <session-id>\n" +
			"unweave exec python train.py",
		Args:          cobra.MinimumNArgs(0),
		SilenceUsage:  false,
		SilenceErrors: false,
	}
)

func init() {
	rootCmd.Version = ""
	rootCmd.Flags().BoolP("version", "v", false, "Get the version of current Unweave CLI")
	rootCmd.AddGroup(&cobra.Group{ID: groupDev, Title: "Dev:"})
	rootCmd.AddGroup(&cobra.Group{ID: groupManagement, Title: "Account Management:"})

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&config.AuthToken, "token", "t", "", "Use a specific token to authenticate - overrides login token")
	flags.StringVarP(&config.ProjectPath, "path", "p", "", "ProjectPath to an Unweave project to run")

	rootCmd.AddCommand(&cobra.Command{
		Use:     "config",
		Short:   "Show the current config",
		GroupID: groupDev,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.Config.String())
		},
	})

	linkCmd := &cobra.Command{
		Use:     "link",
		Short:   "Link your local directory to an Unweave project",
		GroupID: groupManagement,
		Example: "unweave link project-id",
		Args:    cobra.ExactArgs(1),
		RunE:    cmd.Link,
	}
	linkCmd.Flags().StringP("path", "p", "", "Path to the project directory")
	rootCmd.AddCommand(linkCmd)

	// Auth
	loginCmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to Unweave",
		GroupID: groupManagement,
		RunE:    cmd.Login,
	}
	rootCmd.AddCommand(loginCmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:    "logout",
		Short:  "Logout of Unweave",
		RunE:   cmd.Logout,
		Hidden: true,
	})

	// Provider commands
	providerCmd := &cobra.Command{
		Use:     "provider",
		Short:   "Manage providers",
		GroupID: groupManagement,
		Args:    cobra.NoArgs,
	}
	lsNodeType := &cobra.Command{
		Use:   "list-node-types <provider>",
		Short: "List node types available on a provider",
		Args:  cobra.ExactArgs(1),
		RunE:  cmd.ProviderListNodeTypes,
	}
	lsNodeType.Flags().BoolVarP(&config.All, "all", "a", false, "Including out of capacity node types")
	providerCmd.AddCommand(lsNodeType)
	rootCmd.AddCommand(providerCmd)

	// Session commands
	sessionCmd := &cobra.Command{
		Use:     "sessions",
		Short:   "Manage Unweave sessions: create | ls | terminate",
		GroupID: groupDev,
		Args:    cobra.NoArgs,
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Unweave session.",
		Long: wordwrap.String("Create a new Unweave session. If no region is provided,"+
			"the first available one will be selected.", ui.MaxOutputLineLength),
		Args: cobra.NoArgs,
		RunE: cmd.SessionCreate,
	}
	createCmd.Flags().StringVar(&config.Provider, "provider", "", "Provider to use")
	createCmd.Flags().StringVar(&config.NodeTypeID, "type", "", "Node type to use, eg. `gpu_1x_a100`")
	createCmd.Flags().StringVar(&config.NodeRegion, "region", "", "Region to use, eg. `us_west_2`")
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

	// SSH Key commands
	sshKeyCmd := &cobra.Command{
		Use:     "ssh-keys",
		Short:   "Manage Unweave SSH keys: add | ls",
		GroupID: groupDev,
		Args:    cobra.NoArgs,
	}
	sshKeyCmd.AddCommand(&cobra.Command{
		Use:   "add <public-key-path> [name]",
		Short: "Add a new SSH key to Unweave",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  cmd.SSHKeyAdd,
	})
	sshKeyCmd.AddCommand(&cobra.Command{
		Use:   "ls",
		Short: "List Unweave SSH keys",
		Args:  cobra.NoArgs,
		RunE:  cmd.SSHKeyList,
	})
	rootCmd.AddCommand(sshKeyCmd)
}

func main() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
