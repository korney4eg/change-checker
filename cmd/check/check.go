package check

import (
	"fmt"

	"github.com/korney4eg/change-checker/pkg/compare"
)

type Command struct {
	BaseBranch     string `short:"b" long:"base-branch" default:"master" required:"false" description:"branch to which compare "`
	UpstreamBranch string `short:"u" long:"upstream-branch" required:"true" description:"branch to which compare "`
}

func (c *Command) Execute(_ []string) error {
	fmt.Println("check")
	changedItems, err := compare.Run(c.BaseBranch, c.UpstreamBranch, &task)
	if err != nil {
		return err
	}
	for _, item := range changedItems {
	}
	return nil
}
