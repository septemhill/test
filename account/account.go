package account

import (
	"context"
	"time"

	"gopkg.in/guregu/null.v4"
)

type Account struct {
	ID       int         `db:"id"`
	Username string      `db:"username" json:"username" uri:"user"`
	Email    string      `db:"email" json:"email"`
	Password string      `db:"password" json:"password"`
	Phone    null.String `db:"phone" json:"phone"`
	CreateAt time.Time   `db:"create_at" json:"createAt"`
	UpdateAt time.Time   `db:"update_at" json:"updateAt"`
}

type AccountUseCase interface {
	Create(context.Context, *Account) error
	GetInfo(context.Context, int) (*Account, error)
	UpdateInfo(context.Context, *Account) error
	Delete(context.Context, int) error
	ChangePassword(context.Context, int) error
}

type AccountService interface {
	Create(context.Context, *Account) error
	GetInfo(context.Context, int) (*Account, error)
	UpdateInfo(context.Context, *Account) error
	Delete(context.Context, int) error
	ChangePassword(context.Context, int) error
}

type AccountRepository interface {
	Create(context.Context, *Account) (int, error)
	GetInfo(context.Context, int) (*Account, error)
	UpdateInfo(context.Context, *Account) (int, error)
	Delete(context.Context, int) (int, error)
	ChangePassword(context.Context, int) (int, error)
}
