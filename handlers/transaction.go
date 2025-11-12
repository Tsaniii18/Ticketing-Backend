package handlers

import (
	"time"

	"github.com/Tsaniii18/Ticketing-Backend/config"
	"github.com/Tsaniii18/Ticketing-Backend/models"
	"github.com/Tsaniii18/Ticketing-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Checkout(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Get user's cart items
	var cart []models.Cart
	if err := config.DB.Preload("TicketCategory").Where("owner_id = ?", user.UserID).Find(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch cart",
		})
	}

	if len(cart) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cart is empty",
		})
	}

	// Calculate total
	var total float64
	for _, item := range cart {
		total += item.PriceTotal
	}

	// Create transaction
	transaction := models.TransactionHistory{
		TransactionID:   utils.GenerateTransactionID(),
		OwnerID:         user.UserID,
		TransactionTime: time.Now(),
		PriceTotal:      total,
		CreatedAt:       time.Now(),
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	// Create transaction details and tickets
	for _, cartItem := range cart {
		// Create transaction detail
		transactionDetail := models.TransactionDetail{
			TransactionDetailID: utils.GenerateTransactionDetailID(),
			TicketCategoryID:    cartItem.TicketCategoryID,
			TransactionID:       transaction.TransactionID,
			OwnerID:             user.UserID,
			Quantity:            cartItem.Quantity,
			Subtotal:            cartItem.PriceTotal,
		}

		if err := config.DB.Create(&transactionDetail).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create transaction detail",
			})
		}

		var ticketCategory models.TicketCategory
		if err := config.DB.First(&ticketCategory, "ticket_category_id = ?", cartItem.TicketCategoryID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Ticket category not found",
			})
		}

		// Create tickets
		for i := 0; i < int(cartItem.Quantity); i++ {
			ticket := models.Ticket{
				TicketID:         utils.GenerateTicketID(),
				EventID:          ticketCategory.EventID,
				TicketCategoryID: cartItem.TicketCategoryID,
				OwnerID:          user.UserID,
				Status:           "active",
				Code:             utils.GenerateTicketCode(),
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			if err := config.DB.Create(&ticket).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create ticket",
				})
			}
		}

		// Update sold count
		config.DB.Model(&ticketCategory).Update("sold", ticketCategory.Sold+cartItem.Quantity)

		// Update event total sales
		config.DB.Model(&models.Event{}).Where("event_id = ?", ticketCategory.EventID).
			Update("total_sales", gorm.Expr("total_sales + ?", cartItem.PriceTotal))
	}

	// Clear cart
	if err := config.DB.Where("owner_id = ?", user.UserID).Delete(&models.Cart{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear cart",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Checkout successful",
		"transaction": transaction,
	})
}

func GetTransactions(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var transactions []models.TransactionHistory
	if err := config.DB.Preload("TransactionDetails").Preload("TransactionDetails.TicketCategory").Where("owner_id = ?", user.UserID).Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	return c.JSON(transactions)
}

func GetTransaction(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	transactionID := c.Params("id")

	var transaction models.TransactionHistory
	if err := config.DB.Preload("TransactionDetails").Preload("TransactionDetails.TicketCategory").
		Where("transaction_id = ? AND owner_id = ?", transactionID, user.UserID).
		First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(transaction)
}
