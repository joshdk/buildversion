// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package buildversion asdf...
package buildversion

import (
	"bytes"
	"runtime"
	"runtime/debug"
	"text/template"
)

// versionData holds several build-time and run-time properties for use in
// templating.
type versionData struct {
	// Path is the import path to the main module.
	Path string

	// Version is an arbitrary version string.
	Version string

	// Revision is a git commit SHA.
	Revision string

	// ShortRevision is a short (truncated to the first 7 characters) variant
	// of the Revision.
	ShortRevision string

	// Timestamp is and RFC3339 timestamp.
	Timestamp string

	// OS is the current operating system.
	OS string

	// Arch is the current CPU architecture.
	Arch string

	// Runtime is the current Go runtime version.
	Runtime string
}

var data = versionData{
	OS:      runtime.GOOS,
	Arch:    runtime.GOARCH,
	Runtime: runtime.Version(),
}

// init is used to read the current build data and extracts path, version,
// revision, & timestamp if possible.
func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	// Extract main package import path.
	data.Path = buildInfo.Main.Path

	// Extract the "vcs.revision" value from build settings.
	for _, setting := range buildInfo.Settings {
		if setting.Key == "vcs.revision" {
			data.Revision = setting.Value
			data.ShortRevision = setting.Value[:7]

			break
		}
	}

	// Extract the "vcs.time" value from build settings.
	for _, setting := range buildInfo.Settings {
		if setting.Key == "vcs.time" {
			data.Timestamp = setting.Value

			break
		}
	}

	// Extract main package version.
	if buildInfo.Main.Version != "(devel)" {
		data.Version = buildInfo.Main.Version
	}

	// Set version to "development" as a fallback.
	if data.Version == "" {
		data.Version = "development"
	}
}

// Override replaces the built-in version and timestamp with custom values.
// Intended to be used with values from e.g. a release build.
func Override(version, revision, timestamp string) any {
	if version != "" {
		data.Version = version
	}

	if revision != "" {
		data.Revision = revision
		data.ShortRevision = revision[:7]
	}

	if timestamp != "" {
		data.Timestamp = timestamp
	}

	return nil
}

// Template renders the provided gotemplate using the built-in/overridden build
// data. Not intended to be used with a dynamic template, so the function
// panics if any error is encountered.
func Template(body string) string {
	var buf bytes.Buffer

	if err := template.Must(template.New("").Parse(body)).Execute(&buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}
