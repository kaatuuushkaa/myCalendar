package user_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"myCalendar/internal/user"
	"os"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	host := getEnv("TEST_DB_TEST", "localhost")
	port := getEnv("TEST_DB_PORT", "5434")
	userDB := getEnv("TEST_DB_USER", "postgres")
	password := getEnv("TEST_DB_PASSWORD", "yourpassword")
	dbname := getEnv("TEST_DB_DBNAME", "postgres_test")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, userDB, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "не удалось подключиться к тестовой БД")

	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	return db
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func TestRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	repo := user.NewRepository(db)
	u := &user.User{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "alice",
		Password: "hashpassword",
		Email:    "alice@test.com",
		Name:     "Alice",
		Surname:  "Smith",
		Birth:    "2000-01-15",
	}

	err := repo.Create(context.Background(), u)
	assert.NoError(t, err)

	var count int64
	db.Model(&user.User{}).Where("username = ?", "alice").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestRepository_GetByUsername(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	repo := user.NewRepository(db)

	u := &user.User{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "alice",
		Password: "hashpassword",
		Email:    "alice@test.com",
		Name:     "Alice",
		Surname:  "Smith",
		Birth:    "2000-01-15",
	}
	require.NoError(t, repo.Create(context.Background(), u))

	found, err := repo.GetByUsername(context.Background(), "alice")
	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", found.ID)
	assert.Equal(t, "Alice", found.Name)
}

func TestRepository_GetByUsername_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	repo := user.NewRepository(db)

	found, err := repo.GetByUsername(context.Background(), "ghost")
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestRepository_GetByLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	repo := user.NewRepository(db)

	u := &user.User{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "alice",
		Password: "hashpassword",
		Email:    "alice@test.com",
		Name:     "Alice",
		Surname:  "Smith",
		Birth:    "2000-01-15",
	}
	require.NoError(t, repo.Create(context.Background(), u))

	byUsername, err := repo.GetByLogin(context.Background(), "alice")
	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", byUsername.ID)

	byEmail, err := repo.GetByLogin(context.Background(), "alice@test.com")
	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", byEmail.ID)
}

func TestRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	repo := user.NewRepository(db)

	u := &user.User{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "alice",
		Password: "hashpassword",
		Email:    "alice@test.com",
		Name:     "Alice",
		Surname:  "Smith",
		Birth:    "2000-01-15",
	}
	require.NoError(t, repo.Create(context.Background(), u))

	updated, err := repo.Update(context.Background(), "alice", "alice@test.com", "Alicia", "Kun", "2000-01-15")
	assert.NoError(t, err)
	assert.Equal(t, "Alicia", updated.Name)
	assert.Equal(t, "Kun", updated.Surname)

	fromDB, _ := repo.GetByUsername(context.Background(), "alice")
	assert.Equal(t, "Alicia", fromDB.Name)
}

func TestRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	repo := user.NewRepository(db)

	u := &user.User{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "alice",
		Password: "hashpassword",
		Email:    "alice@test.com",
		Name:     "Alice",
		Surname:  "Smith",
		Birth:    "2000-01-15",
	}
	require.NoError(t, repo.Create(context.Background(), u))

	err := repo.Delete(context.Background(), "alice")
	assert.NoError(t, err)

	found, err := repo.GetByUsername(context.Background(), "alice")
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestRepository_DuplicateUsername(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	repo := user.NewRepository(db)

	u1 := &user.User{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "alice",
		Password: "hashpassword",
		Email:    "alice@test.com",
		Name:     "Alice",
		Surname:  "Smith",
		Birth:    "2000-01-15",
	}
	u2 := &user.User{
		ID:       "550e8400-e29b-41d4-a716-446655440001",
		Username: "alice",
		Password: "hashpassword",
		Email:    "alice@test.com",
		Name:     "Alice",
		Surname:  "Smith",
		Birth:    "2000-01-15",
	}
	require.NoError(t, repo.Create(context.Background(), u1))

	err := repo.Create(context.Background(), u2)
	assert.Error(t, err)
}
