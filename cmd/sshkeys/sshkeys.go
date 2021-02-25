package sshkeys

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	sshkeyLong    = `This will delete sshkeys on your account.`
	sshkeyExample = `
	# To delete all sshkeys on your account
	$ vultr-b-gone sshkeys
`
)

// NewCmdSSHKey returns the instance cobra command
func NewCmdSSHKey(config *util.VultrBGone) *cobra.Command {
	return &cobra.Command{
		Use:     "sshkeys",
		Short:   "delete sshkeys",
		Long:    sshkeyLong,
		Example: sshkeyExample,
		Run: func(cmd *cobra.Command, args []string) {
			run(config)
		},
	}
}

func run(config *util.VultrBGone) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	wg := sync.WaitGroup{}
	for {
		i, meta, err := config.Config.SSHKey.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.SSHKey) {
				if err := config.Config.SSHKey.Delete(context.Background(), v.ID); err != nil {
					fmt.Println("error : ", err.Error())
					defer wg.Done()
					return
				}
				fmt.Println("deleted sshkey:", v.ID)
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
