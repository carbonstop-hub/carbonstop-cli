// Package formatter handles output formatting (JSON, table, raw text).
package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Formatter formats and prints API responses.
type Formatter struct {
	Raw bool
}

// New creates a new Formatter.
func New(raw bool) *Formatter {
	return &Formatter{Raw: raw}
}

// Print formats and outputs the response body.
func (f *Formatter) Print(status int, body string) {
	if f.Raw {
		output := body
		if !strings.HasSuffix(output, "\n") {
			output += "\n"
		}
		os.Stdout.WriteString(output)
		if status < 200 || status >= 300 {
			fmt.Fprintf(os.Stderr, "[carbonstop] http status %d\n", status)
		}
		return
	}

	var parsed interface{}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		output := body
		if !strings.HasSuffix(output, "\n") {
			output += "\n"
		}
		os.Stdout.WriteString(output)
	} else {
		pretty, err := json.MarshalIndent(parsed, "", "  ")
		if err != nil {
			os.Stdout.WriteString(body + "\n")
		} else {
			os.Stdout.WriteString(string(pretty) + "\n")
		}
	}

	if status < 200 || status >= 300 {
		fmt.Fprintf(os.Stderr, "[carbonstop] http status %d\n", status)
	}
}

// Info prints an informational message to stderr.
func Info(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[carbonstop] "+format+"\n", args...)
}
