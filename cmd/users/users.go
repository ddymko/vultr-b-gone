package users

import (
	"context"
	"fmt"
	"sync"

	"github.com/ddymko/vultr-b-gone/cmd/util"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
)

var (
	userLong    = `This will delete users on your account.`
	userExample = `
	# To delete all users on your account
	$ vultr-b-gone users
`
)

// NewCmdUser returns the instance cobra command
func NewCmdUser(config *util.VultrBGone, parentWait *sync.WaitGroup) *cobra.Command {
	return &cobra.Command{
		Use:     "users",
		Short:   "delete users",
		Long:    userLong,
		Example: userExample,
		Run: func(cmd *cobra.Command, args []string) {
			Run(config, parentWait)
		},
	}
}

func Run(config *util.VultrBGone, wg *sync.WaitGroup) {
	listOptions := &govultr.ListOptions{PerPage: 100}
	for {
		i, meta, err := config.Config.User.List(context.Background(), listOptions)
		if err != nil {
			_ = fmt.Errorf("error retrieving list %s", err.Error())
			return
		}
		wg.Add(len(i))

		for _, v := range i {
			go func(v govultr.User) {
				if err := config.Config.User.Delete(context.Background(), v.ID); err != nil {
					fmt.Println("error : ", err.Error())
					defer wg.Done()
					return
				}
				fmt.Println("deleted user:", v.ID)
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
