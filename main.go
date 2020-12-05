package main

import (
	"context"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/septemhill/test/account/repository/postgresql"
)

func main() {
	ctx := context.Background()

	repo := postgresql.NewAccountRepository()
	_, _ = ctx, repo
}
