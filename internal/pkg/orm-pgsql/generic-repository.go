package ormpgsql

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// GenericRepository provides a generic repository for CRUD operations on any entity type.
type GenericRepository[T any] struct {
	DB *gorm.DB
}

// NewGenericRepository creates a new instance of GenericRepository.
func NewGenericRepository[T any](DB *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{DB: DB}
}

// Create inserts a new record into the database.
func (r *GenericRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

// CreateMany inserts multiple records into the database.
func (r *GenericRepository[T]) CreateMany(ctx context.Context, entities *[]T) error {
	return r.DB.WithContext(ctx).Create(entities).Error
}

// GetById retrieves a single record based on the provided ID.
func (r *GenericRepository[T]) GetById(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).
		Model(&entity).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Get retrieves a single record based on the provided parameters.
func (r *GenericRepository[T]) Get(ctx context.Context, params *T) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).
		Where(params).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// GetAll retrieves all records from the database.
func (r *GenericRepository[T]) GetAll(ctx context.Context) (*[]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

// Where retrieves records based on the provided parameters.
func (r *GenericRepository[T]) Where(ctx context.Context, params *T) (*[]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).
		Where(params).
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

// Update modifies an existing record in the database.
func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

// UpdateMany modifies multiple records in the database.
func (r *GenericRepository[T]) UpdateMany(ctx context.Context, entities *[]T) error {
	tx := r.DB.WithContext(ctx).Begin()
	for _, entity := range *entities {
		if err := tx.Save(&entity).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// Delete sets the deleted_at timestamp for soft deletion.
func (r *GenericRepository[T]) Delete(ctx context.Context, entityId string) error {
	var entity T
	err := r.DB.WithContext(ctx).
		Model(&entity).
		Where("id = ?", entityId).
		UpdateColumn("deleted_at", time.Now().UTC()).
		Error
	if err != nil {
		return err
	}
	return nil
}

// SkipTake retrieves records with pagination support.
func (r *GenericRepository[T]) SkipTake(ctx context.Context, skip int, take int) (*[]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).
		Offset(skip).
		Limit(take).
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

// Count returns the total number of records.
func (r *GenericRepository[T]) Count(ctx context.Context) (int64, error) {
	var entity T
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&entity).
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountWhere returns the number of records that match the provided parameters.
func (r *GenericRepository[T]) CountWhere(ctx context.Context, params *T) (int64, error) {
	var entity T
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&entity).
		Where(params).
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
