package firewalls

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	firewallLong    = `This will delete firewalls on your account.`
	firewallExample = `
	# To delete all firewalls on your account
	$ vultr-b-gone firewalls
`
)

// NewCmdFirewall returns the instance cobra command
func NewCmdFirewall(config *util.VultrBGone, parentWait *sync.WaitGroup) *cobra.Command {
	return &cobra.Command{
		Use:     "firewalls",
		Short:   "delete firewalls",
		Long:    firewallLong,
		Example: firewallExample,
		Run: func(cmd *cobra.Command, args []string) {
			Run(config, parentWait)
		},
	}
}

func Run(config *util.VultrBGone, wg *sync.WaitGroup) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	for {
		i, meta, err := config.Config.FirewallGroup.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.FirewallGroup) {
				if err := config.Config.FirewallGroup.Delete(context.Background(), v.ID); err != nil {
					fmt.Println("error : ", err.Error())
					defer wg.Done()
					return
				}
				fmt.Println("deleted firewall:", v.ID)
				defer wg.Done()
				return

			}(v)
		}
		if meta.Links.Next == "" {
			break
		} else {
			listOptions.Cursor = meta.Links.Next
			continue
		}
	}
	wg.Wait()
}
