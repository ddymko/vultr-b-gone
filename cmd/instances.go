package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	instanceLong    = `This will delete instances on your account. You can specify whether you want all instances or on in a specific region`
	instanceExample = `

	# To delete all instances on your account
	$ vultr-b-gone instances --method=all

	# To delete instances only in specific regions
	$ vultr-b-gone instances --method=region --regions="ewr,sea"

	# To delete all instance expect for specific ones
	$ vultr-b-gone instances --method=region --omit-regions="ewr,sea"
`
)

func newInstance(cobra *cobra.Command, instance *SubscriptionScheme) (*SubscriptionScheme, error) {
	method, err := cobra.Flags().GetString("method")
	if err != nil {
		return nil, err
	}

	regions, err := cobra.Flags().GetStringSlice("regions")
	if err != nil {
		return nil, err
	}

	omitRegions, err := cobra.Flags().GetStringSlice("omit-regions")
	if err != nil {
		return nil, err
	}

	instance.Regions = regions
	instance.Method = method
	instance.OmitRegions = omitRegions

	return instance, nil
}

// NewCmdInstance returns the instance cobra command
func NewCmdInstance(instance *SubscriptionScheme) *cobra.Command {
	command := &cobra.Command{
		Use:     "instances",
		Short:   "delete instances",
		Long:    instanceLong,
		Example: instanceExample,
		Args: func(cmd *cobra.Command, args []string) error {
			instance, _ = newInstance(cmd, instance)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := instance.Validate(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			run(instance)
		},
	}

	command.Flags().StringP("method", "m", "all", "method of delete: all (default) or region")
	if err := command.MarkFlagRequired("method"); err != nil {
		panic(err.Error())
	}

	command.Flags().StringSliceP("regions", "r", []string{}, "list of region(s) that instances will be deleted from. Must be provided with `method` of `region`")
	command.Flags().StringSliceP("omit-regions", "o", []string{}, "list of region(s) that instances will be omitted from deletion. Must be provided with `method` of `region`")

	return command
}

func run(instance *SubscriptionScheme) {
	fmt.Println(instance.Method)
	fmt.Println(instance.OmitRegions)
	fmt.Println(instance.Regions)
}
