package usecase

import (
	"context"

	"github.com/septemhill/test/account"
)

type accountUseCase struct {
	repo account.AccountRepository
}

func NewAccountUseCase(repo account.AccountRepository) *accountUseCase {
	return &accountUseCase{
		repo: repo,
	}
}

func (u *accountUseCase) Create(ctx context.Context, acc *account.Account) error {
	_, err := u.repo.Create(ctx, acc)
	return err
}

func (u *accountUseCase) GetInfo(ctx context.Context, id int) (*account.Account, error) {
	_, err := u.repo.GetInfo(ctx, id)
	return err
}

func (u *accountUseCase) UpdateInfo(ctx context.Context, acc *account.Account) error {
	_, err := u.repo.UpdateInfo(ctx, acc)
	return err
}

func (u *accountUseCase) Delete(ctx context.Context, id int) error {
	_, err := u.repo.Delete(ctx, id)
	return err
}

func (u *accountUseCase) ChangePassword(ctx context.Context, id int, passwd string) error {
	_, err := u.repo.ChangePassword(ctx, id, passwd)
	return err
}

var _ account.AccountUseCase = (*accountUseCase)(nil)
