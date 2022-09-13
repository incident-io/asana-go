package main

import (
	"context"
	"fmt"

	asana "github.com/SnoozeThis-org/asana-go"
)

func ListWorkspaces(ctx context.Context, c *asana.Client) error {
	// List workspaces
	workspaces, nextPage, err := c.Workspaces(ctx)
	if err != nil {
		return err
	}
	_ = nextPage

	for _, workspace := range workspaces {
		if workspace.IsOrganization {
			fmt.Printf("  Organization %s %q\n", workspace.ID, workspace.Name)
		} else {
			fmt.Printf("  Workspace %s %q\n", workspace.ID, workspace.Name)
		}
	}
	return nil
}

func ListProjects(ctx context.Context, client *asana.Client, w *asana.Workspace) error {
	// List projects
	projects, err := w.AllProjects(ctx, client, &asana.Options{
		Fields: []string{"name", "section_migration_status", "layout"},
	})
	if err != nil {
		return err
	}

	for _, project := range projects {
		fmt.Printf("  Project %s %q\n", project.ID, project.Name)
	}
	return nil
}

func ListTasks(ctx context.Context, client *asana.Client, p *asana.Project) error {
	// List projects
	tasks, nextPage, err := p.Tasks(ctx, client, asana.Fields(asana.Task{}))
	if err != nil {
		return err
	}
	_ = nextPage

	for _, task := range tasks {
		fmt.Printf("  Task %s %q (separator: %v)\n", task.ID, task.Name, task.IsRenderedAsSeparator)
	}
	return nil
}

func ListSections(ctx context.Context, client *asana.Client, p *asana.Project) error {
	// List sections
	sections, nextPage, err := p.Sections(ctx, client)
	if err != nil {
		return err
	}
	_ = nextPage

	for _, section := range sections {
		fmt.Printf("  Section %s %q\n", section.ID, section.Name)
	}
	return nil
}
