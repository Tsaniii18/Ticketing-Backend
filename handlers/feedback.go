package handlers

import (
	// "time"

	// "log"
	"context"
	"fmt"

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
		Status:           "unseen",
		Comment:          comment,
		Image:            imageURL,
	}

	if err := config.DB.Create(&feedResp).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create feedback: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Tickets fetched successfully",
		"feed_respon": feedResp,
	})

}
