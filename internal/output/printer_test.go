package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/output"
)

func TestJSONAndQuery(t *testing.T) {
	var buf bytes.Buffer
	p := &output.Printer{Format: output.JSON, Query: "slug", Out: &buf}
	if err := p.Print(map[string]any{"slug": "senior", "status": "active"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "senior") {
		t.Fatalf("%s", buf.String())
	}
}

func TestTableResults(t *testing.T) {
	var buf bytes.Buffer
	p := &output.Printer{Format: output.Table, Out: &buf}
	err := p.Print(map[string]any{
		"results": []any{
			map[string]any{"slug": "a", "status": "active"},
			map[string]any{"slug": "b", "status": "draft"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "slug") || !strings.Contains(buf.String(), "a") {
		t.Fatalf("%s", buf.String())
	}
}
