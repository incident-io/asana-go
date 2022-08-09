package asana

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"time"
)

// Signature headers
const (
	hSecret    = "X-Hook-Secret"
	hSignature = "X-Hook-Signature"
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
func (w *Workspace) Webhooks(ctx context.Context, client *Client, options ...*Options) ([]*Webhook, *NextPage, error) {
	client.trace("Listing webhooks for workspace %s...\n", w.ID)
	var result []*Webhook

	workspace := &Options{
		Workspace: w.ID,
	}

	allOptions := append([]*Options{workspace}, options...)

	// Make the request
	nextPage, err := client.get(ctx, "/webhooks", nil, &result, allOptions...)
	return result, nextPage, err
}

// AllWebhooks repeatedly pages through all available webhooks for a user
func (w *Workspace) AllWebhooks(ctx context.Context, client *Client, options ...*Options) ([]*Webhook, error) {
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
		webhooks, nextPage, err = w.Webhooks(ctx, client, allOptions...)
		if err != nil {
			return nil, err
		}

		allWebhooks = append(allWebhooks, webhooks...)
	}
	return allWebhooks, nil
}

// CreateWebhook registers a new webhook
func (c *Client) CreateWebhook(ctx context.Context, resource, target string, filters []Filter) (*Webhook, error) {
	m := map[string]interface{}{}
	m["resource"] = resource
	m["target"] = target
	m["filters"] = filters

	result := &Webhook{}

	err := c.post(ctx, "/webhooks", m, result)

	return result, err
}

// DeleteWebhook deletes an existing webhook
func (c *Client) DeleteWebhook(ctx context.Context, ID string) error {
	err := c.delete(ctx, fmt.Sprintf("/webhooks/%s", ID))
	return err
}

func ParseHook(body io.ReadCloser) ([]Event, error) {
	ed := EventData{}
	if err := json.NewDecoder(body).Decode(&ed); err != nil {
		return nil, errors.New("Cannot decode payload")
	}

	return ed.Events, nil
}

// SecretsVerifier contains the information needed to verify that the request comes from Asana
type SecretsVerifier struct {
	signature []byte
	hmac      hash.Hash
}

// NewSecretsVerifier returns a SecretsVerifier object in exchange for an http.Header object and signing secret
func NewSecretsVerifier(header http.Header, secret string) (sv SecretsVerifier, err error) {
	var bsignature []byte

	signature := header.Get(hSignature)
	if signature == "" {
		return SecretsVerifier{}, errors.New("Missing header")
	}

	if bsignature, err = hex.DecodeString(signature); err != nil {
		return SecretsVerifier{}, err
	}

	hash := hmac.New(sha256.New, []byte(secret))

	return SecretsVerifier{
		signature: bsignature,
		hmac:      hash,
	}, nil
}

func (v *SecretsVerifier) Write(body []byte) (n int, err error) {
	return v.hmac.Write(body)
}

// Ensure compares the signature sent from Slack with the actual computed hash to judge validity
func (v SecretsVerifier) Ensure() error {
	computed := v.hmac.Sum(nil)
	if hmac.Equal(computed, v.signature) {
		return nil
	}
	return fmt.Errorf("Computed unexpected signature of: %s", hex.EncodeToString(computed))
}
