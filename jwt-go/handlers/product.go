package handlers

import (
	"jwt-go/db"
	"jwt-go/models"

	"github.com/gofiber/fiber/v2"
)

// Menambahkan produk
func AddProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := db.DB.Create(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add product"})
	}

	return c.JSON(product)
}

// Mengedit produk
func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := db.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	db.DB.Save(&product)
	return c.JSON(product)
}

// Menghapus produk
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := db.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	db.DB.Delete(&product)
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}
