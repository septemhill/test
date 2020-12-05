package usecase

import (
	"context"

	"github.com/septemhill/test/account"
	"github.com/septemhill/test/account/repository"
)

type accountUseCase struct {
	repo repository.AccountRepository
}

func NewAccountUseCase(repo repository.AccountRepository) *accountUseCase {
	return &accountUseCase{
		repo: repo,
	}
}

func (u *accountUseCase) Create(ctx context.Context) error {
	_, err := u.repo.Create(ctx)
	return err
}

func (u *accountUseCase) GetInfo(ctx context.Context) (*account.Account, error) {
	_, err := u.repo.GetInfo(ctx)
	return err
}

func (u *accountUseCase) UpdateInfo(ctx context.Context) error {
	_, err := u.repo.UpdateInfo(ctx)
	return err
}

func (u *accountUseCase) Delete(ctx context.Context) error {
	_, err := u.repo.Delete(ctx)
	return err
}

func (u *accountUseCase) ChangePassword(ctx context.Context, id int, passwd string) error {
	_, err := u.repo.ChangePassword(ctx)
	return err
}

var _ account.AccountUseCase = (*accountUseCase)(nil)
