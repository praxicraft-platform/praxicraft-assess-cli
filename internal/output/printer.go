package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

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

const maxCellRunes = 56
const maxSlugRunes = 64

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
	first, ok := rows[0].(map[string]any)
	if !ok {
		return json.NewEncoder(w).Encode(rows)
	}

	preferred := []string{
		"slug", "email", "name", "title", "status", "provider",
		"case_count", "time_limit_minutes", "invitation_count", "id",
	}
	var cols []string
	seen := map[string]bool{}
	for _, k := range preferred {
		if v, ok := first[k]; ok && isScalar(v) {
			cols = append(cols, k)
			seen[k] = true
		}
	}
	// Only backfill sparse objects (fewer than 4 preferred hits).
	if len(cols) < 4 {
		var extra []string
		for k, v := range first {
			if seen[k] || !isScalar(v) {
				continue
			}
			extra = append(extra, k)
		}
		sortStrings(extra)
		for _, k := range extra {
			if len(cols) >= 6 {
				break
			}
			cols = append(cols, k)
			seen[k] = true
		}
	}
	if len(cols) == 0 {
		return json.NewEncoder(w).Encode(rows)
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(cols, "\t"))
	for _, row := range rows {
		m, ok := row.(map[string]any)
		if !ok {
			continue
		}
		vals := make([]string, len(cols))
		for i, c := range cols {
			vals[i] = formatCell(c, m[c])
		}
		fmt.Fprintln(tw, strings.Join(vals, "\t"))
	}
	return tw.Flush()
}

func isScalar(v any) bool {
	switch v.(type) {
	case nil, bool, string, json.Number,
		float32, float64,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

func formatCell(col string, v any) string {
	if v == nil {
		return ""
	}
	s := fmt.Sprint(v)
	switch col {
	case "id", "invite_token":
		s = shortenID(s)
		return s
	case "slug", "email":
		return truncateRunes(s, maxSlugRunes)
	default:
		return truncateRunes(s, maxCellRunes)
	}
}

func shortenID(s string) string {
	// UUID-looking values: show first segment for scanability.
	if len(s) == 36 && strings.Count(s, "-") == 4 {
		return s[:8]
	}
	if utf8.RuneCountInString(s) > 12 {
		return truncateRunes(s, 12)
	}
	return s
}

func truncateRunes(s string, max int) string {
	if max <= 1 || utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max-1]) + "…"
}

func sortStrings(a []string) {
	for i := 1; i < len(a); i++ {
		j := i
		for j > 0 && a[j] < a[j-1] {
			a[j], a[j-1] = a[j-1], a[j]
			j--
		}
	}
}
