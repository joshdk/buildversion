// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package main

import (
	_ "embed"
	"fmt"

	"github.com/joshdk/buildversion"
)

var (
	// version is an arbitrary version string which can be replaced at build
	// time by using ldflags:
	// -ldflags "-X 'main.version=...'"
	version string

	// version is the current git sha which can be replaced at build time by
	// using ldflags:
	// -ldflags "-X main.revision=...'"
	revision string

	// timestamp is the application build time which can be replaced at build
	// time by using ldflags:
	// -ldflags "-X main.timestamp=...'"
	timestamp string

	//go:embed version.tmpl
	templateBody string

	_ = buildversion.Override(version, revision, timestamp)
)

func main() {
	fmt.Print(buildversion.Template(templateBody))
}
