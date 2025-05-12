package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/00mark0/macva-press/utils"

	"context"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword("123456")
	arg := CreateUserParams{
		Username: utils.RandomUser(),
		Email:    utils.RandomEmail(),
		Password: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)
	require.False(t, user.Banned.Bool)
	require.False(t, user.IsDeleted.Bool)
	require.NotEmpty(t, user.Pfp)
	require.Equal(t, user.Role, "user")

	return user
}

func createAdminUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Println(err)
	}

	arg := CreateUserAdminParams{
		Username: os.Getenv("ADMIN_USERNAME"),
		Email:    os.Getenv("EMAIL"),
		Password: hashedPassword,
		Role:     "admin",
	}

	user, err := testQueries.CreateUserAdmin(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)
	require.False(t, user.Banned.Bool)
	require.False(t, user.IsDeleted.Bool)
	require.NotEmpty(t, user.Pfp)
	require.Equal(t, user.Role, "admin")

	return user
}

func TestCreateUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}
}

func TestCreateUserAdmin(t *testing.T) {
	createAdminUser(t)
}

func TestGetUserByID(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByID(context.Background(), user1.UserID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.NotZero(t, user2.UserID)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.NotZero(t, user2.UserID)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)
	arg := UpdateUserParams{
		UserID:   user1.UserID,
		Username: utils.RandomUser(),
		Pfp:      utils.RandomString(20),
	}

	err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
}

func TestUpdateUserPassword(t *testing.T) {
	user1 := createRandomUser(t)
	arg := UpdateUserPasswordParams{
		UserID:   user1.UserID,
		Password: utils.RandomString(6),
	}

	err := testQueries.UpdateUserPassword(context.Background(), arg)

	require.NoError(t, err)
}

func TestBanUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.BanUser(context.Background(), user.UserID)

	require.NoError(t, err)
}

func TestUnbanUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.BanUser(context.Background(), user.UserID)

	require.NoError(t, err)

	err = testQueries.UnbanUser(context.Background(), user.UserID)

	require.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.UserID)

	require.NoError(t, err)
}

func TestCheckEmailExists(t *testing.T) {
	user := createRandomUser(t)

	exists, err := testQueries.CheckEmailExists(context.Background(), user.Email)

	require.NoError(t, err)
	require.Equal(t, exists, int32(1))
}

func TestGetActiveUsersCount(t *testing.T) {
	count, err := testQueries.GetActiveUsersCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)

	_ = createRandomUser(t)

	count2, err := testQueries.GetActiveUsersCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, count2, count+1)
}

func TestGetActiveUsers(t *testing.T) {
	users, err := testQueries.GetActiveUsers(context.Background(), 5)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, user := range users {
		require.False(t, user.Banned.Bool)
		require.False(t, user.IsDeleted.Bool)
	}

	require.Len(t, users, 5)
}

func TestGetBannedUsersCount(t *testing.T) {
	count, err := testQueries.GetBannedUsersCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)

	user := createRandomUser(t)

	err = testQueries.BanUser(context.Background(), user.UserID)
	require.NoError(t, err)

	count2, err := testQueries.GetBannedUsersCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, count2, count+1)
}

func TestGetBannedUsers(t *testing.T) {
	users, err := testQueries.GetBannedUsers(context.Background(), 5)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, user := range users {
		require.True(t, user.Banned.Bool)
	}

	require.LessOrEqual(t, len(users), 5)
}

func TestGetDeletedUsersCount(t *testing.T) {
	count, err := testQueries.GetDeletedUsersCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)

	user := createRandomUser(t)

	err = testQueries.DeleteUser(context.Background(), user.UserID)
	require.NoError(t, err)

	count2, err := testQueries.GetDeletedUsersCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, count2, count+1)
}

func TestGetDeletedUsers(t *testing.T) {
	users, err := testQueries.GetDeletedUsers(context.Background(), 5)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, user := range users {
		require.True(t, user.IsDeleted.Bool)
	}

	require.LessOrEqual(t, len(users), 5)
}

func createRandomUserInteractive(username string, email string) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	if err != nil {
		log.Println(err)
	}

	arg := CreateUserParams{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	return user
}

func TestSearchActiveUsers(t *testing.T) {
	searchTerm := "test_user_" + utils.RandomString(5)
	_ = createRandomUserInteractive(searchTerm, utils.RandomEmail())

	arg := SearchActiveUsersParams{
		Limit:      50,
		SearchTerm: searchTerm[:8],
	}

	users, err := testQueries.SearchActiveUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.LessOrEqual(t, len(users), 50)

	for _, user := range users {
		require.False(t, user.Banned.Bool)
		require.False(t, user.IsDeleted.Bool)
		require.Equal(t, searchTerm[:8], user.Username[:8])
	}
}

func TestSearchDeletedUsers(t *testing.T) {
	searchTerm := "test_user_" + utils.RandomString(5)
	user := createRandomUserInteractive(searchTerm, utils.RandomEmail())

	// Mark the user as deleted
	err := testQueries.DeleteUser(context.Background(), user.UserID)
	require.NoError(t, err)

	// First search: by username
	arg := SearchDeletedUsersParams{
		Limit:      50,
		SearchTerm: searchTerm[:8],
	}

	users, err := testQueries.SearchDeletedUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.LessOrEqual(t, len(users), 50)

	// Check that the username matches the search term
	for _, user := range users {
		require.True(t, user.IsDeleted.Bool)
		require.Equal(t, searchTerm[:8], user.Username[:8])
	}

	// Second search: by deleted email
	deletedEmailSearchTerm := fmt.Sprintf("deleted_%v@", user.UserID) // Construct search term for the email
	arg = SearchDeletedUsersParams{
		Limit:      50,
		SearchTerm: deletedEmailSearchTerm,
	}

	users, err = testQueries.SearchDeletedUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.LessOrEqual(t, len(users), 50)

	// Check that the email matches the deleted format
	for _, user := range users {
		require.True(t, user.IsDeleted.Bool)
		expectedEmail := fmt.Sprintf("deleted_%v@example.com", user.UserID)
		require.Equal(t, expectedEmail, user.Email)
	}
}

func TestSearchBannedUsers(t *testing.T) {
	searchTerm := "test_user_" + utils.RandomString(5)
	user := createRandomUserInteractive(searchTerm, utils.RandomEmail())

	// Mark the user as banned
	err := testQueries.BanUser(context.Background(), user.UserID)
	require.NoError(t, err)

	// First search: by username
	arg := SearchBannedUsersParams{
		Limit:      50,
		SearchTerm: searchTerm[:8],
	}

	users, err := testQueries.SearchBannedUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.LessOrEqual(t, len(users), 50)

	// Check that the username matches the search term
	for _, user := range users {
		require.True(t, user.Banned.Bool)
		require.Equal(t, searchTerm[:8], user.Username[:8])
	}
}

func TestSetEmailVerified(t *testing.T) {
	user := createRandomUser(t)
	err := testQueries.SetEmailVerified(context.Background(), user.UserID)
	require.NoError(t, err)
}
