package api

import (
	"context"
	"net/url"
)

// --- Org ---

func (c *Client) OrgGet(ctx context.Context) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/org/", nil, &out)
}

func (c *Client) OrgStats(ctx context.Context) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/org/stats/", nil, &out)
}

func (c *Client) OrgTeam(ctx context.Context) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/org/team/", nil, &out)
}

func (c *Client) OrgAuditLog(ctx context.Context, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/org/audit-log/", q, &out)
}

func (c *Client) OrgSquadsList(ctx context.Context) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/org/squads/", nil, &out)
}

func (c *Client) OrgSquadGet(ctx context.Context, teamID string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/org/squads/"+url.PathEscape(teamID)+"/", nil, &out)
}

func (c *Client) OrgSquadMembers(ctx context.Context, teamID string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/org/squads/"+url.PathEscape(teamID)+"/members/", nil, &out)
}

// --- Assessments ---

func (c *Client) AssessmentsList(ctx context.Context, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/assessments/", q, &out)
}

func (c *Client) AssessmentsGet(ctx context.Context, slug string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/assessments/"+pathSeg(slug)+"/", nil, &out)
}

func (c *Client) AssessmentsCreate(ctx context.Context, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/assessments/create/", body, &out)
}

func (c *Client) AssessmentsUpdate(ctx context.Context, slug string, body any) (any, error) {
	var out any
	return out, c.PatchJSON(ctx, "/assessments/"+pathSeg(slug)+"/update/", body, &out)
}

func (c *Client) AssessmentsDuplicate(ctx context.Context, slug string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/assessments/"+pathSeg(slug)+"/duplicate/", body, &out)
}

func (c *Client) AssessmentsCasesList(ctx context.Context, slug string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/assessments/"+pathSeg(slug)+"/cases/", nil, &out)
}

func (c *Client) AssessmentsCasesAttach(ctx context.Context, slug string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/assessments/"+pathSeg(slug)+"/cases/attach/", body, &out)
}

func (c *Client) AssessmentsCasesReplace(ctx context.Context, slug string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/assessments/"+pathSeg(slug)+"/cases/replace/", body, &out)
}

func (c *Client) AssessmentsCasesRemove(ctx context.Context, slug string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/assessments/"+pathSeg(slug)+"/cases/remove/", body, &out)
}

func (c *Client) AssessmentsResults(ctx context.Context, slug string, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/assessments/"+pathSeg(slug)+"/results/", q, &out)
}

// --- Invites ---

func (c *Client) InvitesList(ctx context.Context, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/invites/", q, &out)
}

func (c *Client) InvitesCreate(ctx context.Context, slug string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/assessments/"+pathSeg(slug)+"/invites/", body, &out)
}

func (c *Client) InvitesBulkCreate(ctx context.Context, slug string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/assessments/"+pathSeg(slug)+"/invites/bulk/", body, &out)
}

func (c *Client) InvitesGet(ctx context.Context, token string) (any, error) {
	// Detail is via cancel path GET isn't separate — use result or list.
	// Public API: GET cancel path doesn't exist; DELETE cancels. Retrieve via list filter or result.
	var out any
	return out, c.GetJSON(ctx, "/invites/"+url.PathEscape(token)+"/result/", nil, &out)
}

func (c *Client) InvitesRemind(ctx context.Context, token string) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/invites/"+url.PathEscape(token)+"/remind/", map[string]any{}, &out)
}

func (c *Client) InvitesCancel(ctx context.Context, token string) (any, error) {
	var out any
	return out, c.DeleteJSON(ctx, "/invites/"+url.PathEscape(token)+"/", &out)
}

func (c *Client) InvitesResult(ctx context.Context, token string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/invites/"+url.PathEscape(token)+"/result/", nil, &out)
}

// --- Results ---

func (c *Client) ResultsList(ctx context.Context, assessmentSlug string, q url.Values) (any, error) {
	return c.AssessmentsResults(ctx, assessmentSlug, q)
}

func (c *Client) ResultsGet(ctx context.Context, token string) (any, error) {
	return c.InvitesResult(ctx, token)
}

// --- Cases ---

func (c *Client) CasesPlatformList(ctx context.Context, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/platform-cases/", q, &out)
}

func (c *Client) CasesList(ctx context.Context, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/cases/", q, &out)
}

func (c *Client) CasesCreate(ctx context.Context, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/cases/create/", body, &out)
}

func (c *Client) CasesGet(ctx context.Context, id string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/cases/"+url.PathEscape(id)+"/", nil, &out)
}

func (c *Client) CasesUpdate(ctx context.Context, id string, body any) (any, error) {
	var out any
	return out, c.PatchJSON(ctx, "/cases/"+url.PathEscape(id)+"/", body, &out)
}

func (c *Client) CasesDelete(ctx context.Context, id string) (any, error) {
	var out any
	return out, c.DeleteJSON(ctx, "/cases/"+url.PathEscape(id)+"/", &out)
}

// --- Pipelines ---

func (c *Client) PipelinesList(ctx context.Context, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/pipelines/", q, &out)
}

func (c *Client) PipelinesGet(ctx context.Context, slug string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/pipelines/"+pathSeg(slug)+"/", nil, &out)
}

func (c *Client) PipelinesEnroll(ctx context.Context, slug string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/pipelines/"+pathSeg(slug)+"/enroll/", body, &out)
}

