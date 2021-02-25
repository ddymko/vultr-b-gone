package cmd

import (
	"fmt"
	"github.com/ddymko/vultr-b-gone/cmd/all"
	"os"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/backups"
	"github.com/ddymko/vultr-b-gone/cmd/baremetals"
	"github.com/ddymko/vultr-b-gone/cmd/blocks"
	"github.com/ddymko/vultr-b-gone/cmd/domains"
	"github.com/ddymko/vultr-b-gone/cmd/firewalls"
	"github.com/ddymko/vultr-b-gone/cmd/instances"
	"github.com/ddymko/vultr-b-gone/cmd/isos"
	"github.com/ddymko/vultr-b-gone/cmd/loadbalancers"
	"github.com/ddymko/vultr-b-gone/cmd/networks"
	"github.com/ddymko/vultr-b-gone/cmd/objects"
	"github.com/ddymko/vultr-b-gone/cmd/reservedips"
	"github.com/ddymko/vultr-b-gone/cmd/scripts"
	"github.com/ddymko/vultr-b-gone/cmd/snapshots"
	"github.com/ddymko/vultr-b-gone/cmd/sshkeys"
	"github.com/ddymko/vultr-b-gone/cmd/users"
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
	wg := &sync.WaitGroup{}
	rootCmd.AddCommand(all.NewCmdAll(config, wg))
	rootCmd.AddCommand(backups.NewCmdBackup(config, wg))
	rootCmd.AddCommand(baremetals.NewCmdBareMetal(config, wg))
	rootCmd.AddCommand(blocks.NewCmdBlock(config, wg))
	rootCmd.AddCommand(domains.NewCmdDomain(config, wg))
	rootCmd.AddCommand(firewalls.NewCmdFirewall(config, wg))
	rootCmd.AddCommand(instances.NewCmdInstance(config, wg))
	rootCmd.AddCommand(isos.NewCmdISO(config, wg))
	rootCmd.AddCommand(loadbalancers.NewCmdLoadBalancer(config, wg))
	rootCmd.AddCommand(networks.NewCmdNetwork(config, wg))
	rootCmd.AddCommand(objects.NewCmdObject(config, wg))
	rootCmd.AddCommand(reservedips.NewCmdReservedIP(config, wg))
	rootCmd.AddCommand(scripts.NewCmdScript(config, wg))
	rootCmd.AddCommand(snapshots.NewCmdSnapshot(config, wg))
	rootCmd.AddCommand(sshkeys.NewCmdSSHKey(config, wg))
	rootCmd.AddCommand(users.NewCmdUser(config, wg))
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
