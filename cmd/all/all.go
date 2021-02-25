package all

import (
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
	"github.com/ddymko/vultr-b-gone/cmd/sshkeys"
	"github.com/ddymko/vultr-b-gone/cmd/users"
	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"sync"
)

var (
	allLong    = `This will delete everything on your account`
	allExample = `
	# To delete everything on your account
	$ vultr-b-gone all
`
)

// NewCmdBackup returns the instance cobra command
func NewCmdAll(config *util.VultrBGone, parentWait *sync.WaitGroup) *cobra.Command {
	return &cobra.Command{
		Use:     "all",
		Short:   "delete everything",
		Long:    allLong,
		Example: allExample,
		Run: func(cmd *cobra.Command, args []string) {
			Run(config, parentWait)
		},
	}
}

func Run(config *util.VultrBGone, wg *sync.WaitGroup) {
	wg2 := sync.WaitGroup{}
	wg2.Add(14)

	go func() {
		backups.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		baremetals.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		blocks.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		domains.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		firewalls.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		instances.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		isos.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		loadbalancers.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		networks.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		objects.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		reservedips.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		scripts.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		sshkeys.Run(config, wg)
		defer wg2.Done()
	}()

	go func() {
		users.Run(config, wg)
		defer wg2.Done()
	}()

	wg2.Wait()
}
