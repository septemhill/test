package testing

import (
	"context"

	"github.com/septemhill/test/db"
	"github.com/septemhill/test/module"
	"github.com/septemhill/test/utils"
	"gopkg.in/guregu/null.v4"
)

func NewAccount(ctx context.Context, d *db.DB, withPassword bool) *module.Account {
	name := utils.GenerateRandomString(utils.RANDOM_ALL, 7)
	pass := utils.GenerateRandomString(utils.RANDOM_ALL, 12)
	phone := utils.GenerateRandomString(utils.RANDOM_DIGIT_ONLY, 10)

	acc := &module.Account{
		Username: name,
		Password: pass,
		Email:    name + "@balabababa.com",
		Phone:    null.StringFrom(phone),
	}

	_ = module.CreateAccount(ctx, d, *acc)

	if !withPassword {
		acc.Password = ""
	}

	return acc
}

func DeleteAccounts(ctx context.Context, d *db.DB, accs ...*module.Account) {
	for _, acc := range accs {
		_ = module.DeleteAccount(ctx, d, *acc)
	}
}
