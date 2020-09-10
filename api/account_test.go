package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/septemhill/test/module"
	test "github.com/septemhill/test/testing"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestCreateAccount(t *testing.T) {
	ts := test.NewTestRouter(gin.Default(), AccountService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()

	asserts := assert.New(t)

	tests := []struct {
		Description string
		Account     module.Account
		StatusCode  int
		Clean       bool
	}{
		{
			Description: "New username and new email",
			Account: module.Account{
				Username: "user0001",
				Password: "user0001",
				Email:    "user0001@gmail.com",
			},
			StatusCode: http.StatusOK,
			Clean:      true,
		}, {
			Description: "New username, but existed email",
			Account: module.Account{
				Username: "user0002",
				Password: "user0002",
				Email:    "user0001@gmail.com",
			},
			StatusCode: http.StatusInternalServerError,
			Clean:      false,
		}, {
			Description: "New email, but existed username",
			Account: module.Account{
				Username: "user0001",
				Password: "user0003",
				Email:    "user0003@gmail.com",
			},
			StatusCode: http.StatusInternalServerError,
			Clean:      false,
		}, {
			Description: "New username and email",
			Account: module.Account{
				Username: "user0004",
				Password: "user0004",
				Email:    "user0004@gmail.com",
			},
			StatusCode: http.StatusOK,
			Clean:      true,
		},
	}

	defer func() {
		for _, test := range tests {
			if test.Clean {
				_ = module.DeleteAccount(context.Background(), d, test.Account)
			}
		}
	}()

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			b, err := json.Marshal(&test.Account)
			asserts.NoError(err)

			rsp, err := http.Post(ts.URL+"/account", "application/json", bytes.NewBuffer(b))
			asserts.NoError(err)
			defer rsp.Body.Close()

			asserts.Equal(test.StatusCode, rsp.StatusCode)
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	ts := test.NewTestRouter(gin.Default(), AccountService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()

	asserts := assert.New(t)

	users := []module.Account{
		{
			Username: "user0001",
			Password: "user0001",
			Email:    "user0001@gmail.com",
			Phone:    null.StringFrom("12345"),
		},
	}

	for _, user := range users {
		_ = module.CreateAccount(context.Background(), d, user)
	}

	tests := []struct {
		Description string
		Account     module.Account
		StatusCode  int
	}{
		{
			Description: "Delete matched username and email pair",
			Account: module.Account{
				Username: "user0001",
				Email:    "user0001@gmail.com",
			},
			StatusCode: http.StatusOK,
		}, {
			Description: "Delete user already be deleted",
			Account: module.Account{
				Username: "user0001",
				Email:    "user0001@gmail.com",
			},
			StatusCode: http.StatusNotFound,
		}, {
			Description: "Delete user which never registered",
			Account: module.Account{
				Username: "user0099",
				Email:    "user0099@gmail.com",
			},
			StatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			b, err := json.Marshal(&test.Account)
			asserts.NoError(err)

			req, err := http.NewRequest("DELETE", ts.URL+"/account/"+test.Account.Username, bytes.NewBuffer(b))
			asserts.NoError(err)

			rsp, err := http.DefaultClient.Do(req)
			asserts.NoError(err)
			defer rsp.Body.Close()

			asserts.Equal(test.StatusCode, rsp.StatusCode)
		})
	}
}

func TestUpdateAndGetAccountInfo(t *testing.T) {
	ts := test.NewTestRouter(gin.Default(), AccountService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()

	asserts := assert.New(t)

	users := []module.Account{
		{
			Username: "user0001",
			Email:    "user0001@gmail.com",
			Phone:    null.StringFrom("0912345678"),
		},
	}

	for _, user := range users {
		_ = module.CreateAccount(context.Background(), d, user)
	}

	defer func() {
		for _, user := range users {
			_ = module.DeleteAccount(context.Background(), d, user)
		}
	}()

	tests := []struct {
		Description      string
		Account          module.Account
		UpdateStatusCode int
		GetStatusCode    int
	}{
		{
			Description: "Update exised account information and get it back to verify",
			Account: module.Account{
				Username: "user0001",
				Email:    "user0001@gmail.com",
				Phone:    null.StringFrom("0909111222"),
			},
			UpdateStatusCode: http.StatusOK,
			GetStatusCode:    http.StatusOK,
		}, {
			Description: "Update an non-exist account information",
			Account: module.Account{
				Username: "user0004",
			},
			UpdateStatusCode: http.StatusNotFound,
			GetStatusCode:    http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			b, err := json.Marshal(&test.Account)
			asserts.NoError(err)

			ureq, err := http.NewRequest("PUT", ts.URL+"/account/"+test.Account.Username, bytes.NewBuffer(b))
			asserts.NoError(err)

			ursp, err := http.DefaultClient.Do(ureq)
			asserts.NoError(err)
			defer ursp.Body.Close()

			asserts.Equal(test.UpdateStatusCode, ursp.StatusCode)

			greq, err := http.NewRequest("GET", ts.URL+"/account/"+test.Account.Username, bytes.NewBuffer(b))
			asserts.NoError(err)

			grsp, err := http.DefaultClient.Do(greq)
			asserts.NoError(err)
			defer grsp.Body.Close()

			asserts.Equal(test.GetStatusCode, grsp.StatusCode)

			if test.GetStatusCode == http.StatusNotFound {
				return
			}

			acc := module.Account{}
			body, err := ioutil.ReadAll(grsp.Body)
			asserts.NoError(err)

			err = json.Unmarshal(body, &acc)
			asserts.NoError(err)

			asserts.Equal(test.Account, acc)
		})
	}
}
