// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
