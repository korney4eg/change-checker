package listen

import (
	"fmt"
	"log"
	"net/http"

	"github.com/korney4eg/change-checker/pkg/compare"
	"gopkg.in/go-playground/webhooks.v5/github"
)

type Command struct {
	Secret string `short:"s" long:"secret" required:"false" description:"Secret for Github webhook"`
	Config string `short:"c" long:"config" required:"true"  description:"path to configuration yaml file"`
	// Day           string `short:"o" long:"only-day" required:"false" description:"Get statistics only for provided date. Example '01.02.2020'"`
	// SplitPerYear  bool   `short:"y" long:"year-split" required:"true" description:"Will split files by year"`
	// SplitPerMonth bool   `short:"m" long:"month-split" required:"true" description:"Will split files by month"`
}

const (
	path = "/webhooks"
)

func (c *Command) Execute(_ []string) error {
	// var err error
	hook, _ := github.New(github.Options.Secret(c.Secret))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PushEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
			fmt.Println("Error:")
			fmt.Println(err)
		}
		switch payload.(type) {

		case github.PushPayload:
			push := payload.(github.PushPayload)
			err = compare.Run(c.Config, "push", push.Ref, push.Before, push.After)
			if err != nil {
				log.Fatal(err)
			}

			// case github.ReleasePayload:
			// 	release := payload.(github.ReleasePayload)
			// 	// Do whatever you want from here...
			// 	fmt.Printf("%+v", release)

			// case github.PullRequestPayload:
			// 	pullRequest := payload.(github.PullRequestPayload)
			// 	// Do whatever you want from here...
			// 	fmt.Printf("%+v", pullRequest)
			fmt.Println("Processing done.")
		}
	})
	http.ListenAndServe(":3000", nil)
	return nil
}
