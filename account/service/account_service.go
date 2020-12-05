package service

import (
	"context"

	"github.com/septemhill/re/account"
)

type accountService struct {
	usecase account.AccountUseCase
}

func NewAccountService(uc account.AccountUseCase) *accountService {
	return &accountService{
		usecase: uc,
	}
}

func (s *accountService) Create(ctx context.Context) error {
	return s.usecase.Create(ctx)
}

func (s *accountService) GetInfo(ctx context.Context) (*account.Account, error) {
	return s.usecase.GetInfo(ctx)
}

func (s *accountService) UpdateInfo(ctx context.Context) error {
	return s.usecase.UpdateInfo(ctx)
}

func (s *accountService) Delete(ctx context.Context) error {
	return s.usecase.Delete(ctx)
}

func (s *accountService) ChangePassword(ctx context.Context) error {
	return s.usecase.ChangePassword(ctx)
}

var _ account.AccountService = (*accountService)(nil)
