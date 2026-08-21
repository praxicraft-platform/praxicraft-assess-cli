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
			map[string]any{
				"id":                  "5a5a3f99-95ef-4bb1-8424-e362b2280fcc",
				"slug":                "mcp-integration-test-jul-2026",
				"status":              "active",
				"title":               "MCP Integration Test — Jul 2026",
				"case_count":          3,
				"time_limit_minutes":  90,
			},
			map[string]any{
				"id":                 "5188d828-ba30-427d-a7ee-48f22c1a9bac",
				"slug":               "core-data-engineers-cohort-admission-assessment",
				"status":             "active",
				"title":              "Core Data Engineers:  Cohort Admission Assessment",
				"case_count":         60,
				"time_limit_minutes": 20,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "slug") || !strings.Contains(out, "mcp-integration-test-jul-2026") {
		t.Fatalf("%s", out)
	}
	// Aligned columns: header and first data cell should not be raw-tab only mess;
	// short id should appear instead of full UUID.
	if strings.Contains(out, "5a5a3f99-95ef-4bb1-8424-e362b2280fcc") {
		t.Fatalf("expected shortened id, got:\n%s", out)
	}
	if !strings.Contains(out, "5a5a3f99") {
		t.Fatalf("expected short id prefix, got:\n%s", out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected header + rows:\n%s", out)
	}
}

