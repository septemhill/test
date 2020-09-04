package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/module"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	ts := newTestRouter(gin.Default(), AccountService)
	assert := assert.New(t)

	tests := []struct {
		Description string
		Account     module.Account
		Err         error
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
			Err:        nil,
			StatusCode: http.StatusOK,
			Clean:      true,
		}, {
			Description: "New username, but existed email",
			Account: module.Account{
				Username: "user0002",
				Password: "user0002",
				Email:    "user0001@gmail.com",
			},
			Err:        nil,
			StatusCode: http.StatusInternalServerError,
			Clean:      false,
		}, {
			Description: "New email, but existed username",
			Account: module.Account{
				Username: "user0001",
				Password: "user0003",
				Email:    "user0003@gmail.com",
			},
			Err:        nil,
			StatusCode: http.StatusInternalServerError,
			Clean:      false,
		}, {
			Description: "New username and email",
			Account: module.Account{
				Username: "user0004",
				Password: "user0004",
				Email:    "user0004@gmail.com",
			},
			Err:        nil,
			StatusCode: http.StatusOK,
			Clean:      true,
		},
	}

	defer func() {
		for _, test := range tests {
			_ = test
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
