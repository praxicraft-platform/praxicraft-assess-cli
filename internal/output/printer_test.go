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
				"id":                 "5a5a3f99-95ef-4bb1-8424-e362b2280fcc",
				"slug":               "mcp-integration-test-jul-2026",
				"status":             "active",
				"title":              "MCP Integration Test — Jul 2026",
				"case_count":         3,
				"time_limit_minutes": 90,
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
	if strings.Contains(out, "5a5a3f99-95ef-4bb1-8424-e362b2280fcc") {
		t.Fatalf("expected shortened id, got:\n%s", out)
	}
	if !strings.Contains(out, "5a5a3f99") {
		t.Fatalf("expected short id prefix, got:\n%s", out)
	}
}

func TestTableInviteCreateShowsFullToken(t *testing.T) {
	var buf bytes.Buffer
	p := &output.Printer{Format: output.Table, Out: &buf}
	token := "inv_live_abcdefghijklmnopqrstuvwxyz0123456789"
	err := p.Print(map[string]any{
		"email":        "candidate@example.com",
		"invite_token": token,
		"status":       "pending",
		"take_url":     "https://assess.praxicraft.com/take/" + token,
		"name":         "Jane Doe",
	})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"email", "candidate@example.com", "invite_token", token, "take_url", "pending", "Jane Doe"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestTableIntegrationsEnvelope(t *testing.T) {
	var buf bytes.Buffer
	p := &output.Printer{Format: output.Table, Out: &buf}
	err := p.Print(map[string]any{
		"ats_allowed": true,
		"integrations": []any{
			map[string]any{
				"provider":    "greenhouse",
				"is_active":   false,
				"has_api_key": false,
				"auth_mode":   "api_key",
				"connect_url": "https://example.com/connect/greenhouse",
			},
			map[string]any{
				"provider":    "lever",
				"is_active":   true,
				"has_api_key": true,
				"auth_mode":   "api_key",
				"connect_url": "https://example.com/connect/lever",
			},
		},
		"meta": map[string]any{
			"plan":              "growth",
			"invite_limit":      500,
			"invites_used":      1,
			"invites_remaining": 499,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, `"integrations"`) {
		t.Fatalf("should not dump raw JSON, got:\n%s", out)
	}
	for _, want := range []string{"provider", "greenhouse", "lever", "is_active", "plan", "growth"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}