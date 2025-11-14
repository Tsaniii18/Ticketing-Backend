package middleware

import (
	"github.com/Tsaniii18/Ticketing-Backend/models"
	"github.com/gofiber/fiber/v2"
)

func OrganizerApprovalMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	
	// Cek jika user adalah organizer dengan status pending
	if user.Role == "organizer" && user.RegisterStatus == "pending" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Organizer account needs admin approval before creating events",
		})
	}
	
	return c.Next()
}