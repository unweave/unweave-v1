package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/cli/config"
	"github.com/unweave/unweave/cli/ui"
	"github.com/unweave/unweave/tools"
)

func nodeTypesToTable(nodeTypes []types.NodeType) ([]ui.Column, []ui.Row) {
	cols := []ui.Column{
		{Title: "Name", Width: 25},
		{Title: "ID", Width: 21},
		{Title: "Price", Width: 10},
		{Title: "Regions", Width: 50},
	}

	var rows []ui.Row

	for _, nodeType := range nodeTypes {
		regions := "-"
		if len(nodeType.Regions) > 0 {
			regions = strings.Join(nodeType.Regions, ", ")
		}
		row := ui.Row{
			fmt.Sprintf("%s", dashIfZeroValue(tools.StringInv(nodeType.Name))),
			fmt.Sprintf("%s", dashIfZeroValue(nodeType.ID)),
			fmt.Sprintf("$%2.2f", float32(tools.IntInv(nodeType.Price))/100),
			fmt.Sprintf("%s", regions),
		}
		rows = append(rows, row)
	}
	return cols, rows
}

func ProviderListNodeTypes(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	provider := types.RuntimeProvider(args[0])
	uwc := InitUnweaveClient()
	filterAvailable := !config.All

	res, err := uwc.Provider.ListNodeTypes(cmd.Context(), provider, filterAvailable)
	if err != nil {
		return ui.HandleError(err)
	}

	cols, rows := nodeTypesToTable(res)
	ui.Table("Available Instances", cols, rows)
	return nil
}
