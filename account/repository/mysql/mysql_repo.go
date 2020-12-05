package mysql

import (
	"context"
	"fmt"
)

type accountRepository struct{}

func NewAccountRepository() *accountRepository {
	return &accountRepository{}
}

func (repo *accountRepository) Create(ctx context.Context) error {
	fmt.Println("mysql create")
	return nil
}
func (repo *accountRepository) GetInfo(ctx context.Context) error {
	fmt.Println("mysql getinfo")
	return nil
}
func (repo *accountRepository) UpdateInfo(ctx context.Context) error {
	fmt.Println("mysql updateinfo")
	return nil
}
func (repo *accountRepository) Delete(ctx context.Context) error {
	fmt.Println("mysql delete")
	return nil
}
func (repo *accountRepository) ChangePassword(ctx context.Context) error {
	fmt.Println("mysql changepasswd")
	return nil
}
