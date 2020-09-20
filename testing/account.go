package testing

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/module"
	"github.com/septemhill/test/utils"
	"gopkg.in/guregu/null.v4"
)

type LoggedAccount struct {
	*module.Account
	Token string
}

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

	_, _ = module.CreateAccount(ctx, d, acc)

	if !withPassword {
		acc.Password = ""
	}

	return acc
}

func NewAccountWithLogin(ctx context.Context, d *db.DB, r *redis.Client) LoggedAccount {
	acc := NewAccount(ctx, d, true)
	tk, _ := module.Login(ctx, d, r, acc.Email, acc.Password)
	return LoggedAccount{Account: acc, Token: tk}
}

func DeleteAccounts(ctx context.Context, d *db.DB, accs ...*module.Account) {
	for _, acc := range accs {
		_, _ = module.DeleteAccount(ctx, d, acc)
	}
}

func DeleteLoggedAccounts(ctx context.Context, d *db.DB, accs ...LoggedAccount) {
	for _, acc := range accs {
		_, _ = module.DeleteAccount(ctx, d, acc.Account)
	}
}
