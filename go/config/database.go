package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"fmt"
	"github.com/Shota0616/go-sns/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "user:password@tcp(mysql:3306)/sns?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	fmt.Println("Database connected!")
	DB = database
}

func MigrateDatabase() {
	if DB == nil {
		panic("Database connection is not initialized!")
	}
	DB.AutoMigrate(&models.User{})
	fmt.Println("Database migrated!")
}

// func GetDB() *gorm.DB {
// 	return DB
// }
