[![License][license-badge]][license-link]
[![Actions][github-actions-badge]][github-actions-link]

# Build Version

🏗️ Library for determining and overriding the version of a Go application

## Motivations

There are currently a number of ways to distribute a Go application. 
A user could download a prebuilt release binary, or run `go install`, or even clone the repo and `go build`/`go run` the code directly.

Information such as version strings, git commit SHAs, and timestamps may be available to programs executed in each of these scenarios.

This library aims to surface this information for inclusion in bug reports, or when displaying help text.

## Usage

### Basic Example

Here is a simple program which templates out the default version, revision, & timestamp properties.

```go
package main

import (
	"fmt"

	"github.com/joshdk/buildversion"
)

func main() {
	tpl := `{{ .Version }}
	{{- if .Revision }} ({{ .Revision }}){{ end }}
	{{- if .Timestamp }} built on {{ .Timestamp }}{{ end }}`

	fmt.Print(buildversion.Template(tpl))
}
```

When executed with `go install` or `go build` we can observe that the current commit, and the time that commit was originally authored, is printed.

```shell
$ go install .

$ demo
v0.0.1-0.20250707135443-82e199f920b9 (82e199f) built on 2025-07-07T13:54:43Z
```

```shell
$ go build .

$ ./demo
v0.0.1-0.20250707135443-82e199f920b9 (82e199f) built on 2025-07-07T13:54:43Z
```

Commit information is not included when executed via `go run`, so a default version is printed instead.

```shell
$ go run .
development
```

### Overriding the Version by Using ldflags

This program can be enhanced to consume build-time provided version, revision, & timestamp values as custom overrides.

```go
var version, revision, timestamp string

func main() {
	tpl := `{{ .Version }}
	{{- if .Revision }} ({{ .Revision }}){{ end }}
	{{- if .Timestamp }} built on {{ .Timestamp }}{{ end }}`

	buildversion.Override(version, revision, timestamp)
	
	fmt.Print(buildversion.Template(tpl))
}
```

These values can be provided via `-ldflags`/`-X` (see [https://pkg.go.dev/cmd/link](https://pkg.go.dev/cmd/link)).

```shell
$ go build -buildvcs=false -ldflags "-X 'main.version=$(git describe --tags --always)' -X 'main.revision=$(git rev-parse HEAD)' -X 'main.timestamp=$(date -u '+%Y-%m-%dT%H:%M:%SZ')'" .             

$ ./demo
v1.2.3 (82e199f) built on 2025-07-19T03:16:33Z
```

### Templating

The following properties can be referenced from a Go `text/template` style template body: 

| Property               | Description                                                 |
|------------------------|-------------------------------------------------------------|
| `{{ .Path }}`          | Go import path of the main package being run                |
| `{{ .Version }}`       | Package version string                                      |
| `{{ .Revision }}`      | Current git commit SHA                                      |
| `{{ .ShortRevision }}` | Current git commit SHA, truncated to the first 7 characters |
| `{{ .Timestamp }}`     | RFC3339 timestamp of the current git commit                 |
| `{{ .OS }}`            | Current operating system type (`GOOS`)                      |
| `{{ .Arch }}`          | Current CPU architecture (`GOARCH`)                         |
| `{{ .Runtime }}`       | Go runtime version                                          |

The `{{ .Version }}`, `{{ .Revision }}`/`{{ .ShortRevision }}`, & `{{ .Timestamp }}` are not guaranteed to contain a value if not explicitly overridden.

## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

[github-actions-badge]:  https://github.com/joshdk/buildversion/workflows/Build/badge.svg
[github-actions-link]:   https://github.com/joshdk/buildversion/actions
[license-badge]:         https://img.shields.io/badge/license-MIT-green.svg
[license-file]:          https://github.com/joshdk/buildversion/blob/master/LICENSE.txt
[license-link]:          https://opensource.org/licenses/MIT
