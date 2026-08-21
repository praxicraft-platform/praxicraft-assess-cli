package cmdutil

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/runtime"
	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/ui"
)

// ResultsRows extracts list rows from a Public API list payload.
func ResultsRows(v any) []map[string]any {
	switch t := v.(type) {
	case map[string]any:
		if results, ok := t["results"].([]any); ok {
			return anyMaps(results)
		}
		return nil
	case []any:
		return anyMaps(t)
	default:
		return nil
	}
}

func anyMaps(rows []any) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		if m, ok := r.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func fieldString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			s := strings.TrimSpace(fmt.Sprint(v))
			if s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// ChoiceFromRow builds a picker option from an API row.
type ChoiceFromRow func(map[string]any) (ui.Choice, bool)

// PickPaged walks list pages with a "Load more…" option until the user picks or cancels.
func PickPaged(rt *runtime.Runtime, title string, emptyMsg string, fetch func(url.Values) (any, error), fromRow ChoiceFromRow) (string, error) {
	if err := rt.EnsureAPI(); err != nil {
		return "", err
	}
	q := url.Values{}
	q.Set("page_size", "50")

	var accumulated []ui.Choice
	seen := map[string]struct{}{}

	for page := 0; page < MaxListPages; page++ {
		raw, err := fetch(q)
		if err != nil {
			return "", err
		}
		for _, m := range ResultsRows(raw) {
			c, ok := fromRow(m)
			if !ok || c.Value == "" {
				continue
			}
			if _, dup := seen[c.Value]; dup {
				continue
			}
			seen[c.Value] = struct{}{}
			accumulated = append(accumulated, c)
		}

		nextQ := CloneQuery(q)
		more := AdvanceQuery(nextQ, raw)

		if len(accumulated) == 0 {
			if more {
				q = nextQ
				continue
			}
			return "", &api.UsageError{Msg: emptyMsg}
		}

		choices := make([]ui.Choice, len(accumulated))
		copy(choices, accumulated)
		if more {
			choices = append(choices, ui.Choice{
				Label: fmt.Sprintf("↓ Load more… (%d so far)", len(accumulated)),
				Value: loadMoreValue,
			})
		}

		selected, err := ui.Select(rt.UI, title, choices)
		if err != nil {
			return "", err
		}
		if selected == loadMoreValue {
			if !more {
				return "", &api.UsageError{Msg: "no more pages"}
			}
			q = nextQ
			continue
		}
		return selected, nil
	}
	return "", &api.UsageError{Msg: fmt.Sprintf("stopped after %d pages", MaxListPages)}
}

// PickAssessmentSlug lists assessments (with Load more) and returns a slug.
func PickAssessmentSlug(rt *runtime.Runtime) (string, error) {
	return PickPaged(rt, "Select assessment", "no assessments found",
		func(q url.Values) (any, error) { return rt.API.AssessmentsList(rt.Context(), q) },
		func(m map[string]any) (ui.Choice, bool) {
			slug := fieldString(m, "slug")
			if slug == "" {
				return ui.Choice{}, false
			}
			title := fieldString(m, "title", "name")
			status := fieldString(m, "status")
			label := slug
			if title != "" {
				label = fmt.Sprintf("%s — %s", slug, title)
			}
			if status != "" {
				label = label + " [" + status + "]"
			}
			return ui.Choice{Label: label, Value: slug}, true
		},
	)
}

// PickInviteToken lists invites (with Load more) and returns invite_token.
func PickInviteToken(rt *runtime.Runtime) (string, error) {
	return PickPaged(rt, "Select invitation", "no invitations found",
		func(q url.Values) (any, error) { return rt.API.InvitesList(rt.Context(), q) },
		func(m map[string]any) (ui.Choice, bool) {
			token := fieldString(m, "invite_token", "token", "id")
			if token == "" {
				return ui.Choice{}, false
			}
			email := fieldString(m, "email")
			status := fieldString(m, "status")
			label := token
			if len(token) > 12 {
				label = token[:12] + "…"
			}
			if email != "" {
				label = email + " · " + label
			}
			if status != "" {
				label = label + " [" + status + "]"
			}
			return ui.Choice{Label: label, Value: token}, true
		},
	)
}

// PickPipelineSlug lists pipelines (with Load more) and returns a slug.
func PickPipelineSlug(rt *runtime.Runtime) (string, error) {
	return PickPaged(rt, "Select pipeline", "no pipelines found",
		func(q url.Values) (any, error) { return rt.API.PipelinesList(rt.Context(), q) },
		func(m map[string]any) (ui.Choice, bool) {
			slug := fieldString(m, "slug")
			if slug == "" {
				return ui.Choice{}, false
			}
			name := fieldString(m, "name", "title")
			label := slug
			if name != "" {
				label = fmt.Sprintf("%s — %s", slug, name)
			}
			return ui.Choice{Label: label, Value: slug}, true
		},
	)
}

// PickWebhookID lists webhooks (with Load more) and returns an id.
func PickWebhookID(rt *runtime.Runtime) (string, error) {
	return PickPaged(rt, "Select webhook", "no webhooks found",
		func(q url.Values) (any, error) { return rt.API.WebhooksList(rt.Context(), q) },
		func(m map[string]any) (ui.Choice, bool) {
			id := fieldString(m, "id")
			if id == "" {
				return ui.Choice{}, false
			}
			urlStr := fieldString(m, "url", "endpoint_url")
			label := id
			if len(id) > 8 {
				label = id[:8]
			}
			if urlStr != "" {
				label = label + " · " + urlStr
			}
			return ui.Choice{Label: label, Value: id}, true
		},
	)
}

// ConfirmDestructive confirms unless --yes (arrow-key Yes/No when TTY).
func ConfirmDestructive(rt *runtime.Runtime, msg string) error {
	ok, err := ui.ConfirmTUI(rt.UI, msg)
	if err != nil {
		return err
	}
	if !ok {
		return &api.AbortError{}
	}
	return nil
}
