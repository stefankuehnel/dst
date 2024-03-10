// Copyright 2024 Stefan Kühnel. All rights reserved.
//
// SPDX-License-Identifier: EUPL-1.2
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"stefanco.de/dst/dst"
)

var currentYear = time.Now().Year()

var usage = fmt.Sprintf(`
Usage:
    dst [--all --output <file>]
    dst (--start-year <start-year> --end-year <end-year>) [-o <file>]

Options:
    -a, --all                   download everything from 1957 to %d 
    -s, --start-year            download interval start year 
                                [default: 1957]
    -e, --end-year              download interval end year 
                                [default: %d]
    -o <file>, --output <file>  write the result to the file at 
                                path <file>
    -v, --version               output version information and exit
    -h, --help                  display this help and exit

dst is a utility tool for downloading Disturbance Storm Time (DST) 
index data (final, provisional, real-time) from 1957 up to %d.

Example:
    $ dst --all
    DST5701*01  X219 000 011 013 012 012 009 007 007 006 002-001-007-007-008-001 009 008 004 000 001 003 002 004 009 009 004
    DST5701*02  X219 000 011 003 006 009 010 012 007 005 008 032 007-007-007-002-001 001 002 005 005-014-041-065-065-059-006
    ...
    DST2403*31RRX020   09999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999
    $ dst --all --output=dst.txt
    $ dst --start-year=1957 --end-year=2024
    DST5701*01  X219 000 011 013 012 012 009 007 007 006 002-001-007-007-008-001 009 008 004 000 001 003 002 004 009 009 004
    DST5701*02  X219 000 011 003 006 009 010 012 007 005 008 032 007-007-007-002-001 001 002 005 005-014-041-065-065-059-006
    ...
    DST2403*31RRX020   09999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999
    $ dst --start-year=1957 --end-year=2024 --output=dst.txt
`, currentYear, currentYear, currentYear)

// Version can be set at link time to override debug.BuildInfo.Main.Version,
// which is "(devel)" when building from within the module.
//
// See: https://golang.org/issue/29814
// See: https://golang.org/issue/29228
// See: Dockerfile
var Version string

func main() {
	flag.Usage = getUsage

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	var (
		err           error
		dstIndexData  []byte
		allFlag       bool
		startYearFlag int
		endYearFlag   int
		outputFlag    string
		versionFlag   bool
		helpFlag      bool
	)

	flag.BoolVar(&allFlag, "a", false, fmt.Sprintf("download everthing from 1957 to %d", currentYear))
	flag.BoolVar(&allFlag, "all", false, fmt.Sprintf("download everthing from 1957 to %d", currentYear))
	flag.IntVar(&startYearFlag, "s", 1957, "download interval start year ")
	flag.IntVar(&startYearFlag, "start-year", 1957, "download interval start year ")
	flag.IntVar(&endYearFlag, "e", currentYear, "download interval end year ")
	flag.IntVar(&endYearFlag, "end-year", currentYear, "download interval end year ")
	flag.StringVar(&outputFlag, "o", "", "write the result to the file at path <file>")
	flag.StringVar(&outputFlag, "output", "", "write the result to the file at path <file>")
	flag.BoolVar(&versionFlag, "v", false, "output version information and exit")
	flag.BoolVar(&versionFlag, "version", false, "output version information and exit")
	flag.BoolVar(&helpFlag, "h", false, "display this help and exit")
	flag.BoolVar(&helpFlag, "help", false, "display this help and exit")

	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		fmt.Println(getVersion())
		os.Exit(0)
	}

	isAllFlagPassed := isFlagPassed("a") || isFlagPassed("all")
	isStartYearFlagPassed := isFlagPassed("s") || isFlagPassed("start-year")
	isEndYearFlagPassed := isFlagPassed("e") || isFlagPassed("end-year")
	isOutputFlagPassed := isFlagPassed("o") || isFlagPassed("output")

	if isAllFlagPassed {
		dstIndexData, err = dst.New().FetchAll()
		if err != nil {
			fail(err)
		}
	}

	if isStartYearFlagPassed && !isEndYearFlagPassed {
		fail(fmt.Errorf("missing -e/--end-year"))
	}

	if !isStartYearFlagPassed && isEndYearFlagPassed {
		fail(fmt.Errorf("missing -s/--start-year"))
	}

	if isStartYearFlagPassed && isEndYearFlagPassed {
		if startYearFlag < 1957 {
			fail(fmt.Errorf("expected -s/--start-year to be greater than 1957, got %d", startYearFlag))
		}

		if endYearFlag > currentYear {
			fail(fmt.Errorf("expected -e/--end-year to be less than %d, got %d", currentYear, endYearFlag))
		}

		dstIndexData, err = dst.New().Fetch(startYearFlag, endYearFlag)
		if err != nil {
			fail(err)
		}
	}

	if isOutputFlagPassed {
		err := os.WriteFile(outputFlag, dstIndexData, 0666)
		if err != nil {
			fail(err)
		}
		os.Exit(0)
	}

	fmt.Println(string(dstIndexData))
	os.Exit(0)
}

func getUsage() {
	_, err := fmt.Fprintf(os.Stderr, "%s\n\n", strings.TrimSpace(usage))
	if err != nil {
		fail(err)
	}
}

func getVersion() string {
	if Version != "" {
		return fmt.Sprintf("dst (%s) %s/%s", Version, runtime.GOOS, runtime.GOARCH)
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		return fmt.Sprintf("dst %s %s/%s", buildInfo.Main.Version, runtime.GOOS, runtime.GOARCH)
	}

	return fmt.Sprintf("dst (unknown) %s/%s", runtime.GOOS, runtime.GOARCH)
}

func fail(err error) {
	fmt.Printf("error: %s\n", err)
	os.Exit(1)
}

func isFlagPassed(name string) bool {
	// See: https://stackoverflow.com/a/54747682
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
