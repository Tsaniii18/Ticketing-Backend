package handlers

import (
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/Tsaniii18/Ticketing-Backend/config"
    "github.com/Tsaniii18/Ticketing-Backend/models"
)

func CreateEvent(c *fiber.Ctx) error {
    user := c.Locals("user").(models.User)
    
    var eventData struct {
        Name        string    `json:"name"`
        DateStart   time.Time `json:"date_start"`
        DateEnd     time.Time `json:"date_end"`
        Location    string    `json:"location"`
        Description string    `json:"description"`
        Image       string    `json:"image"`
        Flyer       string    `json:"flyer"`
        Category    string    `json:"category"`
    }
    
    if err := c.BodyParser(&eventData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }
    
    event := models.Event{
        Name:        eventData.Name,
        OwnerID:     user.UserID,
        Status:      "pending",
        DateStart:   eventData.DateStart,
        DateEnd:     eventData.DateEnd,
        Location:    eventData.Location,
        Description: eventData.Description,
        Image:       eventData.Image,
        Flyer:       eventData.Flyer,
        Category:    eventData.Category,
    }
    
    if err := config.DB.Create(&event).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create event",
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(event)
}

func GetEvents(c *fiber.Ctx) error {
    var events []models.Event
    if err := config.DB.Preload("Owner").Where("status = ?", "approved").Find(&events).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch events",
        })
    }
    
    return c.JSON(events)
}

func GetEvent(c *fiber.Ctx) error {
    eventID := c.Params("id")
    
    var event models.Event
    if err := config.DB.Preload("Owner").Preload("TicketCategories").First(&event, eventID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Event not found",
        })
    }
    
    return c.JSON(event)
}

func UpdateEvent(c *fiber.Ctx) error {
    user := c.Locals("user").(models.User)
    eventID := c.Params("id")
    
    var event models.Event
    if err := config.DB.First(&event, eventID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Event not found",
        })
    }
    
    // Check ownership
    if event.OwnerID != user.UserID && user.Role != "admin" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Not authorized to update this event",
        })
    }
    
    var updateData map[string]interface{}
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }
    
    if err := config.DB.Model(&event).Updates(updateData).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update event",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Event updated successfully",
        "event":   event,
    })
}

func VerifyEvent(c *fiber.Ctx) error {
    eventID := c.Params("id")
    
    var event models.Event
    if err := config.DB.First(&event, eventID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Event not found",
        })
    }
    
    var req struct {
        Status           string `json:"status"`
        ApprovalComment  string `json:"approval_comment,omitempty"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }
    
    event.Status = req.Status
    event.ApprovalComment = req.ApprovalComment
    
    if err := config.DB.Save(&event).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to verify event",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Event verification updated",
        "event":   event,
    })
}

func DeleteEvent(c *fiber.Ctx) error {
    user := c.Locals("user").(models.User)
    eventID := c.Params("id")
    
    var event models.Event
    if err := config.DB.First(&event, eventID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Event not found",
        })
    }
    
    // Check ownership or admin
    if event.OwnerID != user.UserID && user.Role != "admin" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Not authorized to delete this event",
        })
    }
    
    if err := config.DB.Delete(&event).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to delete event",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Event deleted successfully",
    })
}