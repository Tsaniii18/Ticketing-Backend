package config

import (
    "log"
    "os"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        dsn = "root:tsu3182@tcp(127.0.0.1:3306)/ticketingdb?charset=utf8mb4&parseTime=True&loc=Local"
    }

    database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    DB = database
    log.Println("Database connected successfully")
}