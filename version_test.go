// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package buildversion_test

import (
	"encoding/json"
	"os/exec"
	"testing"
)

func TestBuildVersion(t *testing.T) { //nolint:funlen
	tests := []struct {
		title      string
		command    string
		properties map[string]string
	}{
		{
			title:   "go run",
			command: `go run ./testdata/buildversion-cmd`,
			properties: map[string]string{
				"version":        "development",
				"revision":       "",
				"revision-short": "",
				"timestamp":      "",
				"os":             "*",
				"arch":           "*",
				"runtime":        "*",
			},
		},
		{
			title:   "go install",
			command: `go install ./testdata/buildversion-cmd && $(go env GOPATH)/bin/buildversion-cmd`,
			properties: map[string]string{
				"version":        "*",
				"revision":       "*",
				"revision-short": "*",
				"timestamp":      "*",
				"os":             "*",
				"arch":           "*",
				"runtime":        "*",
			},
		},
		{
			title:   "go build",
			command: `go build ./testdata/buildversion-cmd && ./buildversion-cmd`,
			properties: map[string]string{
				"version":        "*",
				"revision":       "*",
				"revision-short": "*",
				"timestamp":      "*",
				"os":             "*",
				"arch":           "*",
				"runtime":        "*",
			},
		},
		{
			title:   "go build with ldflags",
			command: `go build -ldflags "-X main.version=latest -X main.timestamp=now" ./testdata/buildversion-cmd && ./buildversion-cmd`,
			properties: map[string]string{
				"version":        "latest",
				"revision":       "*",
				"revision-short": "*",
				"timestamp":      "now",
				"os":             "*",
				"arch":           "*",
				"runtime":        "*",
			},
		},
		{
			title:   "go build with ldflags and no buildvcs",
			command: `go build -buildvcs=false -ldflags "-X main.version=latest -X main.revision=0000111122223333444455556666777788889999 -X main.timestamp=now" ./testdata/buildversion-cmd && ./buildversion-cmd`,
			properties: map[string]string{
				"version":        "latest",
				"revision":       "0000111122223333444455556666777788889999",
				"revision-short": "0000111",
				"timestamp":      "now",
				"os":             "*",
				"arch":           "*",
				"runtime":        "*",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			output, err := exec.Command("sh", "-c", test.command).CombinedOutput()
			t.Log(string(output))

			if err != nil {
				t.Fatal(err)
			}

			var actual map[string]string
			if err := json.Unmarshal(output, &actual); err != nil {
				t.Fatal(err)
			}

			for key, expectedValue := range test.properties {
				actualValue, found := actual[key]
				if !found {
					t.Fatalf("missing property %q", key)
				}

				switch expectedValue {
				case "*":
					if actualValue == "" {
						t.Errorf("expected key %q to have some value", key)
					}
				default:
					if actualValue != expectedValue {
						t.Errorf("expected key %q to have value %q, got %q", key, expectedValue, actualValue)
					}
				}
			}
		})
	}
}
