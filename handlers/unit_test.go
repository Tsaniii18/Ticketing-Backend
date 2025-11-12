package handlers

import (
	// "context"
	"log"
	// "os"
	"testing"
	"time"

	"fmt"
	// "github.com/gofiber/fiber/v2"
	"github.com/Tsaniii18/Ticketing-Backend/config"
	"github.com/Tsaniii18/Ticketing-Backend/models"
	"github.com/Tsaniii18/Ticketing-Backend/utils"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func TestCreateEvent(t *testing.T) {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  Warning: .env file not loaded, using system environment")
	}
	config.ConnectDatabase()
	config.InitCloudinary()
	err := migrateDatabase(config.DB)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	user := models.User{UserID: utils.GenerateEventID(), Username: utils.GenerateCartID(), RegisterStatus: "approved", Name: utils.GenerateRandomName(), Email: utils.GenerateRandomEmail(), Password: "password", Role: "user", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := config.DB.Create(&user).Error; err != nil {
		log.Fatal("Failed to create test user:", err)
	} else {
		log.Println("Test user created successfully")
	}

	event := models.Event{
		EventID:     utils.GenerateEventID(),
		Name:        "Test User Event",
		OwnerID:     user.UserID,
		Status:      "approved",
		DateStart:   time.Now().Add(24 * time.Hour),
		DateEnd:     time.Now().Add(34 * time.Hour),
		Location:    "Test Location",
		Description: "Test Description",
		Image:       "http://example.com/image.jpg",
		Flyer:       "http://example.com/flyer.jpg",
		Category:    "Test Category",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&event).Error; err != nil {
		panic(err)
	} else {
		log.Println("Test event created successfully")
	}

	// ticketCategory := models.TicketCategory{
	// 	Name:             "VIP",
	// 	TicketCategoryID: utils.GenerateTicketCategoryID(),
	// 	EventID:          event.EventID,
	// 	Price:            100.0,
	// 	Quota:            100,
	// 	Description:      "Test Ticket Category",
	// 	CreatedAt:        time.Now(),
	// 	UpdatedAt:        time.Now(),
	// }

	// if err := config.DB.Create(&ticketCategory).Error; err != nil {
	// 	panic(err)
	// }

}

// func TestTicketCategory(t *testing.T) {

//     if err := godotenv.Load(".env"); err != nil {
//         log.Println("⚠️  Warning: .env file not loaded, using system environment")
//     }
//     config.ConnectDatabase()
// 	config.InitCloudinary()
// 	err := migrateDatabase(config.DB)
// 	if err != nil {
// 		log.Fatal("Failed to migrate database:", err)
// 	}

// }

func migrateDatabase(db *gorm.DB) error {
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.Event{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.TicketCategory{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.Cart{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.TransactionHistory{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.Ticket{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.TransactionDetail{})
	if err != nil {
		return err
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	fmt.Println("Database migrated successfully")
	return nil
}
