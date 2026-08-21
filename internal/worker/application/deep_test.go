package application

import (
	"context"
	worker "github.com/example/regional-authoritative-dns/internal/worker/domain"
	"log/slog"
	"testing"
	"time"
)

type deepRepo struct{ jobs []worker.Job }

func (r *deepRepo) Save(_ context.Context, j worker.Job) error {
	r.jobs = append(r.jobs, j)
	return nil
}
func (r *deepRepo) List(context.Context) ([]worker.Job, error) {
	return append([]worker.Job{}, r.jobs...), nil
}
func TestRunnerDoesNotExecuteJobsBeforeNextRun(t *testing.T) {
	r := &deepRepo{jobs: []worker.Job{{ID: "j", Kind: "publish", Status: "queued", NextRun: time.Now().Add(time.Hour)}}}
	s := New(r, slog.Default())
	called := false
	s.Handlers["publish"] = func(context.Context, worker.Job) error { called = true; return nil }
	s.once(context.Background())
	if called {
		t.Fatal("future job ran")
	}
}
