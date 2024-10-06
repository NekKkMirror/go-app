package tests

import (
	"context"
	"testing"
	"time"

	postgrescontainer "github.com/NekKkMirror/go-app/internal/pkg/container/test/postgres"
	ormpgsql "github.com/NekKkMirror/go-app/internal/pkg/orm-pgsql"
	"gorm.io/gorm"
)

// Entity for testing
type TestEntity struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func TestGenericRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := ormpgsql.NewGenericRepository[TestEntity](db)

	ctx := context.Background()

	// Test Create
	entity := TestEntity{ID: "1", Name: "Test Entity"}
	if err := repo.Create(ctx, &entity); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test GetById
	result, err := repo.GetById(ctx, "1")
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	if result.Name != "Test Entity" {
		t.Fatalf("expected entity name to be 'Test Entity', got %s", result.Name)
	}

	// Test Update
	entity.Name = "Updated Name"
	if err := repo.Update(ctx, &entity); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Test Get
	result, err = repo.Get(ctx, &TestEntity{Name: "Updated Name"})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if result.Name != "Updated Name" {
		t.Fatalf("expected entity name to be 'Updated Name', got %s", result.Name)
	}

	// Test Delete
	if err := repo.Delete(ctx, "1"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Test soft deletion
	_, err = repo.GetById(ctx, "1")
	if err == nil {
		t.Fatalf("expected error when getting deleted entity, got nil")
	}

	// Test CreateMany
	entities := []TestEntity{
		{ID: "2", Name: "Entity 2"},
		{ID: "3", Name: "Entity 3"},
	}
	if err := repo.CreateMany(ctx, &entities); err != nil {
		t.Fatalf("CreateMany failed: %v", err)
	}

	// Test GetAll
	all, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(*all) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(*all))
	}

	// Test Count
	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count to be 2, got %d", count)
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()
	DB, _, err := postgrescontainer.Start(ctx, t)
	if err != nil {
		t.Fatalf("failed to start PostgreSQL container: %v", err)
	}

	if err := DB.AutoMigrate(&TestEntity{}); err != nil {
		t.Fatalf("failed to migrate test entity: %v", err)
	}

	return DB
}
