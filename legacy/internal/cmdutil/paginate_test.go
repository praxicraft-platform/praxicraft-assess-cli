package cmdutil

import (
	"net/url"
	"testing"
)

func TestNextCursor(t *testing.T) {
	page := map[string]any{
		"next":     "https://assess.praxicraft.com/api/v1/public/assessments/?cursor=abc123&page_size=20",
		"previous": nil,
		"results":  []any{},
	}
	if got := NextCursor(page); got != "abc123" {
		t.Fatalf("NextCursor=%q want abc123", got)
	}
	if NextCursor(map[string]any{"next": nil}) != "" {
		t.Fatal("expected empty cursor")
	}
}

func TestNextPage(t *testing.T) {
	page := map[string]any{
		"next": "https://example.com/results/?page=2",
	}
	if got := NextPage(page); got != "2" {
		t.Fatalf("NextPage=%q want 2", got)
	}
}

func TestFetchAllMergesPages(t *testing.T) {
	calls := 0
	fetch := func(q url.Values) (any, error) {
		calls++
		switch calls {
		case 1:
			if q.Get("page_size") != "100" {
				t.Fatalf("expected default page_size=100, got %q", q.Get("page_size"))
			}
			return map[string]any{
				"next": "https://x/?cursor=c2",
				"results": []any{
					map[string]any{"slug": "a"},
					map[string]any{"slug": "b"},
				},
			}, nil
		case 2:
			if q.Get("cursor") != "c2" {
				t.Fatalf("expected cursor=c2, got %q", q.Get("cursor"))
			}
			return map[string]any{
				"next": nil,
				"results": []any{
					map[string]any{"slug": "c"},
				},
			}, nil
		default:
			t.Fatal("unexpected extra fetch")
			return nil, nil
		}
	}
	out, err := FetchAll(fetch, nil)
	if err != nil {
		t.Fatal(err)
	}
	rows := ResultsRows(out)
	if len(rows) != 3 {
		t.Fatalf("got %d rows want 3", len(rows))
	}
	m := out.(map[string]any)
	if m["next"] != nil {
		t.Fatalf("merged next should be nil, got %v", m["next"])
	}
}

func TestAdvanceQuery(t *testing.T) {
	q := url.Values{}
	ok := AdvanceQuery(q, map[string]any{
		"next": "https://x/?cursor=zz",
	})
	if !ok || q.Get("cursor") != "zz" {
		t.Fatalf("AdvanceQuery cursor failed: ok=%v q=%v", ok, q)
	}
}
