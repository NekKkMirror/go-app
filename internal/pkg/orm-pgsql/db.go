package gormpgsql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/utils/db/pagination"
	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"github.com/uptrace/bun/driver/pgdriver"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbName"`
	SSLMode  bool   `mapstructure:"sslMode"`
	Password string `mapstructure:"password"`
}

type ORM struct {
	DB     *gorm.DB
	config *PostgresConfig
}

// NewORM creates a new instance of the ORM struct with a configured GORM database connection.
// It attempts to connect to the PostgreSQL database using the provided configuration and handles connection retries.
// If the database does not exist, it creates it.
//
// Parameters:
//   - cfg: A pointer to a PostgresConfig struct containing the database connection configuration.
//     The PostgresConfig struct should have the following fields:
//   - Host: The host of the PostgreSQL server.
//   - Port: The port of the PostgreSQL server.
//   - User: The username for the PostgreSQL server.
//   - DBName: The name of the database to connect to.
//   - SSLMode: A boolean indicating whether to use SSL for the connection.
//   - Password: The password for the PostgreSQL server.
//
// Returns:
// - A pointer to a gorm.DB instance representing the database connection.
// - An error if the database connection fails or if the database creation fails.
func NewORM(cfg *PostgresConfig) (*gorm.DB, error) {
	var dataSrcName string

	if cfg.DBName == "" {
		return nil, errors.New("db name is required")
	}

	err := createDB(cfg)
	if err != nil {
		return nil, err
	}

	dataSrcName = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second
	maxRetries := 5

	var ORMDB *gorm.DB

	err = backoff.Retry(func() error {
		ORMDB, err = gorm.Open(postgres.Open(dataSrcName), &gorm.Config{})

		if err != nil {
			return errors.Errorf("failed to connect postgres: %v and connection information: %s", err, dataSrcName)
		}

		return nil
	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	return ORMDB, nil
}

// createDB creates a PostgreSQL database using the provided configuration.
// It first checks if the database already exists. If it does not exist, it creates the database.
//
// Parameters:
//   - cfg: A pointer to a PostgresConfig struct containing the database connection configuration.
//     The PostgresConfig struct should have the following fields:
//   - Host: The host of the PostgreSQL server.
//   - Port: The port of the PostgreSQL server.
//   - User: The username for the PostgreSQL server.
//   - DBName: The name of the database to connect to.
//   - SSLMode: A boolean indicating whether to use SSL for the connection.
//   - Password: The password for the PostgreSQL server.
//
// Returns:
// - An error if the database creation fails.
// - nil if the database creation is successful.
func createDB(cfg *PostgresConfig) error {
	dataSrc := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		"postgres",
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dataSrc)))

	var exists int
	rows, err := sqldb.Query(fmt.Sprintf("SELECT 1 FROM pg_catalog_database WHERE datname = '%s'", cfg.DBName))
	if err != nil {
		return err
	}

	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists == 1 {
		return nil
	}

	_, err = sqldb.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer func(sqldb *sql.DB) {
		err := sqldb.Close()
		if err != nil {

		}
	}(sqldb)

	return nil
}

// Close closes the database connection associated with the ORM instance.
// It retrieves the underlying database connection from the GORM instance and closes it.
//
// Parameters:
//   - orm: A pointer to the ORM struct instance.
//
// Returns:
//   - An error if the database connection cannot be closed.
//   - nil if the database connection is successfully closed.
func (orm *ORM) Close() error {
	d, err := orm.DB.DB()
	if err != nil {
		return err
	}

	err = d.Close()
	if err != nil {
		return err
	}

	return nil
}

// Migrate applies database migrations to the provided GORM database instance.
// It iterates through the provided types and applies AutoMigrate to each one.
// If any migration fails, it returns an error.
//
// Parameters:
//   - gORM: A pointer to a gorm.DB instance representing the database connection.
//   - types: Variadic parameter of type interface{}, representing the types for which migrations need to be applied.
//
// Returns:
//   - An error if any migration fails.
//   - nil if all migrations are successful.
func Migrate(gORM *gorm.DB, types ...interface{}) error {
	for _, t := range types {
		err := gORM.AutoMigrate(t)
		if err != nil {
			return err
		}
	}
	return nil
}

// Paginate retrieves a paginated list of data from the database using GORM.
// It applies filters, sorting, and pagination based on the provided ListQuery.
//
// Parameters:
//   - ctx: A context.Context for managing the lifecycle of the request.
//   - listQuery: A pointer to a pagination.ListQuery struct containing the pagination, sorting, and filtering parameters.
//   - DB: A pointer to a gorm.DB instance representing the database connection.
//
// Returns:
//   - A pointer to a pagination.ListResult[T] containing the paginated data and metadata.
//   - An error if any database operation fails.
func Paginate[T any](ctx context.Context, listQuery *pagination.ListQuery, DB *gorm.DB) (*pagination.ListResult[T], error) {
	var data []T
	var totalCount int64
	DB.Model(data).Count(&totalCount)

	query := DB.Offset(listQuery.GetOffset()).
		Limit(listQuery.GetLimit()).
		Order(listQuery.GetOrderBy())

	var err error

	if listQuery.Filters != nil {
		for _, filter := range listQuery.Filters {
			query, err = pagination.ApplyFilterAction(query, filter.Field, filter.Value, filter.Comparison)
		}
	}

	if err = query.Find(&data).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch data")
	}

	return pagination.NewListResult[T](listQuery.GetSize(), listQuery.GetPage(), totalCount, data), nil
}
