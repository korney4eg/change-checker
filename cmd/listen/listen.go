package listen

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"

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

type Output struct {
	CommitBefore string
	CommitAfter  string
	Title        string
	Url          string
}

func (c *Command) Execute(_ []string) error {
	// var err error
	hook, _ := github.New(github.Options.Secret(c.Secret))
	cfg, err := compare.NewConfig(c.Config)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PushEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
			fmt.Println("Error:")
			fmt.Println(err)
		}
		buf := new(bytes.Buffer)
		generatedTpl := ""
		switch payload.(type) {

		case github.PushPayload:
			push := payload.(github.PushPayload)
			output := Output{push.Before, push.After, "", ""}
			for _, task := range cfg.Filter("push", strings.TrimLeft(push.Ref, "refs/heads/")) {
				changedItems, err := compare.Run(push.Before, push.After, &task)
				if err != nil {
					log.Fatal(err)
				}
				for _, item := range changedItems {

					output.Title = item.Title
					output.Url = "https://learningdevops.makvaz.com" + item.GUID
					tmpl, err := template.New("test").Parse(task.OutputTemplate)
					if err != nil {
						log.Fatal(err)
					}
					err = tmpl.Execute(buf, output)
					if err != nil {
						log.Fatal(err)
					}
					generatedTpl = buf.String()
					// fmt.Println(generatedTpl)

					cmd := exec.Command("/bin/sh", "-c", generatedTpl)
					stdoutStderr, err := cmd.CombinedOutput()
					fmt.Printf("%s\n", stdoutStderr)
					if err != nil {
						log.Fatal(err)
					}
				}

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
