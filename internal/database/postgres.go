package database

import (
	"context"

	"github.com/Rhaqim/creds/internal/config"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB // Global database connection

// InitDB initializes the database connection using the configuration values from the config package.
// It establishes a connection to the PostgreSQL database and assigns the connection to the global DB variable.
// If an error occurs during the connection process, it logs the error and shuts down the logger.
func Init() {
	var err error

	var dsn string = "host=" + config.PgHost + " port=" + config.PgPort + " user=" + config.PgUser + " dbname=" + config.Database + " sslmode=" + config.SSLMode + " password=" + config.PgPassword

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logger.Default.Error(context.Background(), "failed to connect database")
		panic("failed to connect database")
	}
}

func Insert[T any](t *T) error {
	result := DB.Model(&t).Create(&t)
	return result.Error
}

func InsertMany[T any](ts *[]T) error {
	tx := DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, t := range *ts {
		if err := tx.Model(&t).Create(&t).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
func Query[T any](t *T, query interface{}, args ...interface{}) error {
	result := DB.Model(&t).Where(query, args...).First(&t)
	return result.Error
}

func QueryAll[T any](ts *[]T, query interface{}, args ...interface{}) error {
	result := DB.Model(&ts).Where(query, args...).Find(&ts)
	return result.Error
}

func Update[T any](t *T, query interface{}, args ...interface{}) (*gorm.DB, error) {
	DB.Save(&t)
	result := DB.Model(&t).Where(query, args...).Updates(t)
	return result, result.Error
}

func Delete[T any](t *T, query interface{}, args ...interface{}) (*gorm.DB, error) {
	result := DB.Model(&t).Where(query, args...).Delete(&t)
	return result, result.Error
}

// implement DutchDatabase interface

type PostgresDatabase struct {
	db *gorm.DB
}

func NewPostgresDatabase[T any](post_db *gorm.DB, model *T) *PostgresDatabase {

	return &PostgresDatabase{
		db: post_db.Model(model),
	}
}

func (db *PostgresDatabase) Insert(data interface{}) (interface{}, error) {
	result := db.db.Create(&data)

	return data, result.Error
}

func (db *PostgresDatabase) Update(data interface{}, filter interface{}, args ...interface{}) (interface{}, error) {
	result := db.db.Where(filter, args...).Updates(data)

	return data, result.Error
}

func (db *PostgresDatabase) Get(filter interface{}) (interface{}, error) {
	var result interface{}
	result_ := db.db.First(&result, filter)

	return result, result_.Error
}

func (db *PostgresDatabase) GetAll(filter interface{}) ([]interface{}, error) {
	var results []interface{}
	if err := db.db.Find(&results, filter).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (db *PostgresDatabase) Delete(filter interface{}) (interface{}, error) {
	var result interface{}
	tx := db.db.Begin()
	if err := tx.First(&result, filter).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Delete(&result, filter).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return &result, tx.Commit().Error
}

func (p *PostgresDatabase) Disconnect() error {
	// if p.db != nil {
	// 	return p.db.Close()
	// }
	return nil
}
