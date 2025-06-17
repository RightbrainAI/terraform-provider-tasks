// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk_test

import (
	"bytes"
	"encoding/json"
	"terraform-provider-tasks/internal/sdk"
	"terraform-provider-tasks/internal/sdk/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequests(t *testing.T) {

	t.Run("test that it marshalls json correctly", func(t *testing.T) {
		req := sdk.NewCreateTaskRequest()
		req.OutputFormat["foo"] = entities.OutputFormatExtended{
			Type: "str",
		}
		var data = new(bytes.Buffer)
		err := json.NewEncoder(data).Encode(req)
		assert.NoError(t, err)

		expected := `{
			"output_format": {
				"foo": {
					"type": "str"
				}
			}
		}`

		assert.JSONEq(t, expected, data.String())
	})
}
