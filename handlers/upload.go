package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Tsaniii18/Ticketing-Backend/config"
	"github.com/Tsaniii18/Ticketing-Backend/models"
	"github.com/gofiber/fiber/v2"
)

type UploadResponse struct {
	URL string `json:"url"`
}

func UploadImage(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	
	// Get uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image file provided",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size too large. Maximum size is 5MB",
		})
	}

	// Validate file type
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Allowed types: JPG, JPEG, PNG, GIF, WEBP",
		})
	}

	// Open file
	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer fileHeader.Close()

	// Upload to Cloudinary
	folder := fmt.Sprintf("ticketing-app/users/%s", user.UserID)
	imageURL, err := config.UploadImage(context.Background(), fileHeader, folder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload image to Cloudinary",
		})
	}

	return c.JSON(UploadResponse{
		URL: imageURL,
	})
}

func UploadMultipleImages(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No images provided",
		})
	}

	var uploadedURLs []string

	for _, file := range files {
		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			continue // Skip large files
		}

		// Validate file type
		allowedExtensions := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".webp": true,
		}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExtensions[ext] {
			continue // Skip invalid file types
		}

		// Open file
		fileHeader, err := file.Open()
		if err != nil {
			continue // Skip files that can't be opened
		}

		// Upload to Cloudinary
		folder := fmt.Sprintf("ticketing-app/users/%s", user.UserID)
		imageURL, err := config.UploadImage(context.Background(), fileHeader, folder)
		fileHeader.Close()

		if err == nil {
			uploadedURLs = append(uploadedURLs, imageURL)
		}
	}

	if len(uploadedURLs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid images were uploaded",
		})
	}

	return c.JSON(fiber.Map{
		"urls": uploadedURLs,
	})
}