package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
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
	data := normalize(v)
	if p.Query != "" {
		filtered, err := jmespath.Search(p.Query, data)
		if err != nil {
			return fmt.Errorf("query: %w", err)
		}
		data = normalize(filtered)
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
		return writeJSON(p.Out, data)
	}
}

func writeJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}

// normalize re-encodes via JSON so nested maps are map[string]any and
// numbers stay readable (avoids odd decoder leftovers).
func normalize(v any) any {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var out any
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.UseNumber()
	if err := dec.Decode(&out); err != nil {
		return v
	}
	return out
}

func printTable(w io.Writer, data any) error {
	switch t := data.(type) {
	case nil:
		fmt.Fprintln(w, "(null)")
		return nil
	case []any:
		return printRows(w, t)
	case map[string]any:
		if rows, rest, ok := extractListEnvelope(t); ok {
			if err := printRows(w, rows); err != nil {
				return err
			}
			if len(rest) > 0 {
				fmt.Fprintln(w)
				fmt.Fprintln(w, "—")
				return printObject(w, rest)
			}
			return nil
		}
		// Single resource object (invite create, get, whoami, …)
		return printObject(w, t)
	default:
		return writeJSON(w, data)
	}
}

// extractListEnvelope finds the main row array in list-shaped API payloads.
// e.g. { "results": [...] }, { "integrations": [...], "meta": {...}, "ats_allowed": true }
func extractListEnvelope(m map[string]any) (rows []any, rest map[string]any, ok bool) {
	preferred := []string{
		"results", "integrations", "items", "data", "webhooks",
		"deliveries", "enrollments", "members", "cases", "invites",
	}
	var listKey string
	for _, k := range preferred {
		if arr, is := m[k].([]any); is {
			listKey = k
			rows = arr
			break
		}
	}
	if listKey == "" {
		// Fallback: exactly one array-of-objects field.
		var found string
		for k, v := range m {
			arr, is := v.([]any)
			if !is || len(arr) == 0 {
				continue
			}
			if _, isMap := arr[0].(map[string]any); !isMap {
				continue
			}
			if found != "" {
				return nil, nil, false
			}
			found = k
			rows = arr
		}
		if found == "" {
			return nil, nil, false
		}
		listKey = found
	}

	rest = map[string]any{}
	for k, v := range m {
		if k == listKey {
			continue
		}
		rest[k] = v
	}
	// Flatten meta into rest scalars when it's a small object.
	if meta, is := rest["meta"].(map[string]any); is {
		delete(rest, "meta")
		for k, v := range meta {
			if _, exists := rest[k]; !exists && isScalar(v) {
				rest[k] = v
			}
		}
	}
	return rows, rest, true
}

func printObject(w io.Writer, m map[string]any) error {
	if len(m) == 0 {
		fmt.Fprintln(w, "(empty)")
		return nil
	}

	preferred := []string{
		"email", "name", "slug", "title", "status", "invite_token", "token",
		"take_url", "url", "id", "assessment_slug", "assessment",
		"created_at", "updated_at", "expires_at", "sent_at",
		"provider", "role", "plan", "ats_allowed",
		"invite_limit", "invites_used", "invites_remaining",
	}

	keys := make([]string, 0, len(m))
	seen := map[string]bool{}
	for _, k := range preferred {
		if _, ok := m[k]; ok {
			keys = append(keys, k)
			seen[k] = true
		}
	}
	var rest []string
	for k := range m {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	keys = append(keys, rest...)

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	for _, k := range keys {
		v := m[k]
		if isScalar(v) {
			fmt.Fprintf(tw, "%s\t%s\n", k, formatObjectValue(k, v))
			continue
		}
		// Nested object/array: show key, then indented JSON block.
		fmt.Fprintf(tw, "%s\t\n", k)
		_ = tw.Flush()
		b, err := json.MarshalIndent(v, "  ", "  ")
		if err != nil {
			fmt.Fprintf(w, "  %v\n", v)
			continue
		}
		fmt.Fprintf(w, "%s\n", string(b))
		tw = tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	}
	return tw.Flush()
}

func formatObjectValue(key string, v any) string {
	if v == nil {
		return ""
	}
	switch n := v.(type) {
	case json.Number:
		return n.String()
	case bool:
		return fmt.Sprintf("%v", n)
	case string:
		// Don't shorten invite_token / take_url / ids on detail view — users need full values.
		switch key {
		case "invite_token", "token", "take_url", "url", "id":
			return n
		default:
			return n
		}
	default:
		return fmt.Sprint(v)
	}
}

func printRows(w io.Writer, rows []any) error {
	if len(rows) == 0 {
		fmt.Fprintln(w, "(empty)")
		return nil
	}
	first, ok := rows[0].(map[string]any)
	if !ok {
		return writeJSON(w, rows)
	}

	preferred := []string{
		"slug", "email", "name", "title", "status", "provider",
		"is_active", "has_api_key", "auth_mode", "invite_token",
		"case_count", "time_limit_minutes", "invitation_count", "id",
		"connect_url", "webhook_url",
	}
	var cols []string
	seen := map[string]bool{}
	for _, k := range preferred {
		if v, ok := first[k]; ok && isScalar(v) {
			cols = append(cols, k)
			seen[k] = true
		}
	}
	if len(cols) < 4 {
		var extra []string
		for k, v := range first {
			if seen[k] || !isScalar(v) {
				continue
			}
			extra = append(extra, k)
		}
		sort.Strings(extra)
		for _, k := range extra {
			if len(cols) >= 8 {
				break
			}
			cols = append(cols, k)
			seen[k] = true
		}
	}
	if len(cols) == 0 {
		return writeJSON(w, rows)
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
	if n, ok := v.(json.Number); ok {
		return n.String()
	}
	s := fmt.Sprint(v)
	switch col {
	case "id":
		return shortenID(s)
	case "invite_token":
		// List view: shorten; detail view uses printObject (full).
		return shortenID(s)
	case "slug", "email":
		return truncateRunes(s, maxSlugRunes)
	default:
		return truncateRunes(s, maxCellRunes)
	}
}

func shortenID(s string) string {
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
