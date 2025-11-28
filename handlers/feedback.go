package handlers

import (
	// "time"

	"context"
	"fmt"
	"log"

	"github.com/Tsaniii18/Ticketing-Backend/config"
	"github.com/Tsaniii18/Ticketing-Backend/models"
	"github.com/Tsaniii18/Ticketing-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateFeedback(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	feedback_category := c.FormValue("feedback_category")
	// status := c.FormValue("status")
	comment := c.FormValue("comment")

	var imageURL string

	imageFile, err := c.FormFile("image")
	if err == nil {
		file, err := imageFile.Open()
		if err == nil {
			defer file.Close()
			folder := fmt.Sprintf("ticketing-app/feedback/%s/images", user.UserID)
			imageURL, err = config.UploadImage(context.Background(), file, folder)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to upload event image",
				})
			}
		}
	}

	feedResp := models.Feedback{
		FeedbackID:       utils.GenerateFeedID(),
		OwnerID:          user.UserID,
		FeedbackCategory: feedback_category,
		Status:           "waiting",
		Comment:          comment,
		Image:            imageURL,
	}

	if err := config.DB.Create(&feedResp).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create feedback: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Feedback created successfully",
		"feed_respon": feedResp,
	})

}

func GetAllFeedbacks(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	if user.Role != "admin" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only admin accounts permited to acccess",
		})
	}

	var feedbacks []models.Feedback
	if err := config.DB.Model(&feedbacks).Preload("User").Find(&feedbacks).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Feedback Not Fount",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Success fetch all feedback",
		"feedback": feedbacks,
	})

}

func GetMyFeedbacks(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var feedbacks []models.Feedback
	if err := config.DB.Model(&feedbacks).Preload("User").Where("owner_id = ?", user.UserID).Find(&feedbacks).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Feedback Not Fount",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Success fetch all feedback",
		"feedback": feedbacks,
	})
}

func GetFeedback(c *fiber.Ctx) error {
	feedbackID := c.Params("id")
	// user := c.Locals("user").(models.User)

	// if user.Role != "admin" {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Only admin accounts permited to acccess",
	// 	})
	// }

	var feedbackFigure models.Feedback
	if err := config.DB.Model(&feedbackFigure).Preload("User").Where("feedback_id = ?", feedbackID).First(&feedbackFigure).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Feedback Not Fount",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Success fetch the feedback",
		"feedback": feedbackFigure,
	})

}

func UpdateStatusFeedback(c *fiber.Ctx) error {

	var updateStatus map[string]interface{}
	user := c.Locals("user").(models.User)
	idFeedback := c.Params("id")

	if user.Role != "admin" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only admin accounts permited to acccess",
		})
	}

	if err := c.BodyParser(&updateStatus); err != nil {
		log.Printf("Error parsing update status: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	status, ok := updateStatus["status"].(string)

	if !ok {
		log.Printf("Invalid status request")
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status"})
	}

	if err := config.DB.Model(&models.Feedback{}).
		Where("feedback_id = ?", idFeedback).
		Update("status", status).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Feedback Not Fount",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success update status to " + status,
		"reply":   "No good",
	})

}
