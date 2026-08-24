package cmdutil

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/spf13/cobra"
)

// MaxListPages caps --all / picker walks to avoid runaway loops.
const MaxListPages = 250

const loadMoreValue = "__praxicraft_load_more__"

// AllFlag adds --all to follow cursor pagination until exhausted.
func AllFlag(cmd *cobra.Command, dest *bool) {
	cmd.Flags().BoolVar(dest, "all", false, "fetch every page (follow cursor until next is null)")
}

// NextCursor extracts the cursor query param from a page's `next` URL.
func NextCursor(page any) string {
	m, ok := page.(map[string]any)
	if !ok {
		return ""
	}
	next, _ := m["next"].(string)
	next = strings.TrimSpace(next)
	if next == "" {
		return ""
	}
	u, err := url.Parse(next)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(u.Query().Get("cursor"))
}

// NextPage extracts page= from next when cursor pagination is not used.
func NextPage(page any) string {
	m, ok := page.(map[string]any)
	if !ok {
		return ""
	}
	next, _ := m["next"].(string)
	next = strings.TrimSpace(next)
	if next == "" {
		return ""
	}
	u, err := url.Parse(next)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(u.Query().Get("page"))
}

// CloneQuery shallow-copies url.Values.
func CloneQuery(q url.Values) url.Values {
	out := url.Values{}
	if q == nil {
		return out
	}
	for k, vs := range q {
		cp := make([]string, len(vs))
		copy(cp, vs)
		out[k] = cp
	}
	return out
}

// AdvanceQuery sets cursor or page from the current page's next link.
func AdvanceQuery(q url.Values, page any) bool {
	if c := NextCursor(page); c != "" {
		q.Set("cursor", c)
		q.Del("page")
		return true
	}
	if p := NextPage(page); p != "" {
		q.Set("page", p)
		q.Del("cursor")
		return true
	}
	return false
}

// ListOrAll fetches one page, or every page when all is true.
func ListOrAll(all bool, pairs []string, fetch func(url.Values) (any, error)) (any, error) {
	q := QueryFromPairs(pairs)
	if all {
		return FetchAll(fetch, q)
	}
	return fetch(q)
}

// FetchAll walks cursor/page pagination and merges results into one payload.
func FetchAll(fetch func(url.Values) (any, error), q url.Values) (any, error) {
	q = CloneQuery(q)
	if q.Get("page_size") == "" {
		q.Set("page_size", "100")
	}

	var merged []any
	var extra map[string]any
	seen := map[string]struct{}{}

	for i := 0; i < MaxListPages; i++ {
		raw, err := fetch(q)
		if err != nil {
			return nil, err
		}
		if m, ok := raw.(map[string]any); ok && extra == nil {
			extra = map[string]any{}
			for k, v := range m {
				if k == "results" || k == "next" || k == "previous" {
					continue
				}
				extra[k] = v
			}
		}
		rows := ResultsRows(raw)
		if len(rows) == 0 {
			if arr, ok := raw.([]any); ok {
				merged = append(merged, arr...)
				break
			}
		}
		for _, row := range rows {
			key := rowIdentity(row)
			if key != "" {
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
			}
			merged = append(merged, row)
		}
		if !AdvanceQuery(q, raw) {
			break
		}
		if i == MaxListPages-1 {
			return nil, &api.UsageError{Msg: fmt.Sprintf("stopped after %d pages; narrow with --filter", MaxListPages)}
		}
	}

	out := map[string]any{
		"next":     nil,
		"previous": nil,
		"results":  merged,
	}
	for k, v := range extra {
		out[k] = v
	}
	return out, nil
}

func rowIdentity(m map[string]any) string {
	for _, k := range []string{"id", "invite_token", "slug", "token"} {
		if s := fieldString(m, k); s != "" {
			return k + ":" + s
		}
	}
	return ""
}
