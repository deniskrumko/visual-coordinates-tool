package cmd_serve

import (
	"log"

	"github.com/deniskrumko/visual-coordinates-tool/app/api"
	"github.com/spf13/cobra"
)

var (
	configFile *string
	Cmd        = &cobra.Command{
		Use:   "serve",
		Short: "Run API server",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := api.RunServer(cmd.Context(), *configFile); err != nil {
				log.Fatalf(err.Error()) // nolint
			}
		},
	}
)

func init() {
	configFile = Cmd.Flags().StringP("config", "c", "", "path to config file")
}
