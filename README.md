# pacmod

![logo](logo/pacmod_logo.png)

[![GoDoc](https://godoc.org/github.com/plexsystems/pacmod?status.svg)](https://godoc.org/github.com/plexsystems/pacmod)
[![Go Report Card](https://goreportcard.com/badge/github.com/plexsystems/pacmod)](https://goreportcard.com/report/github.com/plexsystems/pacmod)

Pacmod is a small tool that can be used to package up your Go modules for distribution. This will be typically used for pushing artifacts to a Go module store such as [Athens](https://github.com/gomods/athens).

## Installation

`go get github.com/plexsystems/pacmod`

## Usage

Run the `pack` command in the directory containing your `go.mod`. For example:

`pacmod pack github.com/repo/user v1.0.0 outputdirectory`

This will result in the following files:

- `go.mod` - The current mod file when the `pack` command was executed
- `v1.0.0.info` - The info file containing the module version and timestamp
- `v1.0.0.zip` - An archive containing the Go module
