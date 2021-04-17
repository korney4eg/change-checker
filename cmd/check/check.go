package check

import (
	"fmt"

	"github.com/korney4eg/change-checker/pkg/compare"
)

type Command struct {
	BaseBranch     string `short:"b" long:"base-branch" default:"master" required:"false" description:"branch to which compare "`
	UpstreamBranch string `short:"u" long:"upstream-branch" required:"true" description:"branch to which compare "`
	Config         string `short:"c" long:"config" required:"true"  description:"path to configuration yaml file"`
}

func (c *Command) Execute(_ []string) error {
	fmt.Println("check")
	err := compare.Run(c.Config, "push", "master", c.BaseBranch, c.UpstreamBranch)
	if err != nil {
		return err
	}
	return nil
}
