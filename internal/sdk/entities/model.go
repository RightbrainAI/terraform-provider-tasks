// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package entitites

type Model struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Alias          string `json:"alias"`
	Provider       string `json:"provider"`
	Description    string `json:"description"`
	SupportsVision bool   `json:"supports_vision"`
}
