package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/septemhill/test/account"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	uc := account.NewMockAccountUseCase(ctrl)
	uc.EXPECT().Create(ctx, &account.Account{}).Return(nil)
}
