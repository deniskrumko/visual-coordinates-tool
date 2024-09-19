package cmd

import (
	"context"

	serve "github.com/deniskrumko/visual-coordinates-tool/cmd/serve"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Version: "0.0.1",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}
)

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(serve.Cmd)
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
