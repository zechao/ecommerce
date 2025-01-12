package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	pgContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/zechao158/ecomm/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func TestCRUDStore(t *testing.T) {
	ctx := context.Background()
	dbname := "yourdb"
	user := "youruser"
	password := "yourpassword"
	// 1. Start the postgres ctr and run any migrations on it
	ctr, err := pgContainer.Run(
		ctx,
		"postgres:16-alpine",
		pgContainer.WithDatabase(dbname),
		pgContainer.WithUsername(user),
		pgContainer.WithPassword(password),
		pgContainer.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, ctr)
	assert.NoError(t, err)

	dbURL, err := ctr.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	db = db.Debug()
	assert.NoError(t, err)

	err = db.AutoMigrate(&TestModel{})
	assert.NoError(t, err)

	RunTest := func(name string, f func(t *testing.T, tx *gorm.DB)) {
		tx := db.Begin()
		defer tx.Rollback()
		t.Run(name, func(t *testing.T) {
			f(t, tx)
		})
	}

	RunTest("Create", func(t *testing.T, tx *gorm.DB) {
		store := storage.New[TestModel](tx)
		model := TestModel{ID: uuid.New(), Name: "Test"}

		err := store.Create(ctx, &model)
		assert.NoError(t, err)
		var created TestModel
		err = tx.First(&created).Error
		assert.NoError(t, err)
		equalTestModel(t, &model, &created)
	})

	RunTest("Create already exist", func(t *testing.T, tx *gorm.DB) {
		store := storage.New[TestModel](tx)
		model := TestModel{ID: uuid.New(), Name: "Test"}

		err := store.Create(ctx, &model)
		assert.NoError(t, err)
		var created TestModel
		err = tx.First(&created).Error
		assert.NoError(t, err)
		equalTestModel(t, &model, &created)

		err = store.Create(ctx, &model)
		assert.ErrorIs(t, err, storage.ErrDuplicateKey)
	})

	RunTest("Delete", func(t *testing.T, tx *gorm.DB) {
		store := storage.New[TestModel](tx)
		model := TestModel{ID: uuid.New(), Name: "Test"}
		assert.NoError(t, tx.Create(&model).Error)
		err := store.Delete(ctx, &model)
		assert.NoError(t, err)
		var count int64
		tx.Model(&TestModel{}).Where("id = ?", model.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	RunTest("GetByID", func(t *testing.T, tx *gorm.DB) {
		store := storage.New[TestModel](tx)
		model := TestModel{ID: uuid.New(), Name: "Test"}
		assert.NoError(t, tx.Create(&model).Error)
		result, err := store.GetByID(ctx, model.ID, true) // Using select for update
		assert.NoError(t, err)
		equalTestModel(t, &model, result)
	})

	RunTest("Update", func(t *testing.T, tx *gorm.DB) {
		store := storage.New[TestModel](tx)
		model := TestModel{ID: uuid.New(), Name: "Test"}
		assert.NoError(t, tx.Create(&model).Error)
		model.Name = "Updated"
		err := store.Update(ctx, &model)
		assert.NoError(t, err)
		var result TestModel
		assert.NoError(t, tx.First(&result, "id = ?", model.ID).Error)
		assert.Equal(t, "Updated", result.Name)
	})

	RunTest("GetAll", func(t *testing.T, tx *gorm.DB) {
		store := storage.New[TestModel](tx)
		model1 := &TestModel{ID: uuid.New(), Name: "Test1"}
		model2 := &TestModel{ID: uuid.New(), Name: "Test2"}
		assert.NoError(t, tx.Create(model1).Error)
		assert.NoError(t, tx.Create(model2).Error)

		results, err := store.GetAll(ctx, nil)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})

}

func equalTestModel(t *testing.T, expected, actual *TestModel) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Microsecond)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, time.Microsecond)
}
