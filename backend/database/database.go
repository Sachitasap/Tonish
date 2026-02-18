package database

import (
	"errors"
	"log"
	"os"
	"tonish/backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "tonish.db" // Default fallback
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Notebook{},
		&models.Page{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")
}

// SeedDefaultUser ensures a default user exists when credentials are supplied.
func SeedDefaultUser() {
	defaultEmail := os.Getenv("DEFAULT_USER_EMAIL")
	defaultPassword := os.Getenv("DEFAULT_USER_PASSWORD")
	defaultName := os.Getenv("DEFAULT_USER_NAME")

	if defaultEmail == "" || defaultPassword == "" {
		log.Println("Default user credentials not provided; skipping seed")
		return
	}
	if defaultName == "" {
		defaultName = "Default User"
	}

	var user models.User
	if err := DB.Where("email = ?", defaultEmail).First(&user).Error; err == nil {
		log.Printf("Default user already exists: %s\n", defaultEmail)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Failed to check default user: %v\n", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash default user password: %v\n", err)
		return
	}

	user = models.User{
		Email:    defaultEmail,
		Password: string(hashedPassword),
		Name:     defaultName,
	}

	if err := DB.Create(&user).Error; err != nil {
		log.Printf("Failed to create default user: %v\n", err)
		return
	}

	log.Printf("Default user created successfully: %s\n", defaultEmail)
}

// NormalizeTaskTypes backfills task_type for existing records so Eisenhower
// matrix tasks stay isolated from Kanban board items.
func NormalizeTaskTypes() {
	if DB == nil {
		return
	}

	if err := DB.Model(&models.Task{}).
		Where("quadrant IS NOT NULL AND quadrant != ''").
		Where("task_type IS NULL OR task_type = '' OR task_type != ?", "matrix").
		Update("task_type", "matrix").Error; err != nil {
		log.Printf("Failed to backfill matrix task types: %v\n", err)
	}

	if err := DB.Model(&models.Task{}).
		Where("(task_type = '' OR task_type IS NULL) AND (quadrant IS NULL OR quadrant = '')").
		Update("task_type", "kanban").Error; err != nil {
		log.Printf("Failed to backfill kanban task types: %v\n", err)
	}

	if err := DB.Model(&models.Task{}).
		Where("status = ? AND completed_at IS NULL", "done").
		Update("completed_at", gorm.Expr("COALESCE(updated_at, CURRENT_TIMESTAMP)")); err != nil {
		log.Printf("Failed to backfill completed task timestamps: %v\n", err)
	}
}
