package postgrescontainer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NekKkMirror/go-app/internal/pkg/orm-pgsql"
	"github.com/cenkalti/backoff/v4"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

// Options holds configuration for the PostgreSQL container.
type Options struct {
	Database  string
	Host      string
	Port      nat.Port
	UserName  string
	Password  string
	ImageName string
	Name      string
	Tag       string
	Timeout   time.Duration
}

// Start initializes a PostgreSQL container and returns a gorm DB instance, sqlmock, and any error occurred.
func Start(ctx context.Context, t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	options := getDefaultPostgresOptions()
	containerReq := getContainerRequest(options)

	postgresContainer, err := startContainer(ctx, containerReq)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to start PostgreSQL container")
	}

	t.Cleanup(func() {
		if err := stopContainer(ctx, postgresContainer); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	DB, err := createORMConnection(ctx, postgresContainer, options)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create ORM connection")
	}

	mock, err := setupSQLMock()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create mock object")
	}

	if err := loadSeed(DB); err != nil {
		return nil, nil, err
	}

	return DB, mock, nil
}

// getDefaultPostgresOptions returns the default configuration for PostgreSQL container.
func getDefaultPostgresOptions() *Options {
	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		panic(errors.Wrap(err, "failed to create new port")) // handle this appropriately
	}

	return &Options{
		Database:  "test_db",
		Port:      port,
		Host:      "localhost",
		UserName:  "testcontainers",
		Password:  "testcontainers",
		Tag:       "latest",
		ImageName: "postgres",
		Name:      "postgresql-testcontainer",
		Timeout:   5 * time.Minute,
	}
}

// getContainerRequest builds and returns a testcontainers.ContainerRequest using the provided options.
func getContainerRequest(opts *Options) testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		ExposedPorts: []string{opts.Port.Port()},
		WaitingFor:   wait.ForListeningPort(opts.Port),
		Env: map[string]string{
			"POSTGRES_DB":       opts.Database,
			"POSTGRES_PASSWORD": opts.Password,
			"POSTGRES_USER":     opts.UserName,
		},
	}
}

// startContainer starts a new testcontainer with given request configuration.
func startContainer(ctx context.Context, req testcontainers.ContainerRequest) (testcontainers.Container, error) {
	genericReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}
	return testcontainers.GenericContainer(ctx, genericReq)
}

// stopContainer stops and terminates the given testcontainer.
func stopContainer(ctx context.Context, container testcontainers.Container) error {
	return errors.Wrap(container.Terminate(ctx), "failed to terminate container")
}

// createORMConnection establishes a GORM connection using provided PostgreSQL container and options.
func createORMConnection(ctx context.Context, container testcontainers.Container, opts *Options) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second
	const maxRetries = 5

	var (
		DB     *gorm.DB
		err    error
		config *ormpgsql.PostgresConfig
	)

	err = backoff.Retry(func() error {
		opts.Host, err = container.Host(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get container host")
		}

		opts.Port, err = container.MappedPort(ctx, opts.Port)
		if err != nil {
			return errors.Wrap(err, "failed to get exposed container port")
		}

		config = &ormpgsql.PostgresConfig{
			Port:     opts.Port.Int(),
			Host:     opts.Host,
			DBName:   opts.Database,
			User:     opts.UserName,
			Password: opts.Password,
			SSLMode:  false,
		}

		DB, err = ormpgsql.NewORM(config)
		return err
	}, backoff.WithMaxRetries(bo, maxRetries))

	if err != nil {
		return nil, errors.Wrap(err, "failed to create connection after retries")
	}

	return DB, nil
}

// setupSQLMock initializes and returns sqlmock instance.
func setupSQLMock() (sqlmock.Sqlmock, error) {
	_, mock, err := sqlmock.New()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sqlmock")
	}
	return mock, nil
}

// loadSeed inserts dummy data into the database tables for testing purposes.
func loadSeed(DB *gorm.DB) error {
	if err := addUsersSeed(DB); err != nil {
		return fmt.Errorf("failed to load users seed data: %w", err)
	}
	return nil
}

