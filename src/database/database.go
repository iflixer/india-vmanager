package database

import (
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service struct {
	DB *gorm.DB
}

func NewService(mysqlURL string) (*Service, error) {

	db, err := gorm.Open(mysql.Open(mysqlURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to open database:" + err.Error())
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(20)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	debugMode := os.Getenv("MYSQL_DEBUG_MODE")
	if debugMode != "" {
		d, _ := strconv.Atoi(debugMode)
		db.Logger = logger.Default.LogMode(logger.LogLevel(d))
	}

	return &Service{
		DB: db,
	}, nil
}
