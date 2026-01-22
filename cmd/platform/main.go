package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/core/auth"
	"github.com/hoppermq/streamly/internal/core/platform"
	"github.com/hoppermq/streamly/internal/core/platform/bootstrap"
	"github.com/hoppermq/streamly/internal/core/platform/membership"
	"github.com/hoppermq/streamly/internal/core/platform/organization"
	"github.com/hoppermq/streamly/internal/core/platform/user"
	"github.com/hoppermq/streamly/internal/http"
	"github.com/hoppermq/streamly/internal/http/routes"
	"github.com/hoppermq/streamly/internal/storage/postgres"
	"github.com/hoppermq/streamly/pkg/shared/zitadel/client"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/zixyos/glog"
	serviceloader "github.com/zixyos/goloader/service"
)

func main() {
	logger, err := glog.NewDefault()
	if err != nil {
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		)
		os.Exit(84)
	}

	ctx := context.Background()
	platformConf, err := config.LoadPlatformConfig()
	if err != nil {
		logger.Warn("failed to load platform config", "error", err)
		os.Exit(84)
	}

	logger.InfoContext(ctx, "starting platform service")

	sqldb := sql.OpenDB(
		pgdriver.NewConnector(pgdriver.WithDSN(platformConf.DatabaseDSN())),
	)

	db := bun.NewDB(sqldb, pgdialect.New())

	d := postgres.NewClient(
		postgres.WithLogger(logger),
		postgres.WithDB(db),
	)

	if err = d.Bootstrap(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to bootstrap database", "error", err)
		os.Exit(84)
	}

	orgRepos, err := organization.NewRepository(
		organization.RepositoryWithLogger(logger),
		organization.RepositoryWithDB(db),
	)

	if err != nil {
		logger.ErrorContext(ctx, "failed to create organization repository", "error", err)
		os.Exit(84)
	}

	userRepo, err := user.NewRepository(
		user.RepositoryWithLogger(logger),
		user.RepositoryWithDB(db),
	)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create user repository", "error", err)
		os.Exit(84)
	}

	authRepo, err := auth.NewRepository(
		auth.RepositoryWithLogger(logger),
		auth.RepositoryWithDB(db),
	)

	if err != nil {
		logger.ErrorContext(ctx, "failed to create auth repository", "error", err)
		os.Exit(84)
	}

	membershipRepo := membership.NewRepository(
		membership.RepositoryWithLogger(logger),
		membership.RepositoryWithDB(db),
	)

	// Create UnitOfWork factory for transaction management
	uowFactory := postgres.NewUnitOfWorkFactory(
		postgres.FactoryWithDB(db),
		postgres.FactoryWithOrgRepo(orgRepos),
		postgres.FactoryWithMembershipRepo(membershipRepo),
		postgres.FactoryWithUserRepo(userRepo),
	)

	zitadelClient, err := client.NewZitadelClient(
		ctx,
		client.NewZitadel(
			platformConf.Platform.Zitadel.Domain,
			client.WithPort(platformConf.Platform.Zitadel.Port),
			client.WithInsecure(strconv.Itoa(int(platformConf.Platform.Zitadel.Port))),
		),
		client.WithLogger(logger),
		client.WithPATFromFile(platformConf.ZitadelPATPath()),
	)

	if err != nil {
		logger.ErrorContext(ctx, "failed to create zitadel client", "error", err)
		os.Exit(84)
	}

	generator := uuid.New
	uuidParser := uuid.Parse

	userUC, err := user.NewUseCase(
		user.UseCaseWithLogger(logger),
		user.UseCaseWithUserRepository(userRepo),
		user.UseCaseWithAuthRepository(authRepo),
		user.UseCaseWithGenerator(generator),
		user.UseCaseWithUUIDParser(uuidParser),
		user.UseCaseWithZitadelAPI(zitadelClient),
	)

	if err != nil {
		logger.ErrorContext(ctx, "failed to create user usecase", "error", err)
		os.Exit(84)
	}

	membershipUC := membership.NewUseCase(
		membership.UseCaseWithLogger(logger),
		membership.UseCaseWithRepository(membershipRepo),
		membership.UseCaseWithUUIDParser(uuidParser),
		membership.UseCaseWithGenerator(generator),
	)

	organizationUC, err := organization.NewUseCase(
		organization.UseCaseWithLogger(logger),
		organization.UseCaseWithRepository(orgRepos),
		organization.UseCaseWithGenerator(generator),
		organization.UseCaseWithUUIDParser(uuidParser),
		organization.UseCaseWithMembershipUC(membershipUC),
		organization.UseCaseWithUserUC(userUC),
		organization.UseCaseWithUOW(uowFactory),
	)

	if err != nil {
		logger.ErrorContext(ctx, "failed to create organization usecase", "error", err)
		os.Exit(84)
	}

	engine := gin.New()
	httpServer := http.NewHTTPServer(
		http.WithEngine(engine),
		http.WithPlatformHTTPServer(platformConf),
		http.WithLogger(logger),
		http.WithRoutes(
			routes.CreateRouteRegistrar(
				routes.CreatePlatformRegistrar(logger, organizationUC),
				routes.CreateWebhookRegistrar(logger, userUC),
			),
		),
	)

	platformService := platform.NewStreamlyService(
		platform.WithLogger(logger),
		platform.WithHandler(httpServer),
	)

	bootstrapOrchestrator := bootstrap.NewOrchestrator(
		bootstrap.BstWithLogger(logger),
		bootstrap.BstWithZitadel(*zitadelClient),
		bootstrap.BstWithOrgUC(organizationUC),
		bootstrap.BstWithUserUC(userUC),
	)

	if err := bootstrapOrchestrator.Run(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to bootstrap orchestrator", "error", err)
		os.Exit(84)
	}

	app := serviceloader.New(
		serviceloader.WithLogger(logger),
		serviceloader.WithService(platformService),
	)

	app.Run(ctx)
}
