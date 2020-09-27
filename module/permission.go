package module

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
)

type RoleType int
type PermissionType int

const (
	ROLE_ADMIN RoleType = iota
	ROLE_ACCOUNT_MGR
	ROLE_NORMAL_USER
)

const (
	ACCOUNT_PERM PermissionType = iota
	BLOG_PERM
)

type Role struct {
	//Name  string
	Perms []PermissionType
}

func CreateRole(ctx context.Context, d *db.DB, role Role, perms ...PermissionType) ([]int, error) {
	tpl := `INSERT INTO permissions(role, perm) VALUES %s RETURNING id`

	v := fmt.Sprintf(`(%d, %d)`, role, perms[0])
	for i := 1; i < len(perms); i++ {
		v += fmt.Sprintf(`, (%d, %d)`, role, perms[i])
	}

	expr := fmt.Sprintf(tpl, v)
	ids := []int{}

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.SelectContext(ctx, &ids, expr); err != nil {
			return err
		}
		return nil
	})

	return ids, err
}

func DeleteRole(ctx context.Context, d *db.DB, role Role) ([]int, error) {
	expr := `DELETE FROM permissions WHERE role = $1`
	ids := []int{}

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.SelectContext(ctx, &ids, expr, role); err != nil {
			return err
		}
		return nil
	})

	return ids, err
}

func GetRoles(ctx context.Context, d *db.DB) ([]Role, error) {
	expr := `SELECT DISTINCT role FROM permissions`
	roles := []Role{}

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.SelectContext(ctx, &roles, expr); err != nil {
			return err
		}
		return nil
	})

	return roles, err
}

func GetRole(ctx context.Context, d *db.DB, role RoleType) (*Role, error) {
	expr := `SELECT perms FROM permissions WHERE role = $1`
	r := &Role{}

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, r, expr, role); err != nil {
			return err
		}

		return nil
	})

	return r, err
}

func UpdateRolePermissions(ctx context.Context, d *db.DB) {
}

func AddRole(ctx context.Context, d *db.DB) {
}

func RemoveRole(ctx context.Context, d *db.DB) {}

func CreatePermission(ctx context.Context, d *db.DB) {
}

func DeletePermission(ctx context.Context, d *db.DB) {
}

func GetPermissions(ctx context.Context, d *db.DB) {}
