package bootstrap

import (
	"context"
	"log/slog"

	"github.com/hoppermq/streamly/internal/core/platform/organization"
	"github.com/hoppermq/streamly/internal/core/platform/user"
	"github.com/hoppermq/streamly/pkg/shared/zitadel/client"
)

type Orchestrator struct {
	logger  *slog.Logger
	zitadel client.Zitadel
	userUC  *user.UseCase
	orgUC   *organization.UseCase
}

type Options func(*Orchestrator)

func BstWithLogger(logger *slog.Logger) Options {
	return func(o *Orchestrator) {
		o.logger = logger
	}
}

func BstWithZitadel(zitadel client.Zitadel) Options {
	return func(o *Orchestrator) {
		o.zitadel = zitadel
	}
}

func BstWithUserUC(userUC *user.UseCase) Options {
	return func(o *Orchestrator) {
		o.userUC = userUC
	}
}

func BstWithOrgUC(orgUC *organization.UseCase) Options {
	return func(o *Orchestrator) {
		o.orgUC = orgUC
	}
}

func NewOrchestrator(opts ...Options) Orchestrator {
	orc := Orchestrator{}

	for _, opt := range opts {
		opt(&orc)
	}

	return orc
}

// Run will execute all bootstrap methods at the startup before the application process.
func (o *Orchestrator) Run(ctx context.Context) error {
	o.logger.Info("starting orchestrator")
	if err := o.setupDefaultOrg(ctx); err != nil {
		return err
	}
	return nil
}
