package handlers

import (
	"fmt"
	"jwt-go/db"
	"jwt-go/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AddItemToCart(c *fiber.Ctx) error {
	// Ambil userClaims dari context
	userClaims, ok := c.Locals("userClaims").(jwt.MapClaims)
	if !ok {
		fmt.Println(" userClaims tidak ditemukan!")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Debug: Print isi userClaims
	fmt.Println("userClaims di AddItemToCart:", userClaims)

	// Ambil user_id dengan aman
	userIDFloat, ok := userClaims["user_id"].(float64)
	if !ok {
		fmt.Println(" user_id tidak ditemukan di token!")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}
	userID := uint(userIDFloat) // Konversi dari float64 ke uint

	// Ambil permission dari claims
	userPermission, ok := userClaims["permission"].(string)
	if !ok {
		fmt.Println(" Permission tidak ditemukan di token!")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}

	// Debug: Print permission
	fmt.Println(" User Permission:", userPermission)

	// Validasi hanya admin yang bisa tambah item
	if userPermission != "admin" {
		fmt.Println(" Akses ditolak! Hanya admin yang bisa menambah item.")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
	}

	// Parsing body request
	var cartItem models.CartItem
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Set user ID di cart item
	cartItem.UserID = userID

	// Simpan ke database
	if err := db.DB.Create(&cartItem).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add item to cart"})
	}

	return c.JSON(cartItem)
}

func GetCartItem(c *fiber.Ctx) error {
	id := c.Params("id")

	var cartItem models.CartItem
	if err := db.DB.First(&cartItem, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}

	return c.JSON(cartItem)
}

func DeleteCartItem(c *fiber.Ctx) error {
	// Ambil userClaims dari context
	userClaims, ok := c.Locals("userClaims").(jwt.MapClaims)
	if !ok {
		fmt.Println("[ERROR] Tidak dapat membaca userClaims dari context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Ambil permission user dari context
	userPermission, ok := c.Locals("userPermission").(string)
	if !ok || userPermission != "admin" {
		fmt.Println("[ERROR] User tidak memiliki izin untuk menghapus item")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
	}

	// Ambil user_id dengan aman
	userIDFloat, ok := userClaims["user_id"].(float64)
	if !ok {
		fmt.Println("[ERROR] user_id tidak valid dalam token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}
	userID := uint(userIDFloat) // Konversi dari float64 ke uint

	// Ambil ID item dari parameter URL
	id := c.Params("id")
	if id == "" {
		fmt.Println("[ERROR] ID item tidak ditemukan dalam URL")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing item ID"})
	}

	// Cari item di database berdasarkan ID
	var cartItem models.CartItem
	if err := db.DB.First(&cartItem, id).Error; err != nil {
		fmt.Println("[ERROR] Item tidak ditemukan, ID:", id)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	// Pastikan hanya pemilik item atau admin yang bisa menghapus
	if cartItem.UserID != userID {
		fmt.Println("[ERROR] User tidak berhak menghapus item orang lain")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	// Hapus item dari database
	if err := db.DB.Delete(&cartItem).Error; err != nil {
		fmt.Println("[ERROR] Gagal menghapus item, ID:", id, "Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete item"})
	}

	fmt.Println("[SUCCESS] Item berhasil dihapus, ID:", id)
	return c.JSON(fiber.Map{"message": "Item deleted successfully"})
}
