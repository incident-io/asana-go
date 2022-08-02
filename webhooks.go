package asana

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Filter struct {
	Action          string   `json:"action"`
	Fields          []string `json:"fields,omitempty"`
	ResourceType    string   `json:"resource_type"`
	ResourceSubtype string   `json:"resource_subtype,omitempty"`
}

type Webhook struct {
	ID           string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Active       bool   `json:"active"`
	Resource     struct {
		ID           string `json:"gid"`
		ResourceType string `json:"resource_type"`
		Name         string `json:"name"`
	} `json:"resource"`
	Target             string    `json:"target"`
	CreatedAt          time.Time `json:"created_at"`
	Filters            []Filter  `json:"filters"`
	LastFailureAt      time.Time `json:"last_failure_at"`
	LastFailureContent string    `json:"last_failure_content"`
	LastSuccessAt      time.Time `json:"last_success_at"`
}

// Webhooks returns the compact records for all webhooks your app has registered for the authenticated user in the given workspace.
func (w *Workspace) Webhooks(client *Client, options ...*Options) ([]*Webhook, *NextPage, error) {
	client.trace("Listing webhooks for workspace %s...\n", w.ID)
	var result []*Webhook

	workspace := &Options{
		Workspace: w.ID,
	}

	allOptions := append([]*Options{workspace}, options...)

	// Make the request
	nextPage, err := client.get("/webhooks", nil, &result, allOptions...)
	return result, nextPage, err
}

// AllWebhooks repeatedly pages through all available webhooks for a user
func (w *Workspace) AllWebhooks(client *Client, options ...*Options) ([]*Webhook, error) {
	var allWebhooks []*Webhook
	nextPage := &NextPage{}

	var webhooks []*Webhook
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  100,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		webhooks, nextPage, err = w.Webhooks(client, allOptions...)
		if err != nil {
			return nil, err
		}

		allWebhooks = append(allWebhooks, webhooks...)
	}
	return allWebhooks, nil
}

// CreateWebhook registers a new webhook
func (c *Client) CreateWebhook(resource, target string, filters []Filter) (*Webhook, error) {
	m := map[string]interface{}{}
	m["resource"] = resource
	m["target"] = target
	m["filters"] = filters

	result := &Webhook{}

	err := c.post("/webhooks", m, result)

	return result, err
}

// DeleteWebhook deletes an existing webhook
func (c *Client) DeleteWebhook(ID string) error {
	err := c.delete(fmt.Sprintf("/webhooks/%s", ID))
	return err
}

func ParseHook(payload []byte) ([]Event, error) {
	ed := EventData{}
	if err := json.Unmarshal(payload, &ed); err != nil {
		return nil, errors.New("Cannot decode payload")
	}

	return ed.Events, nil
}
