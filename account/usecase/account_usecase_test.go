package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/septemhill/test/account"
	"github.com/septemhill/test/account/usecase"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := account.NewMockAccountRepository(ctrl)
	uc := usecase.NewAccountUseCase(repo)
	repo.EXPECT().Create(ctx, &account.Account{}).Return(nil)
	uc.Create(ctx, &account.Account{})
}
