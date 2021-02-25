package networks

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	networkLong    = `This will delete networks on your account. You can specify whether you want all networks or on in a specific region`
	networkExample = `

	# To delete all networks on your account
	$ vultr-b-gone networks --method=all

	# To delete networks only in specific regions
	$ vultr-b-gone networks --method=region --regions="ewr,sea"

	# To delete all instance expect for specific ones
	$ vultr-b-gone networks --method=region --omit-regions="ewr,sea"
`
)

// NewCmdNetwork returns the instance cobra command
func NewCmdNetwork(config *util.VultrBGone, parentWait *sync.WaitGroup) *cobra.Command {
	options := &util.OptionsScheme{}

	command := &cobra.Command{
		Use:     "networks",
		Short:   "delete networks",
		Long:    networkLong,
		Example: networkExample,
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
	command.Flags().StringSliceP("regions", "r", []string{}, "list of region(s) that networks will be deleted from. Must be provided with `method` of `region`")
	command.Flags().StringSliceP("omit-regions", "o", []string{}, "list of region(s) that networks will be omitted from deletion. Must be provided with `method` of `region`")

	return command
}

func Run(config *util.VultrBGone, wg *sync.WaitGroup) {
	listOptions := &govultr.ListOptions{PerPage: 100}
		for {
		i, meta, err := config.Config.Network.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.Network) {
				if util.LocationCheck(v.Region) {
					if err := config.Config.Network.Delete(context.Background(), v.NetworkID); err != nil {
						fmt.Println("error : ", err.Error())
						defer wg.Done()
						return
					}
					fmt.Println("deleted network:", v.NetworkID, " region:", v.Region)
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
