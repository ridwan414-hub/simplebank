package db

import (
	"context"
	"testing"
	"time"

	"github.com/ridwan414-hub/simplebank/utils"
	"github.com/stretchr/testify/require"
)

// CreateAccountParams defines the input fields for the CreateAccount function
func createRandomUser(t *testing.T)User{
	hashedPassword,err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t,err)
	
	arg :=CreateUserParams{
		Username: utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName: utils.RandomOwner(),
		Email: utils.RandomEmail(),
	}

	user,err := testQueries.CreateUser(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,user)

	require.Equal(t,arg.Username,user.Username)
	require.Equal(t,arg.HashedPassword,user.HashedPassword)
	require.Equal(t,arg.FullName,user.FullName)
	require.Equal(t,arg.Email,user.Email)

	require.True(t,user.PasswordChangedAt.IsZero())
	require.NotZero(t,user.CreatedAt)
	return user
}

// TestCreateAccount tests the CreateAccount function
func TestCreateUser(t *testing.T){
	createRandomUser(t)
}

// TestDeleteAccount tests the DeleteAccount function
func TestGetUser(t *testing.T){
	user1 := createRandomUser(t)
	user2,err := testQueries.GetUser(context.Background(),user1.Username)

	require.NoError(t,err)
	require.NotEmpty(t,user2)

	require.Equal(t,user1.Username,user2.Username)
	require.Equal(t,user1.HashedPassword,user2.HashedPassword)
	require.Equal(t,user1.FullName,user2.FullName)
	require.Equal(t,user1.Email,user2.Email)
	
	require.WithinDuration(t,user1.PasswordChangedAt,user2.PasswordChangedAt,time.Second)
	require.WithinDuration(t,user1.CreatedAt,user2.CreatedAt,time.Second)
}