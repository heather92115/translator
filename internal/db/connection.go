package db

import (
	"fmt"
	"github.com/heather92115/verdure-admin/internal/mdl"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var globalDb *gorm.DB

// CreatePool initializes the global db connection pool using
// environment variables. The function configures the db connection pool with
// predefined settings for maximum idle connections, maximum open connections, and the maximum
// lifetime of a connection. If an error occurs while
// establishing a connection to the db, including setting up the connection pool,
// CreatePool returns an error.
func CreatePool(dsn string) (err error) {

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

	// SetMaxOpenConns sets the maximum number of open connections to the db.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = sqlDB.Ping()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = MigrateTables()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Created %d db connections\n", sqlDB.Stats().OpenConnections)
	return nil
}

// GetConnection returns a reference to the global database connection.
// It checks if the global database connection (globalDb) has been established.
// If not, it returns an error indicating that the database connection is not available.
//
// Returns:
// - db: A pointer to the gorm.DB instance representing the database connection.
// - err: An error if the global database connection has not been initialized.
//
// Example usage:
// db, err := GetConnection()
//
//	if err != nil {
//	    log.Fatalf("Database connection error: %v", err)
//	}
func GetConnection() (db *gorm.DB, err error) {

	if globalDb == nil {
		return nil, fmt.Errorf("db not connected")
	}

	db = globalDb
	return
}

// CreateEnumIfNotExists checks if a custom ENUM type named 'status_type' exists in the PostgreSQL database.
// If it does not exist, the function creates this ENUM type with predefined values: 'pending', 'in_progress', and 'completed'.
// This function is useful for initializing or migrating databases to ensure that the necessary ENUM types are available
// for use in table definitions or elsewhere within the database schema.
//
// The function executes a PostgreSQL DO block to conditionally create the ENUM type. This approach avoids errors that
// would occur from attempting to create a type that already exists, ensuring idempotency in database migrations or setups.
//
// Parameters:
// - db: A pointer to a gorm.DB instance representing an established database connection.
//
// Returns:
//   - An error if the SQL execution fails, otherwise nil if the ENUM type is successfully checked for existence
//     and created if needed.
//
// Example usage:
//
//	if err := CreateEnumIfNotExists(db); err != nil {
//	    log.Fatalf("Failed to create or check ENUM 'status_type': %v", err)
//	}
//
// Note: This function specifically targets PostgreSQL and uses features unique to that RDBMS.
// It may need adjustments for compatibility with other database systems.
func CreateEnumIfNotExists(db *gorm.DB) error {
	sql := `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type') THEN
				CREATE TYPE status_type AS ENUM ('pending', 'in_progress', 'completed');
			END IF;
		END$$;
		`
	return db.Exec(sql).Error
}

// MigrateTables performs the necessary database migrations to ensure that the schema
// matches the expected structure defined by the internal models. This function is
// typically called during application initialization to prepare the database for use.
//
// The migration process includes the following steps:
//  1. Ensuring that a custom ENUM type 'status_type' exists in the PostgreSQL database,
//     creating it if necessary. This ENUM is used by certain table columns.
//  2. Automatically migrating the database schema to match the structure of the Fixit model.
//  3. Automatically migrating the database schema to match the structure of the Audit model.
//
// Note: This function presumes that the 'vocab' table already exists in the database
// and that its schema matches the structure defined by the internal models. It does not
// perform migration for the 'vocab' table. Ensure that any changes to the vocab model
// are manually reflected in the database or through separate migration scripts.
//
// Returns:
//   - An error if any part of the migration process fails, otherwise nil if all migrations
//     are successful.
//
// Example usage:
//
//	if err := MigrateTables(); err != nil {
//	    log.Fatalf("Database migration failed: %v", err)
//	}
//
// This function utilizes the global database connection (globalDb) to perform migrations.
// It's important to ensure that this global connection is properly initialized and connected
// to the target database before calling MigrateTables.
func MigrateTables() (err error) {

	err = CreateEnumIfNotExists(globalDb)
	if err != nil {
		log.Fatalf("Failed to create enum: %v", err)
	}

	err = globalDb.AutoMigrate(mdl.Fixit{})
	if err != nil {
		return err
	}

	err = globalDb.AutoMigrate(mdl.Audit{})
	if err != nil {
		return err
	}

	return
}
