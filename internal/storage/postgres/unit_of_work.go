package postgres

import (
	"context"
	"database/sql"

	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

// UnitOfWork manages a transactional boundary across multiple repositories.
type UnitOfWork struct {
	tx bun.Tx

	orgRepo        domain.OrganizationRepository
	userRepo       domain.UserRepository
	membershipRepo domain.MembershipRepository
}

// UOWOptions defines functional options for UnitOfWork.
type UOWOptions func(*UnitOfWork)

func UowWithOrg(org domain.OrganizationRepository) UOWOptions {
	return func(u *UnitOfWork) {
		u.orgRepo = org
	}
}

func UowWithUser(user domain.UserRepository) UOWOptions {
	return func(u *UnitOfWork) {
		u.userRepo = user
	}
}

func UowWithMembership(membership domain.MembershipRepository) UOWOptions {
	return func(u *UnitOfWork) {
		u.membershipRepo = membership
	}
}

// NewUnitOfWork creates a new UnitOfWork instance with a transaction.
func NewUnitOfWork(ctx context.Context, db bun.IDB, options ...UOWOptions) (*UnitOfWork, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}

	uow := &UnitOfWork{
		tx: tx,
	}

	for _, option := range options {
		option(uow)
	}

	return uow, nil
}

func (u *UnitOfWork) Commit() error {
	return u.tx.Commit()
}

func (u *UnitOfWork) Rollback() error {
	return u.tx.Rollback()
}

func (u *UnitOfWork) Organization() domain.OrganizationRepository {
	return u.orgRepo.WithTx(u.tx)
}

func (u *UnitOfWork) User() domain.UserRepository {
	return u.userRepo.WithTx(u.tx)
}

func (u *UnitOfWork) Membership() domain.MembershipRepository {
	return u.membershipRepo.WithTx(u.tx)
}

// UnitOfWorkFactory creates new UnitOfWork instances (one per request).
type UnitOfWorkFactory struct {
	db bun.IDB

	orgRepoBase        domain.OrganizationRepository
	membershipRepoBase domain.MembershipRepository
	userRepoBase       domain.UserRepository
}

// UnitOfWorkFactoryOption defines functional options for UnitOfWorkFactory.
type UnitOfWorkFactoryOption func(*UnitOfWorkFactory)

func FactoryWithDB(db bun.IDB) UnitOfWorkFactoryOption {
	return func(f *UnitOfWorkFactory) {
		f.db = db
	}
}

func FactoryWithOrgRepo(repo domain.OrganizationRepository) UnitOfWorkFactoryOption {
	return func(f *UnitOfWorkFactory) {
		f.orgRepoBase = repo
	}
}

func FactoryWithMembershipRepo(repo domain.MembershipRepository) UnitOfWorkFactoryOption {
	return func(f *UnitOfWorkFactory) {
		f.membershipRepoBase = repo
	}
}

func FactoryWithUserRepo(repo domain.UserRepository) UnitOfWorkFactoryOption {
	return func(f *UnitOfWorkFactory) {
		f.userRepoBase = repo
	}
}

// NewUnitOfWorkFactory creates a factory that produces UnitOfWork instances.
func NewUnitOfWorkFactory(options ...UnitOfWorkFactoryOption) domain.UnitOfWorkFactory {
	factory := &UnitOfWorkFactory{}

	for _, option := range options {
		option(factory)
	}

	return factory
}

// NewUnitOfWork creates a fresh transaction-scoped UnitOfWork.
func (f *UnitOfWorkFactory) NewUnitOfWork(ctx context.Context) (domain.UnitOfWork, error) {
	return NewUnitOfWork(
		ctx,
		f.db,
		UowWithOrg(f.orgRepoBase),
		UowWithUser(f.userRepoBase),
		UowWithMembership(f.membershipRepoBase),
	)
}
