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

		result.Transfer,err = q.CreateTransfer(ctx,CreateTransferParams(arg)) // create transfer record

		if err != nil{
			return err
		}

		result.FromEntry,err = q.CreateEntry(ctx,CreateEntryParams{ // create entry for the from account
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}

		result.ToEntry,err = q.CreateEntry(ctx,CreateEntryParams{ // create entry for the to account
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return result,err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	ammount1 int64,
	accountID2 int64,
	ammount2 int64,
)(account1 Account,account2 Account,err error){
	account1 , err =q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountID1,
		Amount: ammount1,
	})
	if err!=nil{
		return
	}

	account2 , err =q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountID2,
		Amount: ammount2,
	})
	return
}