package scheduler

import "testing"

func TestNewScheduler(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatalf("new scheduler failed: %v", err)
	}
	if s == nil || s.cron == nil {
		t.Fatalf("scheduler not initialized")
	}
	entries := s.cron.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 cron jobs, got %d", len(entries))
	}
}
