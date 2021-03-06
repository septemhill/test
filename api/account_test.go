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
	"github.com/septemhill/test/utils"
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
	tk := test.NewTestSessionToken(r)

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
			Description: "New another username and email",
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
			_, _ = module.DeleteAccount(context.Background(), d, &test.Account)
		}
	}()

	header := map[string]string{
		utils.HEADER_SESSION_TOKEN: tk,
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			b, err := json.Marshal(&test.Account)
			asserts.NoError(err)

			req, err := NewRequestWithTestHeader("POST", ts.URL+"/account/", bytes.NewBuffer(b), header)
			asserts.NoError(err)

			rsp, err := http.DefaultClient.Do(req)
			asserts.NoError(err)
			defer rsp.Body.Close()

			body, err := ioutil.ReadAll(rsp.Body)
			asserts.NoError(err)

			asserts.Equal(test.StatusCode, rsp.StatusCode, string(body))
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), AccountService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()
	tk := test.NewTestSessionToken(r)

	asserts := assert.New(t)

	user1 := test.NewAccount(ctx, d, false)

	tests := []struct {
		Description string
		Account     module.Account
		StatusCode  int
	}{
		{
			Description: "Delete matched username and email pair",
			Account: module.Account{
				Username: user1.Username,
				Email:    user1.Email,
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

	header := map[string]string{
		"sessionToken": tk,
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			b, err := json.Marshal(&test.Account)
			asserts.NoError(err)

			req, err := NewRequestWithTestHeader("DELETE", ts.URL+"/account/"+test.Account.Username, bytes.NewBuffer(b), header)
			asserts.NoError(err)

			rsp, err := http.DefaultClient.Do(req)
			asserts.NoError(err)
			defer rsp.Body.Close()

			body, err := ioutil.ReadAll(rsp.Body)
			asserts.NoError(err)

			asserts.Equal(test.StatusCode, rsp.StatusCode, string(body))
		})
	}
}

func TestUpdateAndGetAccountInfo(t *testing.T) {
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), AccountService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()
	tk := test.NewTestSessionToken(r)

	asserts := assert.New(t)

	users := []*module.Account{test.NewAccount(ctx, d, false)}

	defer func() {
		test.DeleteAccounts(ctx, d, users...)
	}()

	users[0].Phone = null.StringFrom("0909123123")
	header := map[string]string{
		"sessionToken": tk,
	}

	tests := []struct {
		Description      string
		Account          module.Account
		UpdateStatusCode int
		GetStatusCode    int
	}{
		{
			Description:      "Update exised account information and get it back to verify",
			Account:          *users[0],
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

			ureq, err := NewRequestWithTestHeader("PUT", ts.URL+"/account/"+test.Account.Username, bytes.NewBuffer(b), header)
			asserts.NoError(err)

			ursp, err := http.DefaultClient.Do(ureq)
			asserts.NoError(err)
			defer ursp.Body.Close()

			ubody, err := ioutil.ReadAll(ursp.Body)
			asserts.NoError(err)

			asserts.Equal(test.UpdateStatusCode, ursp.StatusCode, string(ubody))

			greq, err := NewRequestWithTestHeader("GET", ts.URL+"/account/"+test.Account.Username, bytes.NewBuffer(b), header)
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

func TestChangePassword(t *testing.T) {
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), RootService, AccountService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()

	users := []*module.Account{
		test.NewAccount(ctx, d, true),
		test.NewAccount(ctx, d, true),
	}

	tests := []struct {
		Description              string
		Account                  module.Account
		Password                 password
		LoginBeforeChgStatusCode int
		LoginAfterChgStatusCode  int
		PasswordChgStatusCode    int
	}{
		{
			Description: "Login with existed user 1",
			Account:     *users[0],
			Password: password{
				Password: "thisisnewpassword",
			},
			LoginBeforeChgStatusCode: http.StatusOK,
			LoginAfterChgStatusCode:  http.StatusOK,
			PasswordChgStatusCode:    http.StatusOK,
		}, {
			Description: "Login with existed user 2",
			Account:     *users[1],
			Password: password{
				Password: "anothernewpassword",
			},
			LoginBeforeChgStatusCode: http.StatusOK,
			LoginAfterChgStatusCode:  http.StatusOK,
			PasswordChgStatusCode:    http.StatusOK,
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			// Login before change password
			lbb, err := json.Marshal(&test.Account)
			asserts.NoError(err)

			lbreq, err := http.NewRequest("POST", ts.URL+"/login", bytes.NewBuffer(lbb))
			asserts.NoError(err)

			lbrsp, err := http.DefaultClient.Do(lbreq)
			asserts.NoError(err)
			defer lbrsp.Body.Close()

			asserts.Equal(test.LoginBeforeChgStatusCode, lbrsp.StatusCode)

			m := make(map[string]string)
			lbbody, err := ioutil.ReadAll(lbrsp.Body)
			asserts.NoError(err)

			err = json.Unmarshal(lbbody, &m)
			asserts.NoError(err)

			// Change password
			chgb, err := json.Marshal(&test.Password)
			asserts.NoError(err)

			chgreq, err := http.NewRequest("PUT", ts.URL+"/account/"+test.Account.Username+"/chgpasswd", bytes.NewBuffer(chgb))
			asserts.NoError(err)

			code := m["code"]
			chgreq.Header.Add(utils.HEADER_SESSION_TOKEN, code)

			chgrsp, err := http.DefaultClient.Do(chgreq)
			asserts.NoError(err)
			defer chgrsp.Body.Close()

			chgbody, err := ioutil.ReadAll(chgrsp.Body)
			asserts.NoError(err)

			asserts.Equal(test.PasswordChgStatusCode, chgrsp.StatusCode, string(chgbody))

			// Login with new password
			test.Account.Password = test.Password.Password
			lab, err := json.Marshal(&test.Account)
			asserts.NoError(err)

			lareq, err := http.NewRequest("POST", ts.URL+"/login", bytes.NewBuffer(lab))
			asserts.NoError(err)

			larsp, err := http.DefaultClient.Do(lareq)
			asserts.NoError(err)
			defer larsp.Body.Close()

			asserts.Equal(test.LoginAfterChgStatusCode, larsp.StatusCode)
		})
	}
}
