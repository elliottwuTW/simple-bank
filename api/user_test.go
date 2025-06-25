package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/simple_bank/database"
	mockdb "github.com/simple_bank/database/mock"
	"github.com/simple_bank/model"
	"github.com/simple_bank/util"
	"github.com/stretchr/testify/require"
)

// 客製化自己的 matcher
type eqCreateUserParams struct {
	params database.CreateUserParams
	// 原始密碼值
	password string
}

func (e eqCreateUserParams) Matches(x interface{}) bool {
	// x 就是跑測試時，runtime 傳入的參數
	inputParams, ok := x.(database.CreateUserParams)
	if !ok {
		return false
	}

	// 因為比對參數時，只有密碼比較特殊，才特別拉出來比對
	err := util.CheckPassword(e.password, inputParams.HashedPassword)
	if err != nil {
		return false
	}

	// 把密碼 hash 值覆蓋掉，比對剩下其他的欄位值
	e.params.HashedPassword = inputParams.HashedPassword
	return reflect.DeepEqual(e.params, inputParams)
}

func (e eqCreateUserParams) String() string {
	return fmt.Sprintf("matches params %+v and password %+v", e.params, e.password)
}

func EqCreateUserParams(arg database.CreateUserParams, pw string) gomock.Matcher {
	return eqCreateUserParams{params: arg, password: pw}
}

func TestCreateUserAPI(t *testing.T) {
	pwd := util.RandomString(6)
	hashedPwd, err := util.HashPassword(pwd)
	require.NoError(t, err)

	user := model.User{
		Username:       "user",
		HashedPassword: hashedPwd,
		Email:          "user@gmail.com",
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(db *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "user created",
			body: gin.H{
				"username": user.Username,
				"password": pwd,
				"email":    user.Email,
			},
			buildStubs: func(db *mockdb.MockDatabase) {
				arg := database.CreateUserParams{
					Username:       user.Username,
					HashedPassword: user.HashedPassword,
					Email:          user.Email,
				}
				db.EXPECT().
					// CreateUser(gomock.Any(), gomock.Eq(arg)). /// hashedPassword 每次都不一樣
					CreateUser(gomock.Any(), EqCreateUserParams(arg, pwd)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				// 透過 io.ReadAll 讀回 recorder.Body 的資料
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var createdUser model.User
				err = json.Unmarshal(data, &createdUser)
				require.NoError(t, err)

				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},
	}

	for _, tc := range testCases {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mockdb.NewMockDatabase(ctrl)
		tc.buildStubs(db)

		// Start test server and send request
		server := newTestServer(t, db)
		// We don't need to start a REAL HTTP server.
		// Instead, we can just use the Recorder feature
		// of the httptest package
		recorder := httptest.NewRecorder()

		// 產生 request 來送給 HTTP server
		path := "/users"
		marshalled, err := json.Marshal(tc.body)
		require.NoError(t, err)
		request, err := http.NewRequest(
			http.MethodPost, path, bytes.NewReader(marshalled),
		)
		require.NoError(t, err)

		// use recorder to record the response of the API request
		// send to the HTTP server
		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	}
}