func (c *Client) PipelinesBulkEnroll(ctx context.Context, slug string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/pipelines/"+pathSeg(slug)+"/enroll/bulk/", body, &out)
}

func (c *Client) PipelinesEnrollments(ctx context.Context, slug string, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/pipelines/"+pathSeg(slug)+"/enrollments/", q, &out)
}

func (c *Client) PipelinesGetEnrollment(ctx context.Context, id string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/pipelines/enrollments/"+url.PathEscape(id)+"/", nil, &out)
}

func (c *Client) PipelinesReject(ctx context.Context, id string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/pipelines/enrollments/"+url.PathEscape(id)+"/reject/", body, &out)
}

func (c *Client) PipelinesHold(ctx context.Context, id string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/pipelines/enrollments/"+url.PathEscape(id)+"/hold/", body, &out)
}

func (c *Client) PipelinesUnhold(ctx context.Context, id string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/pipelines/enrollments/"+url.PathEscape(id)+"/unhold/", body, &out)
}

// --- Webhooks ---

func (c *Client) WebhooksList(ctx context.Context) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/webhooks/", nil, &out)
}

func (c *Client) WebhooksCreate(ctx context.Context, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/webhooks/create/", body, &out)
}

func (c *Client) WebhooksGet(ctx context.Context, id string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/webhooks/"+url.PathEscape(id)+"/", nil, &out)
}

func (c *Client) WebhooksUpdate(ctx context.Context, id string, body any) (any, error) {
	var out any
	return out, c.PatchJSON(ctx, "/webhooks/"+url.PathEscape(id)+"/", body, &out)
}

func (c *Client) WebhooksDelete(ctx context.Context, id string) (any, error) {
	var out any
	return out, c.DeleteJSON(ctx, "/webhooks/"+url.PathEscape(id)+"/", &out)
}

func (c *Client) WebhooksTest(ctx context.Context, id string) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/webhooks/"+url.PathEscape(id)+"/test/", map[string]any{}, &out)
}

func (c *Client) WebhooksDeliveries(ctx context.Context, id string, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/webhooks/"+url.PathEscape(id)+"/deliveries/", q, &out)
}

func (c *Client) WebhooksRetryDelivery(ctx context.Context, id, deliveryID string) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/webhooks/"+url.PathEscape(id)+"/deliveries/"+url.PathEscape(deliveryID)+"/retry/", map[string]any{}, &out)
}

// --- Interviews ---

func (c *Client) InterviewsList(ctx context.Context, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/interviews/", q, &out)
}

func (c *Client) InterviewsCreate(ctx context.Context, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/interviews/create/", body, &out)
}

func (c *Client) InterviewsBulkCreate(ctx context.Context, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/interviews/bulk/", body, &out)
}

func (c *Client) InterviewsAnalytics(ctx context.Context, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/interviews/analytics/", q, &out)
}

func (c *Client) InterviewsOrgCases(ctx context.Context, q url.Values) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/interviews/org-cases/", q, &out)
}

func (c *Client) InterviewsGet(ctx context.Context, roomID string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/interviews/"+url.PathEscape(roomID)+"/", nil, &out)
}

func (c *Client) InterviewsCancel(ctx context.Context, roomID string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/interviews/"+url.PathEscape(roomID)+"/cancel/", body, &out)
}

func (c *Client) InterviewsReschedule(ctx context.Context, roomID string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/interviews/"+url.PathEscape(roomID)+"/reschedule/", body, &out)
}

func (c *Client) InterviewsAnalysis(ctx context.Context, roomID string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/interviews/"+url.PathEscape(roomID)+"/analysis/", nil, &out)
}

func (c *Client) InterviewsReplay(ctx context.Context, roomID string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/interviews/"+url.PathEscape(roomID)+"/replay/", nil, &out)
}

func (c *Client) InterviewsShare(ctx context.Context, roomID string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/interviews/"+url.PathEscape(roomID)+"/share/", body, &out)
}

func (c *Client) InterviewTemplatesList(ctx context.Context) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/interviews/templates/", nil, &out)
}

func (c *Client) InterviewTemplatesCreate(ctx context.Context, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/interviews/templates/create/", body, &out)
}

func (c *Client) InterviewTemplatesGet(ctx context.Context, id string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/interviews/templates/"+url.PathEscape(id)+"/", nil, &out)
}

func (c *Client) InterviewTemplatesUpdate(ctx context.Context, id string, body any) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/interviews/templates/"+url.PathEscape(id)+"/update/", body, &out)
}

func (c *Client) InterviewTemplatesDelete(ctx context.Context, id string) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/interviews/templates/"+url.PathEscape(id)+"/delete/", map[string]any{}, &out)
}

// --- Integrations ---

func (c *Client) IntegrationsList(ctx context.Context) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/integrations/", nil, &out)
}

func (c *Client) IntegrationsConnectURL(ctx context.Context, provider string) (any, error) {
	var out any
	return out, c.GetJSON(ctx, "/integrations/"+pathSeg(provider)+"/connect/", nil, &out)
}

func (c *Client) IntegrationsTest(ctx context.Context, provider string) (any, error) {
	var out any
	return out, c.PostJSON(ctx, "/integrations/"+pathSeg(provider)+"/test/", map[string]any{}, &out)
}

func pathSeg(s string) string {
	// Keep path segments URL-safe; encode but preserve readability for slugs.
	return url.PathEscape(s)
}
