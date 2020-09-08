package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
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

	assert := assert.New(t)

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
				module.DeleteAccount(context.Background(), d, test.Account)
			}
		}
	}()

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			b, err := json.Marshal(&test.Account)
			assert.NoError(err)

			rsp, err := http.Post(ts.URL+"/account", "application/json", bytes.NewBuffer(b))
			assert.NoError(err)
			defer rsp.Body.Close()

			assert.Equal(test.StatusCode, rsp.StatusCode)
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

	assert := assert.New(t)

	users := []module.Account{
		{
			Username: "user0001",
			Password: "user0001",
			Email:    "user0001@gmail.com",
			Phone:    null.StringFrom("12345"),
		},
	}

	for _, user := range users {
		module.CreateAccount(context.Background(), d, user)
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
			assert.NoError(err)

			req, err := http.NewRequest("DELETE", ts.URL+"/account/"+test.Account.Username, bytes.NewBuffer(b))
			assert.NoError(err)

			rsp, err := http.DefaultClient.Do(req)
			assert.NoError(err)
			defer rsp.Body.Close()

			assert.Equal(test.StatusCode, rsp.StatusCode)
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

	assert := assert.New(t)

	users := []module.Account{
		{
			Username: "user0001",
			Email:    "user0001@gmail.com",
			Phone:    null.StringFrom("0912345678"),
		},
	}

	for _, user := range users {
		module.CreateAccount(context.Background(), d, user)
	}

	defer func() {
		for _, user := range users {
			module.DeleteAccount(context.Background(), d, user)
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
			assert.NoError(err)

			ureq, err := http.NewRequest("PUT", ts.URL+"/account/"+test.Account.Username, bytes.NewBuffer(b))
			assert.NoError(err)

			ursp, err := http.DefaultClient.Do(ureq)
			assert.NoError(err)
			defer ursp.Body.Close()

			assert.Equal(test.UpdateStatusCode, ursp.StatusCode)

			greq, err := http.NewRequest("GET", ts.URL+"/account/"+test.Account.Username, bytes.NewBuffer(b))
			assert.NoError(err)

			grsp, err := http.DefaultClient.Do(greq)
			assert.NoError(err)
			defer grsp.Body.Close()

			assert.Equal(test.GetStatusCode, grsp.StatusCode)

			if test.GetStatusCode == http.StatusNotFound {
				return
			}

			acc := module.Account{}
			body, err := ioutil.ReadAll(grsp.Body)
			assert.NoError(err)

			err = json.Unmarshal(body, &acc)
			assert.NoError(err)

			assert.Equal(test.Account, acc)
		})
	}
}
