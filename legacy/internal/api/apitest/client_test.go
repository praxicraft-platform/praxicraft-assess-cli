package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/praxicraft-platform/praxicraft-assess-cli/internal/api"
)

func TestOrgGetAndErrorMapping(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/public/org/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer ct_test_x" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"code": "AUTH_ERROR", "message": "bad key"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"name": "Acme", "slug": "acme"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c, err := api.New("ct_test_x", srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	c.HTTPClient = srv.Client()

	out, err := c.OrgGet(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	m := out.(map[string]any)
	if m["slug"] != "acme" {
		t.Fatalf("got %#v", m)
	}

	c2, _ := api.New("bad", srv.URL)
	c2.HTTPClient = srv.Client()
	_, err = c2.OrgGet(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if api.ExitCode(err) != api.ExitAuth {
		t.Fatalf("exit code %d", api.ExitCode(err))
	}
}

func TestRetryOn503(t *testing.T) {
	attempts := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/public/org/stats/", func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"invites_remaining": 3})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c, _ := api.New("ct_test_x", srv.URL)
	c.HTTPClient = srv.Client()
	c.MaxRetries = 2
	out, err := c.OrgStats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("attempts=%d", attempts)
	}
	_ = out
}
