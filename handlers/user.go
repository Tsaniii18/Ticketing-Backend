package handlers

import (
    "context"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/Tsaniii18/Ticketing-Backend/config"
    "github.com/Tsaniii18/Ticketing-Backend/models"
)

func GetProfile(c *fiber.Ctx) error {
    user := c.Locals("user").(models.User)
    return c.JSON(user)
}

func UpdateProfile(c *fiber.Ctx) error {
    user := c.Locals("user").(models.User)
    
    name := c.FormValue("name")
    email := c.FormValue("email")
    organization := c.FormValue("organization")
    organizationType := c.FormValue("organization_type")
    organizationDescription := c.FormValue("organization_description")
    ktp := c.FormValue("ktp")

    updateData := map[string]interface{}{
        "Name":                    name,
        "Email":                   email,
        "Organization":            organization,
        "OrganizationType":        organizationType,
        "OrganizationDescription": organizationDescription,
        "KTP":                     ktp,
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

    if err := config.DB.Model(&user).Updates(updateData).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update profile",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Profile updated successfully",
        "user":    user,
    })
}

func GetUsers(c *fiber.Ctx) error {
    var users []models.User
    if err := config.DB.Find(&users).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch users",
        })
    }
    
    return c.JSON(users)
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