package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := initAPICommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initAPICommand() *cobra.Command {
	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "Server Info API",
		Run: func(cmd *cobra.Command, args []string) {
			executeAPI(cmd.PersistentFlags().Lookup("addr").Value.String())
		},
	}
	apiCmd.PersistentFlags().StringP("addr", "a", ":8000", "Service Address")
	return apiCmd
}
