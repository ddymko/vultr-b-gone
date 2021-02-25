package objects

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	objectLong    = `This will delete object storages on your account. You can specify whether you want all objects or on in a specific region`
	objectExample = `

	# To delete all object storages on your account
	$ vultr-b-gone objects --method=all

	# To delete object storages only in specific regions
	$ vultr-b-gone objects --method=region --regions="ewr,sea"

	# To delete all object storages expect for specific ones
	$ vultr-b-gone objects --method=region --omit-regions="ewr,sea"
`
)

// NewCmdObject returns the instance cobra command
func NewCmdObject(config *util.VultrBGone, parentWait *sync.WaitGroup) *cobra.Command {
	options := &util.OptionsScheme{}

	command := &cobra.Command{
		Use:     "objects",
		Short:   "delete objects",
		Long:    objectLong,
		Example: objectExample,
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
	command.Flags().StringSliceP("regions", "r", []string{}, "list of region(s) that object storages will be deleted from. Must be provided with `method` of `region`")
	command.Flags().StringSliceP("omit-regions", "o", []string{}, "list of region(s) that object storages will be omitted from deletion. Must be provided with `method` of `region`")

	return command
}

func Run(config *util.VultrBGone, wg *sync.WaitGroup) {
	listOptions := &govultr.ListOptions{PerPage: 100}
		for {
		i, meta, err := config.Config.ObjectStorage.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.ObjectStorage) {
				if util.LocationCheck(v.Region) {
					if err := config.Config.ObjectStorage.Delete(context.Background(), v.ID); err != nil {
						fmt.Println("error : ", err.Error())
						defer wg.Done()
						return
					}
					fmt.Println("deleted object storage:", v.ID, " region:", v.Region)
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
