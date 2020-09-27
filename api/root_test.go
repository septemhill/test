package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/module"
	test "github.com/septemhill/test/testing"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), RootService, AccountService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()

	users := []*module.Account{
		test.NewAccount(ctx, d, true),
	}

	defer func() {
		test.DeleteAccounts(ctx, d, users...)
	}()

	tests := []struct {
		Description string
		Account     *module.Account
		StatusCode  int
	}{
		{
			Description: "Login a existed account",
			Account:     users[0],
			StatusCode:  http.StatusOK,
		},
		{
			Description: "Login a non-exist account",
			Account: &module.Account{
				Username: "nonexit_account",
				Password: "never_existed",
			},
			StatusCode: http.StatusBadRequest,
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		b, err := json.Marshal(&test.Account)
		asserts.NoError(err)

		req, err := http.NewRequest("POST", ts.URL+"/login", bytes.NewBuffer(b))
		asserts.NoError(err)

		rsp, err := http.DefaultClient.Do(req)
		asserts.NoError(err)
		defer rsp.Body.Close()

		asserts.Equal(test.StatusCode, rsp.StatusCode)
	}
}
