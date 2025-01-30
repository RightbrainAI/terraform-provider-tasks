// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

type Log interface {
	With(args ...any) Log
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type NullLog struct {
}

func (nl NullLog) With(args ...any) Log {
	return nil
}
func (nl NullLog) Debug(args ...any) Log {
	return nil
}
func (nl NullLog) Info(args ...any) Log {
	return nil
}
func (nl NullLog) Warn(args ...any) Log {
	return nil
}
func (nl NullLog) Error(args ...any) Log {
	return nil
}
