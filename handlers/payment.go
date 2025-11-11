package handlers

import (
	"os"

	"github.com/Tsaniii18/Ticketing-Backend/config"
	"github.com/Tsaniii18/Ticketing-Backend/models"
	"github.com/Tsaniii18/Ticketing-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func PaymentMidtrans(c *fiber.Ctx) error {
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

	// Create Snap client
	var s snap.Client
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	var items []midtrans.ItemDetails
	for _, cartItem := range cart {
		item := midtrans.ItemDetails{
			ID:       cartItem.TicketCategoryID,
			Name:     cartItem.TicketCategory.Name,
			Price:    int64(cartItem.PriceTotal),
			Qty:      int32(cartItem.Quantity),
			Category: "Event Ticket",
		}
		items = append(items, item)
	}

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  utils.GenerateTransactionID(),
			GrossAmt: int64(total),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
		Items: &items,
	}

	resp, err := s.CreateTransaction(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction with Midtrans",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment Midtrans Endpoint",
		"total":   total,
		"url":     resp.RedirectURL,
		"token":   resp.Token,
	})

}

func PaymentNotificationHandler(c *fiber.Ctx) error {
	var notifPayload map[string]interface{}
	if err := c.BodyParser(&notifPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	orderID := notifPayload["order_id"].(string)
	status := notifPayload["transaction_status"].(string)

	// Update status transaksi di database kamu
	switch status {
	case "settlement":
		// pembayaran sukses
		// update ke DB: transaksi TRX-001 => "paid"
	case "pending":
		// menunggu pembayaran
	case "deny", "cancel":
		// pembayaran gagal
	}

	return c.JSON(fiber.Map{
		"message": "Notification processed",
		"orderID": orderID,
		"status":  status,
	})
}
