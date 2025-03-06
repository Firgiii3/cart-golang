package handlers

import (
	"fmt"
	"jwt-go/db"
	"jwt-go/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Hash Password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Register
// Register User atau Admin
func Register(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Hash password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error hashing password"})
	}
	user.Password = hashedPassword

	// Cek apakah ini admin (misalnya username "admin")
	if user.Username == "admin" {
		user.Permission = "admin" // Otomatis jadi admin
	} else {
		user.Permission = "user" // Default user biasa
	}

	// Simpan ke database
	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to register user"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// Login
func Login(c *fiber.Ctx) error {
	var input models.User
	var user models.User

	// Parsing body request
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Ambil user dari database berdasarkan username
	err := db.DB.Where("username = ?", input.Username).First(&user).Error
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Cek apakah password cocok
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Pastikan permission tidak kosong
	if user.Permission == "" {
		user.Permission = "user" // Default permission jika tidak ada
	}

	// Generate JWT dengan permission
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"permission": user.Permission,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error generating token"})
	}

	// Debugging: Print user dan token
	fmt.Println("User:", user.Username, "Permission:", user.Permission)
	fmt.Println("Generated Token:", t)

	// Berikan token ke client
	return c.JSON(fiber.Map{"token": t})
}

func GetUserFromToken(c *fiber.Ctx) error {
	// Ambil claims dari middleware
	claims, ok := c.Locals("userClaims").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Ambil user_id dari token
	userIDFloat, ok := claims["user_id"].(float64) // JWT mengembalikan angka sebagai float64
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}
	userID := uint(userIDFloat) // Konversi ke uint

	// Ambil user dari database
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Kembalikan data user (tanpa password)
	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
	})
}

// DeleteUser menghapus user berdasarkan ID yang diberikan
func DeleteUser(c *fiber.Ctx) error {
	// Pastikan token valid
	claims, ok := c.Locals("userClaims").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Ambil user ID dari token
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token data",
		})
	}

	// Cari user di database
	var user models.User
	if err := db.DB.First(&user, uint(userID)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	db.DB.Unscoped().Delete(&user)
	return c.JSON(user)
}

// Logout user
func Logout(c *fiber.Ctx) error {
	// Pastikan token valid
	_, exists := c.Locals("userClaims").(jwt.MapClaims)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Respon sukses (client harus hapus token)
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
