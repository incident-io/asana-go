package asana

// StoryBase contains the text of a story, as used when creating a new comment
type StoryBase struct {
	// Create-only. Human-readable text for the story or comment. This will
	// not include the name of the creator.
	//
	// Note: This is not guaranteed to be stable for a given type of story.
	// For example, text for a reassignment may not always say “assigned to
	// …”. The API currently does not provide a structured way of inspecting
	// the meaning of a story.
	Text string `json:"text,omitempty"`
}

// Story represents an activity associated with an object in the Asana
// system. Stories are generated by the system whenever users take actions
// such as creating or assigning tasks, or moving tasks between projects.
// Comments are also a form of user-generated story.
//
// Stories are a form of history in the system, and as such they are read-
// only. Once generated, it is not possible to modify a story.
type Story struct {
	StoryBase

	HasID
	HasCreated
	HasHearts

	expandable

	// The user who created the story.
	CreatedBy *User `json:"created_by,omitempty"`

	// Read-only. HTML formatted text for a comment. This will not include the
	// name of the creator.
	//
	// Note: This field is only returned if explicitly requested using the
	// opt_fields query parameter.
	HTMLText string `json:"html_text,omitempty"`

	// Read-only. The object this story is associated with. Currently may only
	// be a task.
	Target *Task `json:"target,omitempty"`

	// Read-only. The component of the Asana product the user used to trigger
	// the story.
	Source string `json:"source,omitempty"`

	// Read-only. The type of story this is.
	Type string `json:"type,omitempty"`
}

// Stories lists all stories attached to a task
func (t *Task) Stories(opts ...*Options) ([]*Story, error) {
	t.trace("Listing stories for %q", t.Name)

	var result []*Story

	// Make the request
	err := t.Client.get(fmt.Sprintf("/tasks/%d/stories", t.ID), nil, &result, opts...)
	return result, err
}

// CreateComment adds a comment story to a task
func (t *Task) CreateComment(story *StoryBase) {
	c.info("Creating comment for task %q", t.Name)

	result := &Story{}
	result.expanded = true

	err := c.post(fmt.Sprintf("/tasks/%d/stories", t.ID), result)
	return result, err
}
