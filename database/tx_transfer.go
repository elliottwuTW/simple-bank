package database

import (
	"context"

	"github.com/simple_bank/model"
)

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    model.Transfer `json:"transfer"`
	FromAccount model.Account  `json:"from_account"`
	ToAccount   model.Account  `json:"to_account"`
	FromEntry   model.Entry    `json:"from_entry"`
	ToEntry     model.Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (d *Database) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// 因為 transaction 沒辦法回傳結果，可以把變數放在外面，透過 closure 把結果讀回來
	err := d.execTx(ctx, func(ctx context.Context, d *Database) error {
		var err error

		// 建立 transfer
		// result.Transfer, err = d.accountDAO.CreateTransfer(ctx, CreateTransferParams{
		// 	FromAccountID: arg.FromAccountID,
		// 	ToAccountID:   arg.ToAccountID,
		// 	Amount:        arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// 建立 from entry
		// result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
		// 	AccountID: arg.FromAccountID,
		// 	Amount:    -arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// 建立 to entry
		// result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
		// 	AccountID: arg.ToAccountID,
		// 	Amount:    arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		return err
	})

	return result, err
}
