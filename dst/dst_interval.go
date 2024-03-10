// Copyright 2024 Stefan Kühnel. All rights reserved.
//
// SPDX-License-Identifier: EUPL-1.2

package dst

import (
	"math"
	"time"
)

type downloadInterval struct {
	sCent  int
	sTens  int
	sYear  int
	sMonth time.Month

	eCent  int
	eTens  int
	eYear  int
	eMonth time.Month
}

// downloadIntervals returns a list of download intervals. The list of
// download intervals ranges from the start year up to the end year.
// The period indicates the maximum number of years contained in a
// single download interval.
func (dst *dst) downloadIntervals(startYear int, endYear int, period int) []*downloadInterval {
	var downloadIntervals []*downloadInterval

	numberOfIntervals := math.Ceil((float64(endYear) - float64(startYear)) / float64(period))

	// numberOfIntervals is 0 if startYear and endYear are identical.
	// This happens if the caller wants to download dst data for a
	// single year.
	if numberOfIntervals == 0 {
		numberOfIntervals = 1
	}

	currentPeriodStartYear := startYear
	for range int(numberOfIntervals) {
		sMonth := time.January
		sYear := currentPeriodStartYear % 10                      // 2013 => sYear = 3
		sTens := (currentPeriodStartYear - sYear) % 100           // 2013 => sTens = 10
		sCent := (currentPeriodStartYear - (sYear + sTens)) / 100 // 2013 => sCent = 20

		currentPeriodEndYear := currentPeriodStartYear + period

		if currentPeriodEndYear > endYear {
			currentPeriodEndYear = endYear
		}

		eMonth := time.December
		eYear := currentPeriodEndYear % 10                      // 2013 => sYear = 2
		eTens := (currentPeriodEndYear - eYear) % 100           // 2013 => sTens = 10
		eCent := (currentPeriodEndYear - (eYear + eTens)) / 100 // 2013 => sCent = 20

		if currentPeriodEndYear-currentPeriodStartYear == period {
			eYear = eYear - 1
		}

		downloadIntervals = append(downloadIntervals, &downloadInterval{
			sCent,
			sTens,
			sYear,
			sMonth,
			eCent,
			eTens,
			eYear,
			eMonth,
		})

		currentPeriodStartYear = currentPeriodEndYear
	}

	return downloadIntervals
}