type User struct {
	ID                     int
	Name                   string
	Age                    int
	Email                  string
	Phone                  string
	Address                string
	DateOfBirth            string
	CreatedAt              string
	UpdatedAt              string
	DeletedAt              *string
	IsActive               bool
	IsAdmin                bool
	SubscriptionType       string
	Price                  float64
	Discount               float64
	IsVerified             bool
	LastLogin              string
	TotalSpent             float64
	AccountBalance         float64
	NumberOfOrders         int
	IsEmailSubscribed      bool
	PreferredLanguage      string
	LocationLatitude       float64
	LocationLongitude      float64
	ReferralCode           string
	PreferredCurrency      string
	IsOnNewsletter         bool
	PaymentMethod          string
	Gender                 string
	IPAddress              string
	DeviceType             string
	AgeGroup               string
	MaritalStatus          string
	NumberOfChildren       int
	EducationLevel         string
	Occupation             string
	LastPurchaseDate       string
	FirstPurchaseDate      string
	NewsletterOptIn        bool
	LastTransactionID      string
	BankName               string
	BankAccountNumber      string
	NationalID             string
	DateOfJoining          string
	EmergencyContactName   string
	EmergencyContactPhone  string
	SocialSecurityNumber   string
	PreferredContactMethod string
	PreferredStore         string
	ProfilePictureURL      string
	IsLoyalCustomer        bool
	IsBanned               bool
	LastReviewDate         string
	DeviceOS               string
	Timezone               string
	LoyaltyPoints          int
	ReferralPoints         int
	IsSubscriptionActive   bool
	BirthCountry           string
	TaxID                  string
	PreferredPaymentMethod string
}

