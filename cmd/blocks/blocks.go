package blocks

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	blocksLong    = `This will delete block storages on your account. You can specify whether you want all block storage or on in a specific region`
	blocksExample = `

	# To delete all block storages on your account
	$ vultr-b-gone blocks --method=all

	# To delete block storages only in specific regions
	$ vultr-b-gone blocks --method=region --regions="ewr,lax"

	# To delete all block storages expect for specific ones
	$ vultr-b-gone blocks --method=region --omit-regions="ewr"
`
)

// NewCmdBlock returns the blocks cobra command
func NewCmdBlock(config *util.VultrBGone, parentWait *sync.WaitGroup) *cobra.Command {
	options := &util.OptionsScheme{}

	command := &cobra.Command{
		Use:     "blocks",
		Short:   "delete blocks",
		Long:    blocksLong,
		Example: blocksExample,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckError(options.SetOptions(cmd))
			util.CheckError(options.Validate())
			config.Options = options
			Run(config, parentWait)
		},
	}

	command.Flags().StringP("method", "m", "all", "method of delete: all (default) or region")
	if err := command.MarkFlagRequired("method"); err != nil {
		panic(err.Error())
	}
	command.Flags().StringSliceP("regions", "r", []string{}, "list of region(s) that block storages will be deleted from. Must be provided with `method` of `region`")
	command.Flags().StringSliceP("omit-regions", "o", []string{}, "list of region(s) that block storages will be omitted from deletion. Must be provided with `method` of `region`")

	return command
}

func Run(config *util.VultrBGone, wg *sync.WaitGroup) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	for {
		i, meta, err := config.Config.BlockStorage.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.BlockStorage) {
				if util.LocationCheck(v.Region) {
					if err := config.Config.BlockStorage.Delete(context.Background(), v.ID); err != nil {
						fmt.Println("error : ", err.Error())
						defer wg.Done()
						return
					}
					fmt.Println("deleted block:", v.ID, " region:", v.Region)
					defer wg.Done()
					return
				}
				defer wg.Done()
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
