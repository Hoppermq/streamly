package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/zixyos/glog"
	serviceloader "github.com/zixyos/goloader/service"

	"github.com/hoppermq/streamly/cmd/config"
	"github.com/hoppermq/streamly/internal/core/query"
	"github.com/hoppermq/streamly/internal/core/query/ast"
	"github.com/hoppermq/streamly/internal/http"
	"github.com/hoppermq/streamly/internal/http/routes"
	"github.com/hoppermq/streamly/internal/storage/clickhouse"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/schemas"
)

func main() {
	logger, err := glog.NewDefault()
	if err != nil {
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		).Error("failed to initialize logger", "error", err)
		os.Exit(domain.ExitStatus)
	}

	ctx := context.Background()

	queryConfig, err := config.LoadQueryConfig()
	if err != nil {
		logger.Warn("failed to load query config", "error", err)
	}

	clickhouseDriver := clickhouse.OpenConn(
		// TODO : refactor for multi tenant
		clickhouse.WithQueryConfig(queryConfig),
	)

	queryRepository := query.NewQueryRepository(
		query.WithDriver(clickhouseDriver),
	)

	astValidator := ast.NewValidator(
		ast.ValidatorWithLogger(logger),
	)

	jsonSchemaCompiler := jsonschema.NewCompiler()
	queryTranslator := clickhouse.NewTranslator(
		clickhouse.TranslatorWithLogger(logger),
	)

	astBuilder := ast.NewBuilder(
		ast.BuilderWithLogger(logger),
		ast.BuilderWithSchemaFS(schemas.FileFS),
		ast.BuilderWithValidator(astValidator),
		ast.BuilderWithJsonSchemaCompiler(jsonSchemaCompiler),
		ast.BuilderWithTranslator(queryTranslator),
	)

	queryUseCase := query.NewQueryUseCase(
		query.WithUseCaseLogger(logger),
		query.WithRepository(queryRepository),
		query.WithAstBuilder(astBuilder),
	)

	engine := gin.New()

	httpServer := http.NewHTTPServer(
		http.WithEngine(engine),
		http.WithQueryHTTPServer(queryConfig),
		http.WithLogger(logger),
		http.WithRoutes(
			routes.CreateRouteRegistrar(
				routes.CreateQueryRegistrar(logger, queryUseCase),
			),
		),
	)

	queryService := query.NewQueryService(
		query.WithLogger(logger),
		query.WithHandlers(httpServer, astBuilder),
	)

	app := serviceloader.New(
		serviceloader.WithLogger(logger),
		serviceloader.WithService(queryService),
	)

	app.Run(ctx)
}
