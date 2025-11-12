package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/Tsaniii18/Ticketing-Backend/config"
    "github.com/Tsaniii18/Ticketing-Backend/models"
)

func GetTickets(c *fiber.Ctx) error {
    user := c.Locals("user").(models.User)
    
    var tickets []models.Ticket
    if err := config.DB.Preload("Event").Preload("TicketCategory").Where("owner_id = ?", user.UserID).Find(&tickets).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch tickets",
        })
    }
    
    return c.JSON(tickets)
}

func CheckInTicket(c *fiber.Ctx) error {
    ticketID := c.Params("id")
    
    var ticket models.Ticket
    if err := config.DB.Preload("Event").First(&ticket, ticketID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Ticket not found",
        })
    }
    
    if ticket.Status == "used" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Ticket already used",
        })
    }
    
    ticket.Status = "used"
    if err := config.DB.Save(&ticket).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to check in ticket",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Ticket checked in successfully",
        "ticket":  ticket,
    })
}