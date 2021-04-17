package main

import (
	"github.com/korney4eg/change-checker/cmd/check"
	"github.com/korney4eg/change-checker/cmd/listen"
	log "github.com/sirupsen/logrus"
	"github.com/umputun/go-flags"
)

type opts struct {
	Listen listen.Command `command:"listen" description:"listen for webhooks from GitHub"`
	Check  check.Command  `command:"check" description:"get changes of repository on needed files"`
}
type FilterCmd struct{}

func (f *FilterCmd) Execute(_ []string) error {
	return nil
}

func main() {
	o := opts{}
	if _, err := flags.Parse(&o); err != nil {
		log.Fatal(err)
	}
}
