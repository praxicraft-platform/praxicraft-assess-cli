package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultBaseURL   = "https://assess.praxicraft.com"
	DefaultAPIPrefix = "/api/v1/public"
	DefaultTimeout   = 30 * time.Second
	DefaultMaxRetries = 2
	maxRetryAfter    = 8 * time.Second
)

// Client talks to the Assess Public API.
type Client struct {
	APIKey     string
	BaseURL    string
	APIPrefix  string
	HTTPClient *http.Client
	MaxRetries int
	UserAgent  string
}

// New creates a Client. APIKey and BaseURL must be set by the caller (runtime resolves env/config).
func New(apiKey, baseURL string) (*Client, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, &APIError{
			Code: "MISSING_API_KEY",
			Message: "No API key provided. Run `praxicraft-assess configure`, pass --api-key, or set PRAXICRAFT_API_KEY. " +
				"Create a key at https://assess.praxicraft.com/assess/api — docs: https://docs.praxicraft.com/authentication",
		}
	}
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		APIPrefix:  DefaultAPIPrefix,
		MaxRetries: DefaultMaxRetries,
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		UserAgent:  "praxicraft-assess-cli",
	}, nil
}

// DoJSON performs an HTTP request and decodes JSON into out (may be nil for 204).
func (c *Client) DoJSON(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	var lastErr error
	attempts := c.MaxRetries + 1
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}

		u, err := c.buildURL(path, query)
		if err != nil {
			return err
		}

		var rdr io.Reader
		if payload != nil {
			rdr = bytes.NewReader(payload)
		}
		req, err := http.NewRequestWithContext(ctx, method, u, rdr)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.UserAgent)
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}

		if resp.StatusCode == http.StatusNoContent || len(respBody) == 0 {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				if out != nil {
					return json.Unmarshal([]byte("null"), out)
				}
				return nil
			}
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if out == nil {
				return nil
			}
			if err := json.Unmarshal(respBody, out); err != nil {
				return &APIError{Status: resp.StatusCode, Code: "INVALID_JSON", Message: "success response was not valid JSON", Body: respBody}
			}
			return nil
		}

		ae := parseError(resp.StatusCode, respBody)
		if shouldRetry(resp.StatusCode) && attempt < attempts-1 {
			if ra := retryAfter(resp.Header); ra > 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(ra):
				}
			}
			lastErr = ae
			continue
		}
		return ae
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("request failed after retries")
	}
	return lastErr
}

func (c *Client) buildURL(path string, query url.Values) (string, error) {
	path = "/" + strings.TrimLeft(path, "/")
	base := c.BaseURL + c.APIPrefix
	u, err := url.Parse(base + path)
	if err != nil {
		return "", err
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}
	return u.String(), nil
}

func parseError(status int, body []byte) *APIError {
	var envelope struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	_ = json.Unmarshal(body, &envelope)
	msg := envelope.Message
	if msg == "" {
		msg = envelope.Error
	}
	return mapStatusError(status, envelope.Code, msg, body)
}

func shouldRetry(status int) bool {
	switch status {
	case 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

func retryAfter(h http.Header) time.Duration {
	v := strings.TrimSpace(h.Get("Retry-After"))
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
		d := time.Duration(secs) * time.Second
		if d > maxRetryAfter {
			return maxRetryAfter
		}
		return d
	}
	return 0
}

func backoff(attempt int) time.Duration {
	// ~500ms base with jitter, exponential
	base := 500 * time.Millisecond
	d := time.Duration(float64(base) * math.Pow(2, float64(attempt-1)))
	jitter := time.Duration(rand.Int63n(int64(100 * time.Millisecond)))
	d += jitter
	if d > maxRetryAfter {
		return maxRetryAfter
	}
	return d
}

// GetJSON is a convenience GET.
func (c *Client) GetJSON(ctx context.Context, path string, query url.Values, out any) error {
	return c.DoJSON(ctx, http.MethodGet, path, query, nil, out)
}

// PostJSON is a convenience POST.
func (c *Client) PostJSON(ctx context.Context, path string, body any, out any) error {
	return c.DoJSON(ctx, http.MethodPost, path, nil, body, out)
}

// PatchJSON is a convenience PATCH.
func (c *Client) PatchJSON(ctx context.Context, path string, body any, out any) error {
	return c.DoJSON(ctx, http.MethodPatch, path, nil, body, out)
}

// PutJSON is a convenience PUT.
func (c *Client) PutJSON(ctx context.Context, path string, body any, out any) error {
	return c.DoJSON(ctx, http.MethodPut, path, nil, body, out)
}

// DeleteJSON is a convenience DELETE.
func (c *Client) DeleteJSON(ctx context.Context, path string, out any) error {
	return c.DoJSON(ctx, http.MethodDelete, path, nil, nil, out)
}
