package main

import (
	"fmt"
	"os"

	// "github.com/korney4eg/notifier/cmd/check"
	"github.com/korney4eg/notifier/cmd/listen"
	"github.com/umputun/go-flags"
)

type opts struct {
	Listen listen.Command `command:"listen" description:"listen for webhooks from GitHub"`
	// Check  check.Command  `command:"check" description:"get changes of repository on needed files"`
}
type FilterCmd struct{}

func (f *FilterCmd) Execute(_ []string) error {
	fmt.Println("Filtered ...")
	return nil
}

func main() {
	o := opts{}
	if _, err := flags.Parse(&o); err != nil {
		os.Exit(1)
	}
}
