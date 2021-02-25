package cmd

import (
	"fmt"
	"os"

	"github.com/ddymko/vultr-b-gone/cmd/baremetals"
	"github.com/ddymko/vultr-b-gone/cmd/blocks"
	"github.com/ddymko/vultr-b-gone/cmd/domains"
	"github.com/ddymko/vultr-b-gone/cmd/instances"
	"github.com/ddymko/vultr-b-gone/cmd/loadbalancers"
	"github.com/ddymko/vultr-b-gone/cmd/networks"
	"github.com/ddymko/vultr-b-gone/cmd/objects"
	"github.com/ddymko/vultr-b-gone/cmd/reservedips"
	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vultr-b-gone",
	Short: "vultr-b-gone to remove resources",
	Long:  `A quick and simple way to mass purge resources from your vultr account.`,
}

func init() {
	config := setup()
	rootCmd.AddCommand(baremetals.NewCmdBareMetal(config))
	rootCmd.AddCommand(blocks.NewCmdBlock(config))
	rootCmd.AddCommand(domains.NewCmdDomain(config))
	rootCmd.AddCommand(instances.NewCmdInstance(config))
	rootCmd.AddCommand(loadbalancers.NewCmdLoadBalancer(config))
	rootCmd.AddCommand(networks.NewCmdNetwork(config))
	rootCmd.AddCommand(objects.NewCmdObject(config))
	rootCmd.AddCommand(reservedips.NewCmdReservedIP(config))
	rootCmd.AddCommand(NewCmdSnapshot())
}

// Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setup() *util.VultrBGone {
	return util.NewVultrBGone()
}
