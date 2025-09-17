package migration

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Service struct {
	db            *sql.DB
	logger        *slog.Logger
	migrationPath string
}

type Option func(*Service)

func WithDB(db *sql.DB) Option {
	return func(s *Service) {
		s.db = db
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *Service) {
		s.logger = logger
	}
}

func WithMigrationPath(path string) Option {
	return func(s *Service) {
		s.migrationPath = path
	}
}

func NewService(opts ...Option) *Service {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Service) RunMigrations(ctx context.Context) error {
	s.logger.InfoContext(ctx, "initializing migration system")

	// Ensure migration table exists first
	if err := s.ensureMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	driver, err := clickhouse.WithInstance(s.db, &clickhouse.Config{
		MultiStatementEnabled: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create ClickHouse driver: %w", err)
	}

	s.logger.InfoContext(ctx, "creating migrate instance", "path", s.migrationPath)
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", s.migrationPath),
		"clickhouse",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	s.logger.InfoContext(ctx, "checking current migration version")
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		s.logger.ErrorContext(ctx, "failed to get version", "error", err)
		return fmt.Errorf("failed to get current version: %w", err)
	}

	s.logger.InfoContext(ctx, "migration status", "current_version", version, "dirty", dirty)

	// If migration is dirty, force it to clean state
	if dirty {
		s.logger.InfoContext(ctx, "migration is in dirty state, forcing clean", "version", version)
		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("failed to force clean migration state: %w", err)
		}
	}

	s.logger.InfoContext(ctx, "running migrations")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	newVersion, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to get new version: %w", err)
	}

	s.logger.InfoContext(ctx, "migrations completed", "version", newVersion, "dirty", dirty)
	return nil
}

func (s *Service) ensureMigrationTable(ctx context.Context) error {
	s.logger.InfoContext(ctx, "ensuring schema_migrations table exists")

	// Check if table exists
	var tableName string
	checkQuery := "SELECT name FROM system.tables WHERE database = currentDatabase() AND name = 'schema_migrations'"
	err := s.db.QueryRowContext(ctx, checkQuery).Scan(&tableName)

	if err == sql.ErrNoRows {
		// Table doesn't exist, create it
		s.logger.InfoContext(ctx, "creating schema_migrations table")
		if err := s.createMigrationTable(ctx); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if schema_migrations table exists: %w", err)
	} else {
		s.logger.InfoContext(ctx, "schema_migrations table already exists, checking schema compatibility")

		// Check if it has the right schema
		if compatible, err := s.isTableCompatible(ctx); err != nil {
			return fmt.Errorf("failed to check table compatibility: %w", err)
		} else if !compatible {
			s.logger.InfoContext(ctx, "existing schema_migrations table has incompatible schema, recreating")

			// Drop and recreate with correct schema
			if _, err := s.db.ExecContext(ctx, "DROP TABLE schema_migrations"); err != nil {
				return fmt.Errorf("failed to drop existing schema_migrations table: %w", err)
			}

			if err := s.createMigrationTable(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) createMigrationTable(ctx context.Context) error {
	createQuery := `
		CREATE TABLE schema_migrations (
			version    Int64,
			dirty      UInt8,
			sequence   UInt64
		) ENGINE = TinyLog
	`
	if _, err := s.db.ExecContext(ctx, createQuery); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	s.logger.InfoContext(ctx, "schema_migrations table created successfully")
	return nil
}

func (s *Service) isTableCompatible(ctx context.Context) (bool, error) {
	// Check if table has the expected columns
	rows, err := s.db.QueryContext(ctx, "DESCRIBE TABLE schema_migrations")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"version":  false,
		"dirty":    false,
		"sequence": false,
	}

	for rows.Next() {
		var name, type_, defaultType, defaultExpression, comment, codecExpression, ttlExpression string
		if err := rows.Scan(&name, &type_, &defaultType, &defaultExpression, &comment, &codecExpression, &ttlExpression); err != nil {
			return false, err
		}

		if _, exists := expectedColumns[name]; exists {
			expectedColumns[name] = true
		}
	}

	if err := rows.Err(); err != nil {
		s.logger.WarnContext(ctx, "failed to read schema_migrations table", "error", err)
		return false, err
	}

	// Check if all expected columns are present
	for col, found := range expectedColumns {
		if !found {
			s.logger.InfoContext(ctx, "missing expected column", "column", col)
			return false, nil
		}
	}

	return true, nil
}

