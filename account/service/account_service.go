package service

import (
	"context"

	"github.com/septemhill/test/account"
)

type accountService struct {
	usecase account.AccountUseCase
}

func NewAccountService(uc account.AccountUseCase) *accountService {
	return &accountService{
		usecase: uc,
	}
}

func (s *accountService) Create(ctx context.Context, acc *account.Account) error {
	return s.usecase.Create(ctx, acc)
}

func (s *accountService) GetInfo(ctx context.Context, id int) (*account.Account, error) {
	return s.usecase.GetInfo(ctx, id)
}

func (s *accountService) UpdateInfo(ctx context.Context, acc *account.Account) error {
	return s.usecase.UpdateInfo(ctx, acc)
}

func (s *accountService) Delete(ctx context.Context, id int) error {
	return s.usecase.Delete(ctx, id)
}

func (s *accountService) ChangePassword(ctx context.Context, id int, passwd string) error {
	return s.usecase.ChangePassword(ctx, id, passwd)
}

var _ account.AccountService = (*accountService)(nil)
