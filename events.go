package asana

import (
	"encoding/json"
	"time"
)

type EventData struct {
	Events []Event `json:"events"`
}

type Event struct {
	User struct {
		ID           string `json:"gid"`
		ResourceType string `json:"resource_type"`
	} `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	Action    string    `json:"action"`
	Parent    struct {
		ID              string `json:"gid"`
		ResourceType    string `json:"resource_type"`
		ResourceSubtype string `json:"resource_subtype"`
	} `json:"parent"`
	Change struct {
		Field    string `json:"field"`
		Action   string `json:"action"`
		NewValue struct {
			ID           string `json:"gid"`
			ResourceType string `json:"resource_type"`
		} `json:"new_value"`
	} `json:"change"`
	Resource struct {
		ID              string `json:"gid"`
		ResourceType    string `json:"resource_type"`
		ResourceSubtype string `json:"resource_subtype"`
	} `json:"resource"`
}

func (e Event) String() string {
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}
