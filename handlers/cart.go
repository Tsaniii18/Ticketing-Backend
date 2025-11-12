package handlers

import (
	"time"

	"github.com/Tsaniii18/Ticketing-Backend/config"
	"github.com/Tsaniii18/Ticketing-Backend/models"
	"github.com/Tsaniii18/Ticketing-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AddToCart(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var cartData struct {
		TicketCategoryID string `json:"ticket_category_id"`
		Quantity         uint   `json:"quantity"`
	}

	if err := c.BodyParser(&cartData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	// Get ticket category
	var ticketCategory models.TicketCategory
	if err := config.DB.First(&ticketCategory, "ticket_category_id = ?", cartData.TicketCategoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket category not found",
		})
	}

	// Check quota
	if ticketCategory.Sold+cartData.Quantity > ticketCategory.Quota {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Not enough quota available",
		})
	}

	priceTotal := float64(cartData.Quantity) * ticketCategory.Price

	cart := models.Cart{
		CartID:           utils.GenerateCartID(),
		TicketCategoryID: ticketCategory.TicketCategoryID,
		OwnerID:          user.UserID,
		Quantity:         cartData.Quantity,
		PriceTotal:       priceTotal,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := config.DB.Create(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add to cart",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(cart)
}

func GetCart(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var cart []models.Cart
	if err := config.DB.Preload("TicketCategory").Preload("TicketCategory.Event").Where("owner_id = ?", user.UserID).Find(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch cart",
		})
	}

	return c.JSON(cart)
}

func UpdateCart(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var updateData struct {
		CartID   uint `json:"cart_id"`
		Quantity uint `json:"quantity"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	var cart models.Cart
	if err := config.DB.Preload("TicketCategory").Where("cart_id = ? AND owner_id = ?", updateData.CartID, user.UserID).First(&cart).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart item not found",
		})
	}

	// Check quota
	var ticketCategory models.TicketCategory
	if err := config.DB.First(&ticketCategory, "ticket_category_id = ?", cart.TicketCategoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ticket category not found",
		})
	}

	if ticketCategory.Sold+updateData.Quantity > ticketCategory.Quota {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Not enough quota available",
		})
	}

	cart.Quantity = updateData.Quantity
	cart.PriceTotal = float64(updateData.Quantity) * ticketCategory.Price

	if err := config.DB.Save(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update cart",
		})
	}

	return c.JSON(cart)
}

func DeleteCart(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var deleteData struct {
		CartID uint `json:"cart_id"`
	}

	if err := c.BodyParser(&deleteData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if err := config.DB.Where("cart_id = ? AND owner_id = ?", deleteData.CartID, user.UserID).Delete(&models.Cart{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete cart item",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cart item deleted successfully",
	})
}
