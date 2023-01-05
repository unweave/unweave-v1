package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/cli/config"
	"github.com/unweave/unweave/cli/ui"
	"github.com/unweave/unweave/types"
)

const defaultProjectID = "00000000-0000-0000-0000-000000000002"

func dashIfZeroValue(v interface{}) interface{} {
	if v == reflect.Zero(reflect.TypeOf(v)).Interface() {
		return "-"
	}
	return v
}

func SessionCreate(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	nodeID := args[0]
	uwc := InitUnweaveClient()
	sshKeyName := types.Stringy("")
	sshPublicKey := types.Stringy("")

	if config.SSHKeyName != "" {
		sshKeyName = &config.SSHKeyName
	} else if config.SSHKeyPath != "" {
		f, err := os.ReadFile(config.SSHKeyPath)
		if err != nil {
			return err
		}
		s := string(f)
		sshPublicKey = &s
	} else {
		newKey := ui.Confirm("No SSH key path provided. Do you want to generate a new SSH key")
		if !newKey {
			fmt.Println("No SSH key path provided")
			return nil
		}
		// Leave the sshKey fields empty to generate a new key
	}

	params := api.SessionCreateParams{
		Runtime:      types.LambdaLabsProvider,
		NodeTypeID:   nodeID,
		SSHKeyName:   sshKeyName,
		SSHPublicKey: sshPublicKey,
	}
	session, err := uwc.Session.Create(cmd.Context(), uuid.MustParse(defaultProjectID), params)
	if err != nil {
		var e *api.HTTPError
		if errors.As(err, &e) {

			// If error 503, it's mostly likely an out of capacity error. Try and marshal,
			// the error message into the list of available instances.
			if e.Code == 503 {
				var availableInstances []types.NodeType
				if err = json.Unmarshal([]byte(e.Suggestion), &availableInstances); err == nil {
					cols := []ui.Column{
						{Title: "Name", Width: 25},
						{Title: "ID", Width: 15},
						{Title: "Price", Width: 10},
						{Title: "Regions", Width: 50},
					}

					idx := map[string]int{}
					var rows []ui.Row

					for _, instance := range availableInstances {
						instance := instance

						specID := fmt.Sprintf("%s-%d-%d-%d", instance.ID, instance.Specs.VCPUs, instance.Specs.Memory, instance.Specs.GPUMemory)
						rowID, ok := idx[specID]
						if !ok {
							row := []string{
								fmt.Sprintf("%s", dashIfZeroValue(types.StringInv(instance.Name))),
								fmt.Sprintf("%s", dashIfZeroValue(instance.ID)),
								fmt.Sprintf("$%2.2f", float32(types.IntInv(instance.Price))/100),
								fmt.Sprintf("%s", dashIfZeroValue(types.StringInv(instance.Region))),
							}
							rows = append(rows, row)
							idx[specID] = len(rows) - 1
							continue
						}
						row := rows[rowID]
						ridx := len(cols) - 1
						if !strings.Contains(row[ridx], types.StringInv(instance.Region)) {
							row[ridx] = fmt.Sprintf("%s, %s", row[ridx], types.StringInv(instance.Region))
						}
					}

					fmt.Println(e.Short())
					fmt.Println()
					ui.Table("Available Instances", cols, rows)
				} else {
					fmt.Println(e.Verbose())
				}
			}
			return nil
		}
		return err
	}
	fmt.Printf("Session created: \n"+
		"  ID:      %s\n"+
		"  Status:  %s\n"+
		"  SSH Key: %s\n", session.ID, session.Status, session.SSHKey.Name,
	)
	return nil
}

func SessionList(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	uwc := InitUnweaveClient()
	listTerminated := config.All
	sessions, err := uwc.Session.List(cmd.Context(), uuid.MustParse(defaultProjectID), listTerminated)
	if err != nil {
		var e *api.HTTPError
		if errors.As(err, &e) {
			fmt.Println(e.Verbose())
			return nil
		}
		return err
	}

	fmt.Println("Sessions:")
	for _, s := range sessions {
		fmt.Printf("ID: %s, Status: %s", s.ID, s.Status)
	}
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
