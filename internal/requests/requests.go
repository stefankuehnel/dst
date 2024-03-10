// Copyright 2024 Stefan Kühnel. All rights reserved.
//
// SPDX-License-Identifier: EUPL-1.2

// Package requests provides HTTP utility functions, complementing
// the more common ones in the net/http package.
package requests

import (
	"io"
	"net/http"
	"net/url"
)

// Get issues a GET to the specified URL and closes the response.Body.
// If an error occurs during request, reading of the response.Body or
// closing of the response.Body then the error is returned.
func Get(url url.URL) ([]byte, error) {
	response, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(response.Body)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
