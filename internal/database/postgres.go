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

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Default.Error(context.Background(), "failed to close database")
		panic("failed to close database")
	}
	sqlDB.Close()
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

// implement Database interface
type Database[T any] interface {
	Insert() error
	Get(query interface{}, args ...interface{}) error
	GetAll(query interface{}, args ...interface{}) ([]T, error)
	Update(query interface{}, args ...interface{}) error
	Delete(query interface{}, args ...interface{}) error
}

type PostgresDatabase[T any] struct {
	db    *gorm.DB
	model *T
}

func NewPostgresDatabase[T any](model *T) Database[T] {

	return &PostgresDatabase[T]{
		db: DB.Model(model),
	}
}

func (P *PostgresDatabase[T]) Insert() error {
	return P.db.Create(P.model).Error
}

func (P *PostgresDatabase[T]) Get(query interface{}, args ...interface{}) error {
	return P.db.Where(query, args...).First(P.model).Error
}

func (P *PostgresDatabase[T]) GetAll(query interface{}, args ...interface{}) ([]T, error) {
	var results []T
	result := P.db.Where(query, args...).Find(&results)
	return results, result.Error
}

func (P *PostgresDatabase[T]) Update(query interface{}, args ...interface{}) error {
	return P.db.Where(query, args...).Updates(P.model).Error
}

func (P *PostgresDatabase[T]) Delete(query interface{}, args ...interface{}) error {
	return P.db.Where(query, args...).Delete(P.model).Error
}
