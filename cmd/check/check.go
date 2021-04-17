package check

import (
	"github.com/korney4eg/change-checker/pkg/compare"
	log "github.com/sirupsen/logrus"
)

type Command struct {
	BaseBranch     string `short:"b" long:"base-branch" default:"master" required:"false" description:"branch to which compare "`
	UpstreamBranch string `short:"u" long:"upstream-branch" required:"true" description:"branch to which compare "`
	Config         string `short:"c" long:"config" required:"true"  description:"path to configuration yaml file"`
	Debug          bool   `short:"d" long:"debug"`
}

func (c *Command) Execute(_ []string) error {
	if c.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%+v\n", c)
	err := compare.Run(c.Config, "push", "master", c.BaseBranch, c.UpstreamBranch, c.Debug)
	if err != nil {
		return err
	}
	return nil
}
