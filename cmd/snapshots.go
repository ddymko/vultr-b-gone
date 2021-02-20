package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	snapshotLong    = `This will delete all snapshots on your account`
	snapshotExample = `

	# To delete all snapshots on your account
	$ vultr-b-gone snapshots
`
)

// NewCmdSnapshot returns the instance cobra command
func NewCmdSnapshot() *cobra.Command {
	command := &cobra.Command{
		Use:     "snapshots",
		Short:   "delete snapshots",
		Long:    snapshotLong,
		Example: snapshotExample,
		//Args: func(cmd *cobra.Command, args []string) error {
		//	instance, _ = newInstance(cmd, instance)
		//	return nil
		//},
		Run: func(cmd *cobra.Command, args []string) {
			snapshotRun()
		},
	}
	return command
}

func snapshotRun() {
	fmt.Println("snapshots")
}