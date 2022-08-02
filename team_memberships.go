package asana

import (
	"fmt"
)

// Team is used to group related projects and people together within an
// organization. Each project in an organization is associated with a team.
type TeamMembership struct {

	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	IsGuest bool `json:"is_guest"`

	Team struct {
		ID string `json:"gid"`
		Name string `json:"name"`
	} `json:"team"`

	User struct {
		ID string `json:"gid"`
		Name string `json:"name"`
	} `json:"user"`
}

// Fetch loads the full details for this Team
func (t *TeamMembership) Fetch(client *Client) error {
	client.trace("Loading team membership details for %q\n", t.ID)

	// Use fields options to request Organization field which is not returned by default
	_, err := client.get(fmt.Sprintf("/team_memberships/%s", t.ID), nil, t, Fields(*t))
	return err
}
