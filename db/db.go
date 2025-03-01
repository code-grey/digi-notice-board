package db

import (
    "log"
	"digi-notice-board/models"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// Connect initializes the DB connection
func Connect() {
    dsn := "root:Rubidium@85@tcp(127.0.0.1:3306)/digital_notice_db?charset=utf8mb4&parseTime=True&loc=Local"
    var err error


	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

	DB.AutoMigrate(&models.Announcement{})

	if err := DB.AutoMigrate(&models.Announcement{}); err != nil {
        log.Fatal("Failed to auto-migrate:", err)
    }
	
    /* If you want to auto-migrate models, e.g.:
    / DB.AutoMigrate(&models.Announcement{})*/
}

