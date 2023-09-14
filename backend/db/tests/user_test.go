package db_test

import (
	"context"
	"log"
	db "nextjs/backend/db/sqlc"
	"nextjs/backend/utils"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func clean_up() {
	err := testQuery.DeleteAllUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func createRandomUser(t *testing.T) db.User {
	hashedpw, err := utils.GenerateHashPassword(utils.RandomString(8))
	if err != nil {
		log.Fatal("Unable to generate hashed password.", err)
	}

	arg := db.CreateUserParams{
		Email:          utils.RandomEmail(),
		HashedPassword: hashedpw,
	}

	user, err := testQuery.CreateUser(context.Background(), arg)
	// Test user creation
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, user.Email, arg.Email)
	assert.Equal(t, user.HashedPassword, arg.HashedPassword)
	assert.WithinDuration(t, user.CreatedAt, time.Now(), 10*time.Second)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 10*time.Second)

	return user
}

func TestCreateUser(t *testing.T) {
	defer clean_up()
	user1 := createRandomUser(t)
	// To ensure email provided is unique
	user2, err := testQuery.CreateUser(context.Background(), db.CreateUserParams{
		Email:          user1.Email,
		HashedPassword: user1.HashedPassword,
	})
	assert.Error(t, err)
	assert.Empty(t, user2)

}

func TestUpdateUser(t *testing.T) {
	defer clean_up()
	user := createRandomUser(t)

	newPassword, err := utils.GenerateHashPassword(utils.RandomString(8))
	if err != nil {
		log.Fatal("Unable to generate hash password", err)
	}

	arg := db.UpdateUserPasswordParams{
		HashedPassword: newPassword,
		ID:             user.ID,
		UpdatedAt:      time.Now(),
	}

	newUser, err := testQuery.UpdateUserPassword(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, newUser)
	assert.Equal(t, newUser.HashedPassword, arg.HashedPassword)
	assert.Equal(t, user.Email, newUser.Email)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 10*time.Second)
}

func TestListUsers(t *testing.T) {
	defer clean_up()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			createRandomUser(t)
			defer wg.Done()
		}()
	}

	wg.Wait()

	arg := db.ListUsersParams{
		Offset: 0,
		Limit:  10,
	}

	users, err := testQuery.ListUsers(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, len(users), 10)
}
