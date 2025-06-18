// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package entities_test

import (
	"encoding/json"
	"fmt"
	"os"
	"terraform-provider-tasks/internal/sdk/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRevisionEntity(t *testing.T) {

	t.Run("test that it unmarshalls json correctly", func(t *testing.T) {
		var task entities.Task
		err := json.Unmarshal(getTestFixture(t, "task.json"), &task)
		assert.NoError(t, err)

		rev := task.Revisions[0]
		assert.Equal(t, true, rev.OutputFormat["foo"].IsExtended())
		assert.Equal(t, "str", rev.OutputFormat["foo"].Extended.Type)

		assert.Equal(t, true, rev.OutputFormat["bar"].IsSimple())
		assert.Equal(t, "str", rev.OutputFormat["bar"].Simple.String())
	})
}

// nolint:unparam
func getTestFixture(t *testing.T, fixture string) []byte {
	data, err := os.ReadFile(fmt.Sprintf("../fixtures/%s", fixture))
	assert.NoError(t, err)
	return data
}
