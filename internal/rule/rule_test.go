package rule

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_ = New()
}

func TestRuleAccess_IsValid(t *testing.T) {
	r := New()

	if !r.Access.IsValid() {
		t.Error("expected default Access to be valid")
	}
}

func TestRuleType_IsValid(t *testing.T) {
	r := New()

	if !r.Type.IsValid() {
		t.Error("expected default Type to be valid")
	}
}

func TestRule_Match(t *testing.T) {
	target := "http://example.com"

	req := httptest.NewRequest("GET", target, nil)
	bypassCache := map[string]map[string]time.Duration{}

	r := New()

	r.Pattern = "example\\.com"

	match, allow := r.Match(req, httptest.NewRecorder(), "", &bypassCache, time.Minute)

	if !match {
		t.Errorf("expected rule to match.  pattern: %v, target: %v, type: %v", r.Pattern, target, r.Type)
	}

	if allow {
		t.Errorf("expected the url to be blocked")
	}
}
