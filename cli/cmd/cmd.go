package cmd

import (
	"github.com/unweave/unweave-v2/cli/config"
	"github.com/unweave/unweave-v2/client"
)

func InitUnweaveClient() *client.Client {
	return client.NewClient(
		client.Config{
			ApiURL: config.UnweaveConfig.ApiURL,
			Token:  config.AuthToken,
		})
}
