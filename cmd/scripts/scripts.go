package scripts

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	scriptLong    = `This will delete scripts on your account.`
	scriptExample = `
	# To delete all scripts on your account
	$ vultr-b-gone scripts
`
)

// NewCmdScript returns the instance cobra command
func NewCmdScript(config *util.VultrBGone) *cobra.Command {
	return &cobra.Command{
		Use:     "scripts",
		Short:   "delete scripts",
		Long:    scriptLong,
		Example: scriptExample,
		Run: func(cmd *cobra.Command, args []string) {
			run(config)
		},
	}
}

func run(config *util.VultrBGone) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	wg := sync.WaitGroup{}
	for {
		i, meta, err := config.Config.StartupScript.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.StartupScript) {
				if err := config.Config.StartupScript.Delete(context.Background(), v.ID); err != nil {
					fmt.Println("error : ", err.Error())
					defer wg.Done()
					return
				}
				fmt.Println("deleted script:", v.ID)
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
