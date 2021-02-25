package baremetals

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	baremetalLong    = `This will delete baremetals on your account. You can specify whether you want all baremetals or on in a specific region`
	baremetalExample = `

	# To delete all baremetals on your account
	$ vultr-b-gone baremetals --method=all

	# To delete baremetals only in specific regions
	$ vultr-b-gone baremetals --method=region --regions="ewr,sea"

	# To delete all instance expect for specific ones
	$ vultr-b-gone baremetals --method=region --omit-regions="ewr,sea"
`
)

// NewCmdBareMetal returns the instance cobra command
func NewCmdBareMetal(config *util.VultrBGone, parentWait *sync.WaitGroup) *cobra.Command {
	options := &util.OptionsScheme{}

	command := &cobra.Command{
		Use:     "baremetals",
		Short:   "delete baremetals",
		Long:    baremetalLong,
		Example: baremetalExample,
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
	command.Flags().StringSliceP("regions", "r", []string{}, "list of region(s) that baremetals will be deleted from. Must be provided with `method` of `region`")
	command.Flags().StringSliceP("omit-regions", "o", []string{}, "list of region(s) that baremetals will be omitted from deletion. Must be provided with `method` of `region`")

	return command
}

func Run(config *util.VultrBGone, wg *sync.WaitGroup) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	for {
		i, meta, err := config.Config.BareMetalServer.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.BareMetalServer) {
				if util.LocationCheck(v.Region) {
					if err := config.Config.BareMetalServer.Delete(context.Background(), v.ID); err != nil {
						fmt.Println("error : ", err.Error())
						defer wg.Done()
						return
					}
					fmt.Println("deleted instance:", v.ID, " region:", v.Region)
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
