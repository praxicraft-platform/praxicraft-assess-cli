package cmdutil

import (
	"fmt"
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

// ResolveOrPick returns args[0] or interactively picks from choices.
func ResolveOrPick(rt *runtime.Runtime, args []string, title string, load func() ([]ui.Choice, error)) (string, error) {
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		return strings.TrimSpace(args[0]), nil
	}
	if err := rt.EnsureAPI(); err != nil {
		return "", err
	}
	choices, err := load()
	if err != nil {
		return "", err
	}
	return ui.Select(rt.UI, title, choices)
}

// PickAssessmentSlug lists assessments and returns a slug.
func PickAssessmentSlug(rt *runtime.Runtime) (string, error) {
	return ResolveOrPick(rt, nil, "Select assessment", func() ([]ui.Choice, error) {
		raw, err := rt.API.AssessmentsList(rt.Context(), nil)
		if err != nil {
			return nil, err
		}
		rows := ResultsRows(raw)
		choices := make([]ui.Choice, 0, len(rows))
		for _, m := range rows {
			slug := fieldString(m, "slug")
			if slug == "" {
				continue
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
			choices = append(choices, ui.Choice{Label: label, Value: slug})
		}
		if len(choices) == 0 {
			return nil, &api.UsageError{Msg: "no assessments found"}
		}
		return choices, nil
	})
}

// PickInviteToken lists invites and returns invite_token.
func PickInviteToken(rt *runtime.Runtime) (string, error) {
	return ResolveOrPick(rt, nil, "Select invitation", func() ([]ui.Choice, error) {
		raw, err := rt.API.InvitesList(rt.Context(), nil)
		if err != nil {
			return nil, err
		}
		rows := ResultsRows(raw)
		choices := make([]ui.Choice, 0, len(rows))
		for _, m := range rows {
			token := fieldString(m, "invite_token", "token", "id")
			if token == "" {
				continue
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
			choices = append(choices, ui.Choice{Label: label, Value: token})
		}
		if len(choices) == 0 {
			return nil, &api.UsageError{Msg: "no invitations found"}
		}
		return choices, nil
	})
}

// PickPipelineSlug lists pipelines and returns a slug.
func PickPipelineSlug(rt *runtime.Runtime) (string, error) {
	return ResolveOrPick(rt, nil, "Select pipeline", func() ([]ui.Choice, error) {
		raw, err := rt.API.PipelinesList(rt.Context(), nil)
		if err != nil {
			return nil, err
		}
		rows := ResultsRows(raw)
		choices := make([]ui.Choice, 0, len(rows))
		for _, m := range rows {
			slug := fieldString(m, "slug")
			if slug == "" {
				continue
			}
			name := fieldString(m, "name", "title")
			label := slug
			if name != "" {
				label = fmt.Sprintf("%s — %s", slug, name)
			}
			choices = append(choices, ui.Choice{Label: label, Value: slug})
		}
		if len(choices) == 0 {
			return nil, &api.UsageError{Msg: "no pipelines found"}
		}
		return choices, nil
	})
}

// PickWebhookID lists webhooks and returns an id.
func PickWebhookID(rt *runtime.Runtime) (string, error) {
	return ResolveOrPick(rt, nil, "Select webhook", func() ([]ui.Choice, error) {
		raw, err := rt.API.WebhooksList(rt.Context())
		if err != nil {
			return nil, err
		}
		rows := ResultsRows(raw)
		if len(rows) == 0 {
			// Some APIs return a bare list under a different key or as array.
			if arr, ok := raw.([]any); ok {
				rows = anyMaps(arr)
			}
		}
		choices := make([]ui.Choice, 0, len(rows))
		for _, m := range rows {
			id := fieldString(m, "id")
			if id == "" {
				continue
			}
			url := fieldString(m, "url", "endpoint_url")
			label := id
			if len(id) > 8 {
				label = id[:8]
			}
			if url != "" {
				label = label + " · " + url
			}
			choices = append(choices, ui.Choice{Label: label, Value: id})
		}
		if len(choices) == 0 {
			return nil, &api.UsageError{Msg: "no webhooks found"}
		}
		return choices, nil
	})
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
