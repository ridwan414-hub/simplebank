package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}



// txKey is the context key for the transaction name
var txKey = struct{}{}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer Transfer `json:"transfer"` // the created transfer
	FromAccount Account `json:"from_account"` // the updated from account
	ToAccount Account `json:"to_account"` // the updated to account
	FromEntry Entry `json:"from_entry"` // the created entry for the from account
	ToEntry Entry `json:"to_entry"` // the created entry for the to account
}

// NewStore creates a new Store
func NewStore(db *sql.DB)*Store{
	return &Store{
		Queries: New(db),
		db: db,
	}
}

// ExecTx executes a function within a database transaction
func(store *Store) execTx(ctx context.Context,fn func(*Queries)error)error{
	tx,err := store.db.BeginTx(ctx,nil)

	if err != nil{
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil{
		if rbErr := tx.Rollback(); rbErr != nil{
			return fmt.Errorf("tx error: %v, rb error: %v",err,rbErr)
		}
		return err
	}

	return tx.Commit()
}

/*TransferTx performs a money transfer from one account to the other.
It creates a transfer record, and update account balance within the same transaction.*/
func(store *Store)TransferTx(ctx context.Context,arg TransferTxParams)(TransferTxResult,error){
	var result TransferTxResult

	err := store.execTx(ctx,func(q *Queries)error{
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName,"create transfer")
		result.Transfer,err = q.CreateTransfer(ctx,CreateTransferParams(arg)) // create transfer record

		if err != nil{
			return err
		}

		fmt.Println(txName,"create entry 1")
		result.FromEntry,err = q.CreateEntry(ctx,CreateEntryParams{ // create entry for the from account
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}

		fmt.Println(txName,"create entry 2")
		result.ToEntry,err = q.CreateEntry(ctx,CreateEntryParams{ // create entry for the to account
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		fmt.Println(txName,"get account 1")
		account1,err := q.GetAccountForUpdate(ctx,arg.FromAccountID)
		if err != nil{
			return err
		}

		fmt.Println(txName,"update account 1")
		result.FromAccount,err = q.UpdateAccount(ctx,UpdateAccountParams{// update the from account balance
			ID: arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil{
			return err
		}

		fmt.Println(txName,"get account 2")
		account2,err := q.GetAccountForUpdate(ctx,arg.ToAccountID)
		if err != nil{
			return err
		}

		fmt.Println(txName,"update account 2")
		result.ToAccount,err = q.UpdateAccount(ctx,UpdateAccountParams{// update the to account balance
			ID: arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil{
			return err
		}

		return nil
	})

	return result,err
}