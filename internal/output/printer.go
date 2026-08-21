package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jmespath/go-jmespath"
	"gopkg.in/yaml.v3"
)

// Format is the user-facing output format.
type Format string

const (
	JSON  Format = "json"
	Table Format = "table"
	YAML  Format = "yaml"
)

// Printer writes API results.
type Printer struct {
	Format Format
	Query  string
	Out    io.Writer
}

// NewPrinter creates a Printer writing to stdout.
func NewPrinter(format, query string) *Printer {
	f := Format(strings.ToLower(strings.TrimSpace(format)))
	if f == "" {
		f = JSON
	}
	return &Printer{Format: f, Query: strings.TrimSpace(query), Out: os.Stdout}
}

// Print applies --query and writes formatted output.
func (p *Printer) Print(v any) error {
	var data any = v
	if p.Query != "" {
		filtered, err := jmespath.Search(p.Query, v)
		if err != nil {
			return fmt.Errorf("query: %w", err)
		}
		data = filtered
	}

	switch p.Format {
	case YAML:
		enc := yaml.NewEncoder(p.Out)
		enc.SetIndent(2)
		defer enc.Close()
		return enc.Encode(data)
	case Table:
		return printTable(p.Out, data)
	default:
		enc := json.NewEncoder(p.Out)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	}
}

func printTable(w io.Writer, data any) error {
	// Prefer list under "results"; else encode as JSON pretty for complex shapes.
	switch t := data.(type) {
	case []any:
		return printRows(w, t)
	case map[string]any:
		if results, ok := t["results"].([]any); ok {
			return printRows(w, results)
		}
		b, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(b))
		return err
	default:
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(b))
		return err
	}
}

func printRows(w io.Writer, rows []any) error {
	if len(rows) == 0 {
		fmt.Fprintln(w, "(empty)")
		return nil
	}
	// Collect keys from first object
	first, ok := rows[0].(map[string]any)
	if !ok {
		return json.NewEncoder(w).Encode(rows)
	}
	preferred := []string{"id", "slug", "email", "name", "status", "title", "invite_token", "provider"}
	var cols []string
	seen := map[string]bool{}
	for _, k := range preferred {
		if _, ok := first[k]; ok {
			cols = append(cols, k)
			seen[k] = true
		}
	}
	for k := range first {
		if !seen[k] && len(cols) < 6 {
			cols = append(cols, k)
			seen[k] = true
		}
	}
	fmt.Fprintln(w, strings.Join(cols, "\t"))
	for _, row := range rows {
		m, ok := row.(map[string]any)
		if !ok {
			continue
		}
		vals := make([]string, len(cols))
		for i, c := range cols {
			vals[i] = fmt.Sprint(m[c])
		}
		fmt.Fprintln(w, strings.Join(vals, "\t"))
	}
	return nil
}
