package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5
	// amount := int64(10)

	// channel is designed to connect concurrent goroutines.
	errs := make(chan error)
	results := make(chan interface{})

	// 為了測試多個 operation，利用 concurrent goroutine 打多次操作
	// 連續 5 次轉帳 10 元
	for i := 0; i < n; i++ {
		go func() {
			result, err := testDB.accountDAO.CreateAccount(context.Background(), CreateAccountParams{})

			// 用 channel 一個一個檢查結果
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// // check transfer
		// transfer := result.Transfer
		// require.NotEmpty(t, transfer)
		// require.Equal(t, account1.ID, transfer.FromAccountID)
		// require.Equal(t, account2.ID, transfer.ToAccountID)
		// require.Equal(t, amount, transfer.Amount)
		// require.NotZero(t, transfer.ID)
		// require.NotZero(t, transfer.CreatedAt)

		// _, err = testStore.GetTransfer(context.Background(), transfer.ID)
		// require.NoError(t, err)

		// // check entries
		// fromEntry := result.FromEntry
		// require.NotEmpty(t, fromEntry)
		// require.Equal(t, account1.ID, fromEntry.AccountID)
		// require.Equal(t, -amount, fromEntry.Amount)
		// require.NotZero(t, fromEntry.ID)
		// require.NotZero(t, fromEntry.CreatedAt)

		// _, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		// require.NoError(t, err)

		// toEntry := result.ToEntry
		// require.NotEmpty(t, toEntry)
		// require.Equal(t, account2.ID, toEntry.AccountID)
		// require.Equal(t, amount, toEntry.Amount)
		// require.NotZero(t, toEntry.ID)
		// require.NotZero(t, toEntry.CreatedAt)

		// _, err = testStore.GetEntry(context.Background(), toEntry.ID)
		// require.NoError(t, err)

		// // check accounts
		// fromAccount := result.FromAccount
		// require.NotEmpty(t, fromAccount)
		// require.Equal(t, account1.ID, fromAccount.ID)

		// toAccount := result.ToAccount
		// require.NotEmpty(t, toAccount)
		// require.Equal(t, account2.ID, toAccount.ID)
	}

	// require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	// require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
