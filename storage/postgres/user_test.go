package postgres_test

import (
	"testing"

	"github.com/SaidovZohid/hotel-project/storage/postgres"
	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) int64 {
	user_id, err := dbManager.User().Create(&repo.User{
		FirstName: faker.FirstName(),
		LastName: faker.LastName(),
		Email: faker.Email(),
		Password: faker.Password(),
		Type: postgres.UserType,
	})
	require.NoError(t, err)
	return user_id
}

func deleteUser(t *testing.T, user_id int64) {
	err := dbManager.User().Delete(user_id)
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	user_id := createUser(t)
	deleteUser(t, user_id)
}

func TestGetUser(t *testing.T) {
	user_id := createUser(t)
	user, err := dbManager.User().Get(user_id)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	deleteUser(t, user_id)
}

func TestUpdateUser(t *testing.T) {
	user_id := createUser(t)
	update_at, err := dbManager.User().Update(&repo.User{
		ID: user_id,
		FirstName: faker.FirstName(),
		LastName: faker.LastName(),
		Email: faker.Email(),
		Type: postgres.SuperUserType,
	})
	require.NoError(t, err)
	require.NotEmpty(t, update_at)
	deleteUser(t, user_id)
}

func TestDeleteUser(t *testing.T) {
	user_id := createUser(t)
	deleteUser(t, user_id)
}

func TestGetAllUsers(t *testing.T) {
	user_id := createUser(t)
	users, err := dbManager.User().GetAll(&repo.GetAllUsersParams{
		Limit: 10,
		Page: 1,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users.Users), 1)
	deleteUser(t, user_id)
}