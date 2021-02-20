package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vultr-b-gone",
	Short: "vultr-b-gone to remove resources",
	Long:  `A quick and simple way to mass purge resources from your vultr account.`,
}

func init() {
	subs := &SubscriptionScheme{}
	rootCmd.AddCommand(NewCmdInstance(subs))
	rootCmd.AddCommand(NewCmdSnapshot())
}

// Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
