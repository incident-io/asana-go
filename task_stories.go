package asana

import (
	"context"
	"fmt"
	"time"
)

// TaskStory is a compact representation of a story on a task
type TaskStory struct {
	Gid          string    `json:"gid"`
	ResourceType string    `json:"resource_type"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    struct {
		ID           string `json:"gid"`
		ResourceType string `json:"resource_type"`
		Name         string `json:"name"`
	} `json:"created_by"`
	ResourceSubtype string `json:"resource_subtype"`
	Text            string `json:"text"`
}

// TaskStories returns the compact records for all stories on a task
func (t *Task) TaskStories(ctx context.Context, client *Client, options ...*Options) ([]*TaskStory, *NextPage, error) {
	client.trace("Listing stories for task %s...\n", t.ID)
	var result []*TaskStory

	// Make the request
	nextPage, err := client.get(ctx, fmt.Sprintf("/tasks/%s/stories", t.ID), nil, &result, options...)
	return result, nextPage, err
}

// AllTaskStories repeatedly pages through all available stories for a task
func (t *Task) AllTaskStories(ctx context.Context, client *Client, options ...*Options) ([]*TaskStory, error) {
	var allTaskStories []*TaskStory
	nextPage := &NextPage{}

	var taskStories []*TaskStory
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  100,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		taskStories, nextPage, err = t.TaskStories(ctx, client, allOptions...)
		if err != nil {
			return nil, err
		}

		allTaskStories = append(allTaskStories, taskStories...)
	}
	return allTaskStories, nil
}
