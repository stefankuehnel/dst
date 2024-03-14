# Dst

[![Docker](../../actions/workflows/docker.yml/badge.svg)](../../actions/workflows/docker.yml)

A Go-based, open-source CLI tool without dependencies for downloading [Disturbance Storm Time (DST)](https://wdc.kugi.kyoto-u.ac.jp/dstdir/index.html) index data.

## ⚙️ Get Started

You'll need [Go](https://go.dev) installed.

### Install

First of all, you need to install `dst` locally:

```shell
$ go install stefanco.de/dst/cmd/...@latest
```

This will install `dst` into `$GOROOT/bin`.

### Run locally

Then you're able to run `dst` locally:

```shell
$ dst --help
```

## 👨‍💻 Usage

This message is also available when running `$ dst --help`.

```text
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
```

## 🔨 Technology

The following technologies, tools and platforms were used during development.

- **Code**: [Go](https://go.dev)
- **CI/CD**: [GitHub Actions](https://github.com/actions)

## 👷‍ Error Found?

Thank you for your message! Please fill out a [bug report](../../issues/new?assignees=&labels=&template=bug_report.md&title=).

## License

This project is licensed under the [European Union Public License 1.2](https://choosealicense.com/licenses/eupl-1.2/).