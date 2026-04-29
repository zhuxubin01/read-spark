package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDictionaryLookup_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"word":"hello","phonetic":"həˈləʊ","meanings":[{"definitions":[{"definition":"used as a greeting"}]}]}]`))
	}))
	defer ts.Close()

	svc := NewDictionaryServiceWithClient(ts.URL, ts.Client())
	res, err := svc.Lookup(context.Background(), "Hello")
	if err != nil {
		t.Fatalf("lookup failed: %v", err)
	}
	if res.Word != "hello" || len(res.Meanings) != 1 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestDictionaryLookup_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	svc := NewDictionaryServiceWithClient(ts.URL, ts.Client())
	_, err := svc.Lookup(context.Background(), "missing")
	if err == nil || err.Error() != "word not found" {
		t.Fatalf("expected word not found, got %v", err)
	}
}
