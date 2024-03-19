package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"net/url"
	"os" // Import the os package
	"time"
)

var globalDb *gorm.DB

// getEnv retrieves environment variables or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// createDatabaseURL constructs a PostgreSQL connection string using environment variables
func createDatabaseURL() string {
	host := getEnv("DATABASE_IP", "localhost")
	port := getEnv("DATABASE_PORT", "5433")
	user := getEnv("DATABASE_USER", "fixer")
	password := getEnv("DATABASE_PASSWORD", "")
	dbname := getEnv("DATABASE_NAME", "postgres")

	// Manually construct the URL, ensuring special characters in the password are encoded
	password = url.QueryEscape(password)
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	return connectionString
}

// CreatePool initializes the global database connection pool using
// environment variables. The function configures the database connection pool with
// predefined settings for maximum idle connections, maximum open connections, and the maximum
// lifetime of a connection. If an error occurs while
// establishing a connection to the database, including setting up the connection pool,
// CreatePool returns an error.
func CreatePool() (err error) {

	dsn := createDatabaseURL()

	globalDb, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "palabras.", // schema name
			SingularTable: true,
		},
	})

	sqlDB, err := globalDb.DB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = sqlDB.Ping()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Created %d db connections\n", sqlDB.Stats().OpenConnections)
	return nil
}

func GetConnection() (db *gorm.DB, err error) {

	if globalDb == nil {
		return nil, fmt.Errorf("database not connected")
	}

	db = globalDb
	return
}
