package db

import (
	"fmt"
	"jwt-go/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Ambil koneksi database dari environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=carts port=5432 sslmode=disable"
	}

	// Koneksi ke PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Simpan koneksi ke variabel global
	DB = db

	// Migrasi tabel
	fmt.Println("Running migrations...")
	err = DB.AutoMigrate(&models.User{}, &models.CartItem{}, &models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Migrations completed!")

	// Buat admin default jika belum ada
	createDefaultAdmin()
}

// Buat admin default jika belum ada
func createDefaultAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count)

	if count == 0 {
		fmt.Println("Creating default admin...")

		// Hash password admin dengan fungsi baru
		hashedPassword, _ := models.HashPassword("admin123")

		admin := models.User{
			Username:   "admin",
			Password:   hashedPassword,
			Permission: "admin",
		}

		// Simpan admin ke database
		DB.Create(&admin)
		fmt.Println("Default admin created: username=admin, password=admin123")
	}
}
