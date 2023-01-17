package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/cli/config"
	"github.com/unweave/unweave/cli/ui"
)

func Login(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	uwc := InitUnweaveClient()
	code, err := uwc.Account.PairingTokenCreate(cmd.Context())
	if err != nil {
		var e *types.HTTPError
		if errors.As(err, &e) {
			uie := &ui.Error{HTTPError: e}
			fmt.Println(uie.Verbose())
			os.Exit(1)
			return nil
		}
		return err
	}

	authURL := config.UnweaveConfig.AppURL + "/auth/pair?code=" + code
	openBrowser := ui.Confirm("Do you want to open the browser to login?", "y")

	var openErr error
	if openBrowser {
		openErr = open.Run(authURL)
	}

	if !openBrowser || openErr != nil {
		fmt.Println("Open the following URL in your browser to login: ", authURL)
	}

	var token, email string
	sleep := time.Duration(2)
	timeout := 5 * time.Minute
	retryUntil := time.Now().Add(timeout)

	for {
		if time.Now().After(retryUntil) {
			fmt.Printf("Login timed out after %f minutes \n", timeout.Minutes())
			os.Exit(1)
			return nil
		}

		token, email, err = uwc.Account.PairingTokenExchange(cmd.Context(), code)
		if err != nil {
			var e *types.HTTPError
			if errors.As(err, &e) {
				if e.Code == http.StatusUnauthorized {
					time.Sleep(sleep * time.Second)
					continue
				}
				uie := &ui.Error{HTTPError: e}
				fmt.Println(uie.Verbose())
				os.Exit(1)
				return nil
			}
			return err
		}
		break
	}

	config.UnweaveConfig.User.Token = token
	if err := config.UnweaveConfig.Save(); err != nil {
		return err
	}

	fmt.Printf("Logged in as %q", email)
	return nil
}
