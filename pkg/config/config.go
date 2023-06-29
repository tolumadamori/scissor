package config

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDB() {

	name := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dsn := string(name) + ":" + string(password) + "@tcp(" + string(dbHost) + ":" + string(dbPort) + ")/" + string(dbName) + "?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the DB")
	}

}
