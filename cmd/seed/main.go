package main

import (
	"flag"
	"log"
	"time"

	"github.com/huy/quizme-backend/internal/config"
	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Define flags
	username := flag.String("username", "admin", "Admin username")
	email := flag.String("email", "admin@quizme.com", "Admin email")
	password := flag.String("password", "admin123", "Admin password")
	fullName := flag.String("fullname", "Administrator", "Admin full name")
	updateExisting := flag.Bool("update", false, "Update existing admin user if it already exists")

	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create user repository
	userRepo := repository.NewUserRepository(db)
	userProfileRepo := repository.NewUserProfileRepository(db)

	// Check if user already exists
	existingUser, err := userRepo.FindByUsername(*username)
	if err == nil && existingUser != nil {
		if !*updateExisting {
			log.Printf("Admin user '%s' already exists. Use --update flag to update it.", *username)
			return
		}
		log.Printf("Updating existing admin user: %s", *username)
		updateAdminUser(db, userRepo, existingUser, *password, *email, *fullName)
		return
	}

	// Check if email already exists
	existingUserByEmail, err := userRepo.FindByEmail(*email)
	if err == nil && existingUserByEmail != nil {
		log.Fatalf("Email '%s' is already registered", *email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create admin user
	adminUser := &domain.User{
		Username: *username,
		Email:    *email,
		Password: string(hashedPassword),
		FullName: *fullName,
		Role:     enums.RoleAdmin,
		IsActive: true,
	}

	if err := userRepo.Create(adminUser); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Printf("✓ Admin user created successfully!")
	log.Printf("  Username: %s", adminUser.Username)
	log.Printf("  Email: %s", adminUser.Email)
	log.Printf("  Full Name: %s", adminUser.FullName)
	log.Printf("  Role: %s", adminUser.Role)

	// Create user profile
	profile := &domain.UserProfile{
		UserID: adminUser.ID,
	}
	if err := userProfileRepo.Create(profile); err != nil {
		log.Fatalf("Failed to create user profile: %v", err)
	}

	log.Printf("✓ User profile created successfully!")
}

// updateAdminUser updates an existing admin user's password and other details
func updateAdminUser(db *gorm.DB, userRepo repository.UserRepository, user *domain.User, newPassword, newEmail, newFullName string) {
	// Hash new password if provided
	if newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}
		user.Password = string(hashedPassword)
		log.Printf("  Password updated")
	}

	// Update email if different
	if newEmail != "" && newEmail != user.Email {
		// Check if new email is already used
		existingUser, err := userRepo.FindByEmail(newEmail)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			log.Fatalf("Email '%s' is already registered", newEmail)
		}
		user.Email = newEmail
		log.Printf("  Email updated to: %s", newEmail)
	}

	// Update full name if different
	if newFullName != "" && newFullName != user.FullName {
		user.FullName = newFullName
		log.Printf("  Full name updated to: %s", newFullName)
	}

	// Ensure user is admin and active
	if user.Role != enums.RoleAdmin {
		user.Role = enums.RoleAdmin
		log.Printf("  Role set to: ADMIN")
	}
	if !user.IsActive {
		user.IsActive = true
		log.Printf("  User activated")
	}

	user.UpdatedAt = time.Now()

	if err := userRepo.Update(user); err != nil {
		log.Fatalf("Failed to update admin user: %v", err)
	}

	log.Printf("✓ Admin user updated successfully!")
	log.Printf("  Username: %s", user.Username)
	log.Printf("  Email: %s", user.Email)
	log.Printf("  Full Name: %s", user.FullName)
	log.Printf("  Role: %s", user.Role)
}
