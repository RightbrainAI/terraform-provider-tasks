// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package entitites

import (
	"fmt"
)

// Root represents the overall response structure.
type Task struct {
	AccessToken     string     `json:"access_token"`
	Description     string     `json:"description"`
	Enabled         bool       `json:"enabled"`
	ExposedToAgents bool       `json:"exposed_to_agents"`
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	ProjectID       string     `json:"project_id"`
	Public          bool       `json:"public"`
	Revisions       []Revision `json:"revisions"`
}

func (t *Task) GetActiveRevision() (*Revision, error) {
	for _, r := range t.Revisions {
		if r.Active {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("could not find active revision for task")
}

func (t *Task) GetLatestRevision() (*Revision, error) {
	if len(t.Revisions) == 0 {
		return nil, fmt.Errorf("could not find latest revision for task")
	}
	return &t.Revisions[0], nil
}
