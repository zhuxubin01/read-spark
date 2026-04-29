package scheduler

import (
	"log/slog"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func New() (*Scheduler, error) {
	c := cron.New(cron.WithSeconds())

	// Every day at 08:00: publish daily articles (MVP: log only).
	if _, err := c.AddFunc("0 0 8 * * *", func() {
		slog.Info("scheduler: publish daily articles")
	}); err != nil {
		return nil, err
	}

	// Every day at 08:05: push notifications (MVP: log only).
	if _, err := c.AddFunc("0 5 8 * * *", func() {
		slog.Info("scheduler: send notification batch")
	}); err != nil {
		return nil, err
	}

	// Every hour: sync subscription status (MVP: log only).
	if _, err := c.AddFunc("0 0 * * * *", func() {
		slog.Info("scheduler: sync subscription status")
	}); err != nil {
		return nil, err
	}

	return &Scheduler{cron: c}, nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
