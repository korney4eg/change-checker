package compare

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v2"
)

type Output struct {
	CommitBefore string
	CommitAfter  string
	Title        string
	Url          string
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
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", stdoutStderr)
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
func Run(config, action, branch, commitBefore, commitAfter string) error {
	cfg, err := NewConfig(config)
	buf := new(bytes.Buffer)
	generatedTpl := ""
	output := Output{commitBefore, commitAfter, "", ""}
	if err != nil {
		return err
	}
	for _, task := range cfg.Filter(action, strings.TrimLeft(branch, "refs/heads/")) {
		changedItems, err := GetItems(commitBefore, commitAfter, &task)
		if err != nil {
			return err
		}
		for _, item := range changedItems {

			output.Title = item.Title
			output.Url = "https://learningdevops.makvaz.com" + item.GUID
			tmpl, err := template.New("test").Parse(task.OutputTemplate)
			if err != nil {
				return err
			}
			err = tmpl.Execute(buf, output)
			if err != nil {
				return err
			}
			generatedTpl = buf.String()

			cmd := exec.Command("/bin/sh", "-c", generatedTpl)
			stdoutStderr, err := cmd.CombinedOutput()
			fmt.Printf("%s\n", stdoutStderr)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func GetItems(commitBefore, commitAfter string, task *Task) (items []*gofeed.Item, err error) {
	feedBefore, err := getFileByCommit(task.FileName, commitBefore, task.Command)
	if err != nil {
		return items, err
	}
	feedAfter, err := getFileByCommit(task.FileName, commitAfter, task.Command)
	if err != nil {
		return items, err
	}
	foundItem := false
	for _, item := range feedAfter.Items {
		foundItem = false
		for _, itemBefore := range feedBefore.Items {
			if item.Title == itemBefore.Title {
				foundItem = true
			}
		}
		if !foundItem {
			items = append(items, item)
		}
	}
	return items, nil

}

func find(element string, branches []string) bool {
	for _, branch := range branches {
		if branch == element {
			return true
		}
	}
	return false
}
