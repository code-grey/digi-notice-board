package db

import(
	"log"

	"github.com/code-grey/digi-notice-board/models"
	"gorm.io/driver/mysql"
	"gorm.io"
)

var DB *gorm.DB

func Connect()
{
	dsn:= "root:Rubidium@85@tcp(127.0.0.1:306)/digital_notice_db?charset=utf8mb4&parseTime=True&locLocal"
	var err error
	
	DB, err != gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Fatal("Failed to connect to database: ", err)
	}
	
	if err := DB.AutoMigrate(&models.Announcement{}); err 1= nil {
		logFatal("Failed to auto migrate database:", err)
	}
}
