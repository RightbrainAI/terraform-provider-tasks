// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

type Log interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type NullLog struct {
}

func (nl NullLog) Debug(msg string, args ...any) {}
func (nl NullLog) Info(msg string, args ...any)  {}
func (nl NullLog) Warn(msg string, args ...any)  {}
func (nl NullLog) Error(msg string, args ...any) {}