// addUsersSeed migrates the User struct and inserts dummy data into the database tables for testing purposes.
func addUsersSeed(DB *gorm.DB) error {
	err := DB.AutoMigrate(&User{})

	if err != nil {
		return err
	}

	// Insert dummy data into the database
	users := generateDummyUsers()
	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

// generateDummyUsers creates and returns a slice of dummy User data for testing purposes.
func generateDummyUsers() []User {
	users := []User{
		{
			ID:                     1,
			Name:                   "Alice Smith",
			Age:                    28,
			Email:                  "alice@example.com",
			Phone:                  "+123456789",
			Address:                "123 Main St, City",
			DateOfBirth:            "1995-02-15",
			CreatedAt:              "2023-01-01T10:00:00Z",
			UpdatedAt:              "2023-10-21T15:30:00Z",
			DeletedAt:              nil,
			IsActive:               true,
			IsAdmin:                false,
			SubscriptionType:       "Premium",
			Price:                  99.99,
			Discount:               10.0,
			IsVerified:             true,
			LastLogin:              "2023-10-20T09:45:00Z",
			TotalSpent:             1200.50,
			AccountBalance:         200.75,
			NumberOfOrders:         25,
			IsEmailSubscribed:      true,
			PreferredLanguage:      "English",
			LocationLatitude:       40.7128,
			LocationLongitude:      -74.0060,
			ReferralCode:           "REF123",
			PreferredCurrency:      "USD",
			IsOnNewsletter:         true,
			PaymentMethod:          "Credit Card",
			Gender:                 "Female",
			IPAddress:              "192.168.1.1",
			DeviceType:             "Mobile",
			AgeGroup:               "25-34",
			MaritalStatus:          "Single",
			NumberOfChildren:       0,
			EducationLevel:         "Bachelor's",
			Occupation:             "Software Developer",
			LastPurchaseDate:       "2023-10-10T12:30:00Z",
			FirstPurchaseDate:      "2021-06-15T08:00:00Z",
			NewsletterOptIn:        true,
			LastTransactionID:      "TXN987654321",
			BankName:               "ABC Bank",
			BankAccountNumber:      "1234567890",
			NationalID:             "ID12345",
			DateOfJoining:          "2021-06-01T10:00:00Z",
			EmergencyContactName:   "John Doe",
			EmergencyContactPhone:  "+987654321",
			SocialSecurityNumber:   "SSN123456",
			PreferredContactMethod: "Email",
			PreferredStore:         "Online Store",
			ProfilePictureURL:      "https://example.com/images/alice.jpg",
			IsLoyalCustomer:        true,
			IsBanned:               false,
			LastReviewDate:         "2023-09-25T11:00:00Z",
			DeviceOS:               "iOS",
			Timezone:               "EST",
			LoyaltyPoints:          500,
			ReferralPoints:         150,
			IsSubscriptionActive:   true,
			BirthCountry:           "USA",
			TaxID:                  "TAX123456",
			PreferredPaymentMethod: "Credit Card",
		},
	}

	// Generate additional dummy users
	for i := 2; i <= 40; i++ {
		var deletedAt *string
		if i%10 == 0 {
			date := "2023-08-01T12:00:00Z"
			deletedAt = &date
		}

		users = append(users, User{
			ID:                     i,
			Name:                   fmt.Sprintf("User %d", i),
			Age:                    20 + (i % 30),
			Email:                  fmt.Sprintf("user%d@example.com", i),
			Phone:                  fmt.Sprintf("+12345678%d", i),
			Address:                fmt.Sprintf("%d Example St, City", i),
			DateOfBirth:            fmt.Sprintf("199%d-01-01", i%10),
			CreatedAt:              "2023-01-01T10:00:00Z",
			UpdatedAt:              "2023-10-21T15:30:00Z",
			DeletedAt:              deletedAt,
			IsActive:               i%2 == 0,
			IsAdmin:                i%10 == 0,
			SubscriptionType:       "Basic",
			Price:                  49.99,
			Discount:               float64(i % 5),
			IsVerified:             i%2 == 0,
			LastLogin:              "2023-10-20T09:45:00Z",
			TotalSpent:             float64(100 * i),
			AccountBalance:         float64(50 * i % 100),
			NumberOfOrders:         i * 2,
			IsEmailSubscribed:      i%3 == 0,
			PreferredLanguage:      "English",
			LocationLatitude:       35.0 + float64(i)/10.0,
			LocationLongitude:      -80.0 - float64(i)/10.0,
			ReferralCode:           fmt.Sprintf("REF%d", i),
			PreferredCurrency:      "USD",
			IsOnNewsletter:         i%2 == 0,
			PaymentMethod:          "Credit Card",
			Gender:                 "Other",
			IPAddress:              fmt.Sprintf("192.168.1.%d", i),
			DeviceType:             "Mobile",
			AgeGroup:               "25-34",
			MaritalStatus:          "Single",
			NumberOfChildren:       0,
			EducationLevel:         "High School",
			Occupation:             "Student",
			LastPurchaseDate:       "2023-10-05T12:30:00Z",
			FirstPurchaseDate:      "2020-12-10T08:00:00Z",
			NewsletterOptIn:        i%2 == 0,
			LastTransactionID:      fmt.Sprintf("TXN%d", i),
			BankName:               "Bank",
			BankAccountNumber:      fmt.Sprintf("%d", i*12345),
			NationalID:             fmt.Sprintf("ID%d", i),
			DateOfJoining:          "2020-11-30T10:00:00Z",
			EmergencyContactName:   "Emergency Contact",
			EmergencyContactPhone:  "+999999999",
			SocialSecurityNumber:   "SSN000000",
			PreferredContactMethod: "Email",
			PreferredStore:         "Online",
			ProfilePictureURL:      "https://example.com/default.jpg",
			IsLoyalCustomer:        true,
			IsBanned:               false,
			LastReviewDate:         "2023-09-25T11:00:00Z",
			DeviceOS:               "iOS",
			Timezone:               "EST",
			LoyaltyPoints:          100 * i,
			ReferralPoints:         50 * i,
			IsSubscriptionActive:   true,
			BirthCountry:           "USA",
			TaxID:                  "TAX000000",
			PreferredPaymentMethod: "Credit Card",
		})
	}

	return users
}
