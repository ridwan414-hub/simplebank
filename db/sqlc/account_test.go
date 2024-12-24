package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ridwan414-hub/simplebank/utils"
	"github.com/stretchr/testify/require"
)

// CreateAccountParams defines the input fields for the CreateAccount function
func createRandomAccount(t *testing.T)Account{
		arg :=CreateAccountParams{
		Owner: utils.RandomOwner(),
		Balance : utils.RandomMoney(),
		Currency: utils.RandomCurrrency(),
	}

	account,err := testQueries.CreateAccount(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,account)

	require.Equal(t,arg.Owner,account.Owner)
	require.Equal(t,arg.Balance,account.Balance)
	require.Equal(t,arg.Currency,account.Currency)

	require.NotZero(t,account.ID)
	require.NotZero(t,account.CreatedAt)
	return account
}

// TestCreateAccount tests the CreateAccount function
func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

// TestDeleteAccount tests the DeleteAccount function
func TestGetAccount(t *testing.T){
	account1 := createRandomAccount(t)
	account2,err := testQueries.GetAccount(context.Background(),account1.ID)

	require.NoError(t,err)
	require.NotEmpty(t,account2)

	require.Equal(t,account1.Owner,account2.Owner)
	require.Equal(t,account1.Balance,account2.Balance)
	require.Equal(t,account1.Currency,account2.Currency)
	require.Equal(t,account1.ID,account2.ID)
	require.WithinDuration(t,account1.CreatedAt,account2.CreatedAt,time.Second)
}

// TestUpdateAccount tests the DeleteAccount function
func TestUpdateAccount(t *testing.T){
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2,err := testQueries.UpdateAccount(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,account2)

	require.Equal(t,account1.Owner,account2.Owner)
	require.Equal(t,arg.Balance,account2.Balance)
	require.Equal(t,account1.Currency,account2.Currency)
	require.Equal(t,account1.ID,account2.ID)
	require.WithinDuration(t,account1.CreatedAt,account2.CreatedAt,time.Second)
}

// TestDeleteAccount tests the DeleteAccount function
func TestDeleteAccount(t *testing.T){
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(),account1.ID)

	require.NoError(t,err)

	account2,err := testQueries.GetAccount(context.Background(),account1.ID)
	require.Error(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,account2)
}

// TestListAccount tests the ListAccount function
func TestListAccount(t *testing.T){
	for i:=0;i<10;i++{
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	account,err:=testQueries.ListAccounts(context.Background(),arg)
	require.NoError(t,err)
	require.Len(t,account,5)

	for _,account := range account{
		require.NotEmpty(t,account)
	}
}