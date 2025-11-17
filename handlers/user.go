package handlers

import (
	"context"
	"fmt"

	"github.com/Tsaniii18/Ticketing-Backend/config"
	"github.com/Tsaniii18/Ticketing-Backend/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	return c.JSON(user)
}

func UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Hanya ambil field yang boleh diubah
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	organization := c.FormValue("organization")
	organizationType := c.FormValue("organization_type")
	organizationDescription := c.FormValue("organization_description")
	// KTP sengaja tidak diambil dari form karena tidak boleh diubah

	updateData := map[string]interface{}{
		"Name":  name,
		"Email": email,
	}

	// Handle password update jika ada
	if password != "" {
		hashedPassword, err := HashPassword(password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		updateData["Password"] = hashedPassword
	}

	// Handle profile picture upload
	profilePictFile, err := c.FormFile("profile_pict")
	if err == nil {
		file, err := profilePictFile.Open()
		if err == nil {
			defer file.Close()
			folder := fmt.Sprintf("ticketing-app/users/%s/profile", user.UserID)
			profilePictURL, err := config.UploadImage(context.Background(), file, folder)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to upload profile picture",
				})
			}
			updateData["ProfilePict"] = profilePictURL
		}
	}

	// Untuk organizer, update field organisasi hanya jika ada nilai
	if user.Role == "organizer" {
		if organization != "" {
			updateData["Organization"] = organization
		}
		if organizationType != "" {
			updateData["OrganizationType"] = organizationType
		}
		if organizationDescription != "" {
			updateData["OrganizationDescription"] = organizationDescription
		}
		// KTP tidak diupdate karena tidak boleh diubah
	}

	if err := config.DB.Model(&user).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	// Reload user data untuk mendapatkan data terbaru
	if err := config.DB.First(&user, "user_id = ?", user.UserID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch updated profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

func GetUsers(c *fiber.Ctx) error {
	role := c.Query("role")

	var users []models.User
	query := config.DB

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if err := query.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

func GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	var user models.User
	if err := config.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.JSON(user)
}

func VerifyUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	if err := config.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Only organizers need verification
	if user.Role != "organizer" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only organizer accounts can be verified",
		})
	}

	var req struct {
		Status  string `json:"status"`
		Comment string `json:"comment,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	user.RegisterStatus = req.Status
	user.RegisterComment = req.Comment

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Organizer verification updated",
		"user": fiber.Map{
			"user_id":          user.UserID,
			"username":         user.Username,
			"name":             user.Name,
			"email":            user.Email,
			"role":             user.Role,
			"organization":     user.Organization,
			"ktp":              user.KTP,
			"register_status":  user.RegisterStatus,
			"register_comment": user.RegisterComment,
		},
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}