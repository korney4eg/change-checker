package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"net/http"

	"github.com/korney4eg/notifier/cmd/listen"
	"github.com/mmcdole/gofeed"
	"github.com/umputun/go-flags"

	"gopkg.in/go-playground/webhooks.v5/github"
	"gopkg.in/yaml.v2"
)

type opts struct {
	listner listen.Command `command:"listen"`
}

type Config struct {
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Action         string   `yaml:"action"`
	Command        string   `yaml:"command"`
	FileName       string   `yaml:"file_name"`
	OutputTemplate string   `yaml:"output_template"`
	OnlyBranches   []string `yaml:"only_branches"`
}

type Output struct {
	CommitBefore string
	CommitAfter  string
	Title        string
	Url          string
}

const (
	path = "/webhooks"
)

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func (config *Config) Filter(action string, branch string) (tasks []Task) {
	for _, task := range config.Tasks {
		if task.Action == action && find(branch, task.OnlyBranches) {
			tasks = append(tasks, task)
		}
	}
	return tasks

}

func getFileByCommit(fileName, commit, cmdStr string) (*gofeed.Feed, error) {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf(cmdStr, commit))
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)
	return feed, err
}
func find(element string, branches []string) bool {
	for _, branch := range branches {
		if branch == element {
			return true
		}
	}
	return false
}

func run(commitBefore, commitAfter string, task *Task) {
	feedBefore, err := getFileByCommit(task.FileName, commitBefore, task.Command)
	if err != nil {
		log.Fatal(err)
	}
	feedAfter, err := getFileByCommit(task.FileName, commitAfter, task.Command)
	if err != nil {
		log.Fatal(err)
	}
	output := Output{commitBefore, commitAfter, "", ""}
	foundItem := false
	buf := new(bytes.Buffer)
	generatedTpl := ""
	for _, item := range feedAfter.Items {
		foundItem = false
		for _, itemBefore := range feedBefore.Items {
			if item.Title == itemBefore.Title {
				foundItem = true
			}
		}
		if !foundItem {
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

}

func main() {
	o := opts{}
	if _, err := flags.Parse(&o); err != nil {
		os.Exit(1)
	}
	// hook, _ := github.New(github.Options.Secret("myS3cret1@"))
	hook, _ := github.New()
	cfg, err := NewConfig("config.yaml")
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
		switch payload.(type) {

		case github.PushPayload:
			push := payload.(github.PushPayload)
			for _, task := range cfg.Filter("push", strings.TrimLeft(push.Ref, "refs/heads/")) {
				run(push.Before, push.After, &task)
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
}
