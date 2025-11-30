package main

import (
	"fmt"
	"log"
	"os"

	"go-app-base/auth"
	"go-app-base/config"
	"go-app-base/models"
)

func main() {
	// データベース接続
	config.ConnectDatabase()
	
	// 既存データをクリア
	config.DB.Exec("DELETE FROM users")
	
	// モックユーザーデータ
	users := []struct {
		Username string
		Email    string
		Password string
		IsActive bool
	}{
		{"admin", "admin@example.com", "password123", true},
		{"john_doe", "john@example.com", "password123", true},
		{"jane_smith", "jane@example.com", "password123", true},
		{"bob_wilson", "bob@example.com", "password123", true},
		{"alice_brown", "alice@example.com", "password123", true},
		{"test_user", "test@example.com", "password123", false},
	}

	for _, u := range users {
		hashedPassword, err := auth.HashPassword(u.Password)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", u.Email, err)
			continue
		}

		user := models.User{
			Username: u.Username,
			Email:    u.Email,
			Password: hashedPassword,
			IsActive: u.IsActive,
		}

		if err := config.DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", u.Email, err)
		} else {
			fmt.Printf("✓ Created user: %s (%s)\n", u.Username, u.Email)
		}
	}

	fmt.Println("\n✅ Mock data seeding completed!")
	os.Exit(0)
}
