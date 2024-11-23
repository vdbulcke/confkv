package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/vdbulcke/confkv/src/controller"
	"github.com/vdbulcke/confkv/src/logger"
	"go.uber.org/zap"
)

type Scheduler struct {
	cronTab string

	ctrl *controller.Controller
	cron *cron.Cron

	logger *zap.Logger
}

func NewScheduler(crontab string, ctrl *controller.Controller, l *zap.Logger) (*Scheduler, error) {
	cron := cron.New()
	s := &Scheduler{
		cronTab: crontab,
		ctrl:    ctrl,
		cron:    cron,
		logger:  l,
	}

	entry, err := cron.AddFunc(crontab, s.Job)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("cronjob created", zap.Int("entryID", int(entry)))

	s.logger.Info("created scheduler", zap.Any("opts", crontab))

	return s, nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Shutdown(ctx context.Context) {
	jobCtx := s.cron.Stop()

	// wait for either shutdown ctx cancel
	// of jobCtx completes
	select {
	case <-ctx.Done():
		return
	case <-jobCtx.Done():
		return
	}

}

func (s *Scheduler) Job() {

	ctx := logger.ContextWithJobID(context.Background())

	l := logger.GetLoggerWithContext(ctx, s.logger)

	l.Info("starting job")

	outcome := "success"
	err := s.ctrl.SyncJob(ctx)
	if err != nil {
		l.Error("job error", zap.Error(err))
		outcome = "failure"
	}

	l.Info("job finished", zap.String("status", outcome))
}
