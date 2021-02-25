package snapshots

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	snapshotLong    = `This will delete snapshots on your account.`
	snapshotExample = `
	# To delete all snapshots on your account
	$ vultr-b-gone snapshots
`
)

// NewCmdSnapshot returns the instance cobra command
func NewCmdSnapshot(config *util.VultrBGone, parentWait *sync.WaitGroup) *cobra.Command {
	return &cobra.Command{
		Use:     "snapshots",
		Short:   "delete snapshots",
		Long:    snapshotLong,
		Example: snapshotExample,
		Run: func(cmd *cobra.Command, args []string) {
			Run(config, parentWait)
		},
	}
}

func Run(config *util.VultrBGone, wg *sync.WaitGroup) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	for {
		i, meta, err := config.Config.Snapshot.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.Snapshot) {
				if err := config.Config.Snapshot.Delete(context.Background(), v.ID); err != nil {
					fmt.Println("error : ", err.Error())
					defer wg.Done()
					return
				}
				fmt.Println("deleted snapshot:", v.ID)
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
