package instances

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
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

// NewCmdInstance returns the instance cobra command
func NewCmdInstance(config *util.VultrBGone) *cobra.Command {
	options := &util.OptionsScheme{}

	command := &cobra.Command{
		Use:     "instances",
		Short:   "delete instances",
		Long:    instanceLong,
		Example: instanceExample,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckError(options.SetOptions(cmd))
			util.CheckError(options.Validate())
			config.Options = options
			run(config)
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

func run(config *util.VultrBGone) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	wg := sync.WaitGroup{}
	for {
		i, meta, err := config.Config.Instance.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.Instance) {
				if util.LocationCheck(v.Region) {
					if err := config.Config.Instance.Delete(context.Background(), v.ID); err != nil {
						fmt.Println("error : ", err.Error())
						defer wg.Done()
						return
					}
					fmt.Println("deleted: ", v.ID, " region: ", v.Region)
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
