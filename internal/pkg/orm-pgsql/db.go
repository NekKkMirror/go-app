package ormpgsql

import (
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

// PostgresConfig holds the configuration parameters for the PostgreSQL connection.
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbName"`
	SSLMode  bool   `mapstructure:"sslMode"`
	Password string `mapstructure:"password"`
}

// ORM represents an object-relational mapper with a GORM DB connection and configuration.
type ORM struct {
	DB     *gorm.DB
	config *PostgresConfig
}

// NewORM initializes and returns a new ORM instance with a connected GORM database.
// It handles connection retries using exponential backoff and ensures the database exists.
func NewORM(cfg *PostgresConfig) (*gorm.DB, error) {
	if cfg.DBName == "" {
		return nil, errors.New("database name is required")
	}

	if err := createDB(cfg); err != nil {
		return nil, err
	}

	dataSrcName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second
	maxRetries := 5

	var db *gorm.DB
	err := backoff.Retry(func() error {
		var err error
		db, err = gorm.Open(postgres.Open(dataSrcName), &gorm.Config{})
		if err != nil {
			return errors.Wrapf(err, "failed to connect to postgres: %s", dataSrcName)
		}
		return nil
	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	if err != nil {
		return nil, err
	}

	return db, nil
}

// createDB creates the database if it does not already exist, based on the provided configuration.
func createDB(cfg *PostgresConfig) error {
	// DSN without specifying a database to connect on server level
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	defer sqldb.Close()

	var exists bool
	query := fmt.Sprintf("SELECT 1 FROM pg_catalog.pg_database WHERE datname='%s'", cfg.DBName)
	if err := sqldb.QueryRow(query).Scan(&exists); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(err, "failed to check database existence")
		}
		exists = false
	}

	if exists {
		return nil
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
	if _, err := sqldb.Exec(createDBQuery); err != nil {
		return errors.Wrap(err, "failed to create database")
	}

	return nil
}

// Close closes the GORM database connection associated with the ORM instance.
func (orm *ORM) Close() error {
	db, err := orm.DB.DB()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve db from gorm DB")
	}
	return db.Close()
}

// Paginate fetches the records as per the pagination and filter criteria.
func Paginate[T any](listQuery *pagination.ListQuery, DB *gorm.DB) (*pagination.ListResult[T], error) {
	var data []T
	var totalCount int64
	var query *gorm.DB
	var err error

	if err = DB.Model(new(T)).Count(&totalCount).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count total records")
	}

	query = DB.Offset(listQuery.GetOffset()).
		Limit(listQuery.GetLimit()).
		Order(listQuery.GetOrderBy())

	if listQuery.Filters != nil {
		query, err = pagination.ApplyFilterAction(query, listQuery.Filters, make(map[string]bool))
		if err != nil {
			return nil, err
		}
	}

	if err = query.Find(&data).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch data")
	}

	listResult := pagination.NewListResult(listQuery.Size, listQuery.Page, totalCount, data)

	return listResult, nil
}
