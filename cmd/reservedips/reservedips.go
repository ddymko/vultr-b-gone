package reservedips

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	reservedLong    = `This will delete reservedips on your account. You can specify whether you want all reservedips or on in a specific region`
	reservedExample = `

	# To delete all reservedips on your account
	$ vultr-b-gone reservedips --method=all

	# To delete reservedips only in specific regions
	$ vultr-b-gone reservedips --method=region --regions="ewr,sea"

	# To delete all instance expect for specific ones
	$ vultr-b-gone reservedips --method=region --omit-regions="ewr,sea"
`
)

// NewCmdReservedIP returns the instance cobra command
func NewCmdReservedIP(config *util.VultrBGone) *cobra.Command {
	options := &util.OptionsScheme{}

	command := &cobra.Command{
		Use:     "reservedips",
		Short:   "delete reservedips",
		Long:    reservedLong,
		Example: reservedExample,
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
	command.Flags().StringSliceP("regions", "r", []string{}, "list of region(s) that reservedips will be deleted from. Must be provided with `method` of `region`")
	command.Flags().StringSliceP("omit-regions", "o", []string{}, "list of region(s) that reservedips will be omitted from deletion. Must be provided with `method` of `region`")

	return command
}

func run(config *util.VultrBGone) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	wg := sync.WaitGroup{}
	for {
		i, meta, err := config.Config.ReservedIP.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.ReservedIP) {
				if util.LocationCheck(v.Region) {
					if err := config.Config.ReservedIP.Delete(context.Background(), v.ID); err != nil {
						fmt.Println("error : ", err.Error())
						defer wg.Done()
						return
					}
					fmt.Println("deleted reserved ip:", v.ID, " region:", v.Region)
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
