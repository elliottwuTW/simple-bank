package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/simple_bank/database"
	mockdb "github.com/simple_bank/database/mock"
	"github.com/simple_bank/model"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {
	account := model.Account{
		Owner:    "owner",
		Currency: "USD",
		Balance:  0,
	}

	testCases := []struct {
		name          string
		buildStubs    func(db *mockdb.MockDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "account created",
			buildStubs: func(db *mockdb.MockDatabase) {
				arg := database.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				db.EXPECT().
					CreateAccount(gomock.Any(), arg).
					Times(1).
					Return(model.Account{Owner: account.Owner, Currency: account.Currency}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				// 透過 io.ReadAll 讀回 recorder.Body 的資料
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				var createdAccount model.Account
				err = json.Unmarshal(data, &createdAccount)
				require.NoError(t, err)

				require.Equal(t, account.Owner, createdAccount.Owner)
				require.Equal(t, account.Currency, createdAccount.Currency)
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
		path := "/accounts"
		body := CreateAccountReq{Owner: account.Owner, Currency: account.Currency}
		marshalled, err := json.Marshal(body)
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
