package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		fmt.Println("[ERROR] Tidak ada token dikirim!")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Hapus "Bearer " dari token jika ada
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Debug: Print token setelah "Bearer " dihapus
	fmt.Println("[DEBUG] Token After Bearer Removed:", tokenString)

	// Parsing token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Jika parsing error atau token invalid
	if err != nil || !token.Valid {
		fmt.Println("[ERROR] Token tidak valid!", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Debug: Print isi claims token
	fmt.Println("[DEBUG] User Claims dari Middleware:", claims)

	// Pastikan token mengandung user_id dan permission
	userID, idOk := claims["user_id"].(float64) // JWT menyimpan angka sebagai float64
	permission, permOk := claims["permission"].(string)

	if !idOk || !permOk {
		fmt.Println("[ERROR] Token tidak mengandung user_id atau permission yang valid")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token data",
		})
	}

	// Simpan claims ke context supaya bisa diakses di handler
	c.Locals("userClaims", claims)
	c.Locals("userID", uint(userID)) // Konversi float64 ke uint
	c.Locals("userPermission", permission)

	return c.Next()
}
