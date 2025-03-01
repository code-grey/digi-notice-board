package db

import (
    "log"
    "os"
	"digi-notice-board/models"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// Connect initializes the DB connection
func Connect() {

	user := os.Getenv("MYSQL_USER")       // e.g. "root"
    pass := os.Getenv("MYSQL_PASSWORD")   // e.g. "secretpassword"
    host := os.Getenv("MYSQL_HOST")       // e.g. "127.0.0.1"
    port := os.Getenv("MYSQL_PORT")       // e.g. "3306"
    dbName := os.Getenv("MYSQL_DB")       // e.g. "digital_notice_db"
    
       if user == "" || pass == "" || host == "" || port == "" || dbName == "" {
        log.Fatal("Missing one or more required environment variables")
    }


        dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

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

