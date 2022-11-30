package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"

	asana "github.com/incident-io/asana-go"
)

var options struct {
	Token string `long:"token" description:"Personal Access Token used to authorize access to the API" env:"ASANA_TOKEN" required:"true"`

	Workspace []string `long:"workspace" short:"w" description:"Workspace to access"`
	Project   []string `long:"project" short:"p" description:"Project to access"`
	Task      []string `long:"task" short:"t" description:"Task to access"`

	Attach     string `long:"attach" description:"Attach a file to a task"`
	AddSection string `long:"add-section" description:"Add a new section to a project"`

	Stories bool `long:"stories" description:"List stories for a task"`
	Clean   bool `long:"clean" description:"Clean all stories from a task"`

	Debug   bool   `short:"d" long:"debug" description:"Show debug information"`
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose output"`
}

func authenticate(req *http.Request) (*url.URL, error) {
	req.Header.Add("Authorization", "Bearer "+options.Token)
	return nil, nil
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if _, err := flags.Parse(&options); err != nil {
		return
	}

	// Create a client
	client := asana.NewClient(&http.Client{
		Transport: &http.Transport{
			Proxy: authenticate,
		},
	})
	if options.Debug {
		client.Debug = true
		client.DefaultOptions.Pretty = true
	}
	client.Verbose = options.Verbose
	client.DefaultOptions.Enable = []asana.Feature{asana.StringIDs, asana.NewSections, asana.NewTaskSubtypes}

	ctx := context.TODO()

	// Load a task object
	if options.Task == nil {

		// Load a project object
		if options.Project == nil {

			// Load a workspace object
			if options.Workspace == nil {
				check(ListWorkspaces(ctx, client))
				return
			}

			for _, w := range options.Workspace {
				workspace := &asana.Workspace{ID: w}
				check(ListProjects(ctx, client, workspace))
			}
			return
		}

		for _, p := range options.Project {
			project := &asana.Project{ID: p}

			if options.AddSection != "" {
				request := &asana.SectionBase{
					Name: options.AddSection,
				}

				_, err := project.CreateSection(ctx, client, request)
				check(err)
				return
			}

			fmtProject(ctx, client, project)
		}
		return
	}

	for _, t := range options.Task {
		task := &asana.Task{ID: t}
		check(task.Fetch(ctx, client))

		fmt.Printf("Task %s: %q\n", task.ID, task.Name)
		if options.Attach != "" {
			addAttachment(ctx, task, client)
			return
		}
		if options.Stories {
			listStories(ctx, task, client)
		}
		if options.Clean {
			cleanStories(ctx, task, client)
		}

		fmtTask(ctx, task, client)
	}
}

func listStories(ctx context.Context, task *asana.Task, client *asana.Client) {
	stories, _, _ := task.Stories(ctx, client)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	for _, s := range stories {
		fmt.Printf("Story %s (%s):\n", s.ID, s.CreatedBy.Name)
		check(enc.Encode(s))
	}
}

func cleanStories(ctx context.Context, task *asana.Task, client *asana.Client) {
	stories, _, _ := task.Stories(ctx, client)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	for _, s := range stories {
		check(s.Delete(ctx, client))
	}
}

func fmtProject(ctx context.Context, client *asana.Client, project *asana.Project) {
	fmt.Println("\nSections:")
	check(ListSections(ctx, client, project))
	fmt.Println("\nTasks:")
	check(ListTasks(ctx, client, project))
}

func fmtTask(ctx context.Context, task *asana.Task, client *asana.Client) {
	fmt.Printf("  Completed: %v\n", task.Completed)
	if task.Completed != nil && !*task.Completed {
		fmt.Printf("  Due: %s\n", task.DueAt)
	}
	if task.Notes != "" {
		fmt.Printf("  Notes: %q\n", task.Notes)
	}
	// Get subtasks
	subtasks, nextPage, err := task.Subtasks(ctx, client)
	check(err)
	_ = nextPage
	for _, subtask := range subtasks {
		fmt.Printf("  Subtask %s: %q\n", subtask.ID, subtask.Name)
	}
}

func addAttachment(ctx context.Context, task *asana.Task, client *asana.Client) {
	f, err := os.Open(options.Attach)
	check(err)
	defer f.Close()
	a, err := task.CreateAttachment(ctx, client, &asana.NewAttachment{
		Reader:      f,
		FileName:    f.Name(),
		ContentType: mime.TypeByExtension(filepath.Ext(f.Name())),
	})
	check(err)
	fmt.Printf("Attachment added: %+v", a)
}
