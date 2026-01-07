package role

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/internal/models"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/uptrace/bun"
)

type Repository struct {
	logger *slog.Logger
	db     bun.IDB
}

type RepositoryOption func(*Repository)

func RepositoryWithLogger(logger *slog.Logger) RepositoryOption {
	return func(r *Repository) {
		r.logger = logger
	}
}

func RepositoryWithDB(db bun.IDB) RepositoryOption {
	return func(r *Repository) {
		r.db = db
	}
}

func NewRepository(opts ...RepositoryOption) *Repository {
	r := &Repository{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Repository) WithTx(tx domain.TxContext) domain.RoleRepository {
	bunTx, ok := tx.(bun.IDB)
	if !ok {
		r.logger.Warn("transaction does not implement github.com/uptrace/bun")
		return r
	}

	return &Repository{
		logger: r.logger,
		db:     bunTx,
	}
}

func (r *Repository) Create(ctx context.Context, input domain.CreateTenantRoleInput) error {
	r.logger.InfoContext(ctx, "inserting a new role", "role_name", input.RoleName)
	role := &models.Role{
		Identifier:  input.Identifier,
		RoleName:    input.RoleName,
		Permissions: input.Permissions,
		Metadata:    input.Metadata,
	}

	res, err := r.db.NewInsert().Model(role).Exec(ctx)
	if err != nil {
		r.logger.WarnContext(ctx, "error inserting role", "role_name", input.RoleName, "error", err)
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		r.logger.WarnContext(ctx, "error inserting role", "role_name", input.RoleName, "error", err)
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, input domain.UpdateTenantRoleInput) error {
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*domain.TenantRole, error) {
	return nil, nil
}

func (r *Repository) FindAll(ctx context.Context, limit, offset int) ([]domain.TenantRole, error) {
	return make([]domain.TenantRole, 0), nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *Repository) Exist(ctx context.Context, id uuid.UUID) (bool, error) {
	return false, nil
}

func (r *Repository) Save(ctx context.Context, role models.Role) error {
	r.logger.InfoContext(ctx, "saving role", "role_identifier", role.Identifier, "role_name", role.RoleName)
	_, err := r.db.NewInsert().Model(role).Exec(ctx)
	if err != nil {
		r.logger.WarnContext(ctx, "error inserting role", "role_name", role.RoleName, "error", err)
		return err
	}

	return nil
}

func (r *Repository) SaveAll(ctx context.Context, roles []domain.TenantRole) error {
	newRoles := make([]models.Role, len(roles))
	for i, role := range roles {
		newRoles[i].Identifier = role.Identifier
		newRoles[i].RoleName = role.RoleName
		newRoles[i].Metadata = role.Metadata
		newRoles[i].Permissions = role.Permissions
	}

	res, err := r.db.NewInsert().Model(&newRoles).Exec(ctx)
	if err != nil {
		r.logger.WarnContext(ctx, "error inserting roles", "roles", newRoles, "error", err)
		return err
	}

	if ra, err := res.RowsAffected(); err != nil || int(ra) < len(newRoles) { // TODO: bug here should out.
		r.logger.WarnContext(ctx, "error inserting roles", "roles", newRoles, "error", err)
		return err
	}

	r.logger.InfoContext(ctx, "new roles inserted", "roles", newRoles)
	return nil
}
