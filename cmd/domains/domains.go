package domains

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	domainLong    = `This will delete domains on your account.`
	domainExample = `

	# To delete all domains on your account
	$ vultr-b-gone domains

`
)

// NewCmdDomain returns the instance cobra command
func NewCmdDomain(config *util.VultrBGone) *cobra.Command {
	return &cobra.Command{
		Use:     "domains",
		Short:   "delete domains",
		Long:    domainLong,
		Example: domainExample,
		Run: func(cmd *cobra.Command, args []string) {
			run(config)
		},
	}
}

func run(config *util.VultrBGone) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	wg := sync.WaitGroup{}
	for {
		i, meta, err := config.Config.Domain.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.Domain) {
				if err := config.Config.Domain.Delete(context.Background(), v.Domain); err != nil {
					fmt.Println("error : ", err.Error())
					defer wg.Done()
					return
				}
				fmt.Println("deleted domain:", v.Domain)
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
