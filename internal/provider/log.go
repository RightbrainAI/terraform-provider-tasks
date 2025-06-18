// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	sdk "terraform-provider-tasks/internal/sdk"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func NewTerraformLog() sdk.Log {
	return &TerraformLog{}
}

type TerraformLog struct {
}

func (tl TerraformLog) Debug(msg string, args ...any) {
	tflog.Debug(context.Background(), msg, tl.argsToMap(args))
}
func (tl TerraformLog) Info(msg string, args ...any) {
	tflog.Info(context.Background(), msg, tl.argsToMap(args))
}
func (tl TerraformLog) Warn(msg string, args ...any) {
	tflog.Warn(context.Background(), msg, tl.argsToMap(args))
}
func (tl TerraformLog) Error(msg string, args ...any) {
	tflog.Error(context.Background(), msg, tl.argsToMap(args))
}
func (tl TerraformLog) argsToMap(args ...any) map[string]interface{} {
	result := make(map[string]any)
	for i := 0; i < len(args)-1; i += 2 {
		if key, ok := args[i].(string); ok {
			result[key] = args[i+1]
		}
	}
	return result
}
