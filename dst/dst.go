// Copyright 2024 Stefan Kühnel. All rights reserved.
//
// SPDX-License-Identifier: EUPL-1.2

// Package dst provides the functions for retrieving the Disturbance
// Storm Time (DST) index.
//
// The DST index is a monitor for the axis-symmetric magnetic
// signature of magnetosphere currents, including mainly the ring
// current, the tail currents and also the magnetopause Chapman-
// Ferraro current.
package dst

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"stefanco.de/dst/internal/requests"
)

type Dst interface {
	// Fetch returns the DST index data from startYear until endYear.
	// It returns an error if the start year precedes the year 1957 in
	// which the collection of the DST index data began. It also returns
	// an error if the end year is beyond the current year. Furthermore,
	// it returns an error, if a http error occurs.
	Fetch(startYear int, endYear int) ([]byte, error)

	// FetchAll returns the DST index data from 1957, the year when the
	// collection of the DST index began, up to the present day. It
	// returns an error, if a http error occurs.
	FetchAll() ([]byte, error)
}

type dst struct {
	// period is the 25-year maximum download period
	period int
}

// New returns a new DST index struct
func New() Dst {
	return &dst{
		period: 25,
	}
}

// downloadUrl returns the URL for downloading DST index data.
//
// The parameters correspond to the start ("s"-prefix) and end year ("e"-prefix)
// of the download interval.
func (dst *dst) downloadUrl(sCent int, sTens int, sYear int, sMonth time.Month, eCent int, eTens int, eYear int, eMonth time.Month) url.URL {
	// using url.Values is not possible, as url.Values
	// sorts the values alphabetically and the Kyoto WDC
	// requires the query parameters not to be sorted
	// alphabetically.
	queryParams := []string{
		fmt.Sprintf("SCent=%d", sCent),
		fmt.Sprintf("STens=%d", sTens),
		fmt.Sprintf("SYear=%d", sYear),
		fmt.Sprintf("SMonth=%d", sMonth),
		fmt.Sprintf("ECent=%d", eCent),
		fmt.Sprintf("ETens=%d", eTens),
		fmt.Sprintf("EYear=%d", eYear),
		fmt.Sprintf("EMonth=%d", eMonth),
		fmt.Sprintf("Image+Type=%s", "GIF"),
		fmt.Sprintf("COLOR=%s", "COLOR"),
		fmt.Sprintf("AE+Sensitivity=%d", 0),
		fmt.Sprintf("Dst+Sensitivity=%d", 0),
		fmt.Sprintf("Output=%s", "DST"),
		fmt.Sprintf("Out+format=%s", "WDC"),
	}

	downloadUrl := url.URL{
		Scheme:   "https",
		Host:     "wdc.kugi.kyoto-u.ac.jp",
		Path:     "/cgi-bin/dstae-cgi",
		RawQuery: strings.Join(queryParams, "&"),
	}

	return downloadUrl
}

// Fetch returns the DST index data from startYear until endYear.
// It returns an error if the start year precedes the year 1957 in
// which the collection of the DST index data began. It also returns
// an error if the end year is beyond the current year. Furthermore,
// it returns an error, if a http error occurs.
func (dst *dst) Fetch(startYear int, endYear int) ([]byte, error) {
	var dstData []byte

	if startYear < 1957 {
		return dstData, errors.New(fmt.Sprintf("expected 'startYear' to be greater than 1957, got %d", startYear))
	}

	if endYear > time.Now().Year() {
		return dstData, errors.New(fmt.Sprintf("expected 'endYear' to be smaller than %d, got %d", time.Now().Year(), endYear))
	}

	downloadIntervals := dst.downloadIntervals(startYear, endYear, dst.period)

	for _, downloadInterval := range downloadIntervals {
		downloadUrl := dst.downloadUrl(
			downloadInterval.sCent,
			downloadInterval.sTens,
			downloadInterval.sYear,
			downloadInterval.sMonth,
			downloadInterval.eCent,
			downloadInterval.eTens,
			downloadInterval.eYear,
			downloadInterval.eMonth)

		body, err := requests.Get(downloadUrl)
		if err != nil {
			return dstData, err
		}

		dstData = append(dstData, body...)
	}

	return dstData, nil
}

// FetchAll returns the DST index data from 1957, the year when the
// collection of the DST index began, up to the current year. It
// returns an error, if a http error occurs.
func (dst *dst) FetchAll() ([]byte, error) {
	startYear := 1957
	endYear := time.Now().Year()

	return dst.Fetch(startYear, endYear)
}
