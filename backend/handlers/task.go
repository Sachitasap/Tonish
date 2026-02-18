package handlers

import (
	"time"

	"tonish/backend/database"
	"tonish/backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func applyTaskTypeDefaults(task *models.Task) {
	if task == nil {
		return
	}

	if task.Quadrant != "" {
		task.TaskType = "matrix"
	} else if task.TaskType == "" {
		task.TaskType = "kanban"
	}
}

// GetAllTasks retrieves all tasks for a user
func GetAllTasks(c *fiber.Ctx) error {
	var tasks []models.Task

	// For now, get all tasks (will add user filtering with auth later)
	database.DB.Where("is_archived = ?", false).Find(&tasks)

	return c.JSON(tasks)
}

// GetTask retrieves a single task by ID
func GetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task

	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	return c.JSON(task)
}

// CreateTask creates a new task
func CreateTask(c *fiber.Ctx) error {
	task := new(models.Task)

	if err := c.BodyParser(task); err != nil {
		// Log the actual error for debugging
		println("Body parser error:", err.Error())
		println("Request body:", string(c.Body()))
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Get user_id from auth middleware context
	userID := c.Locals("user_id")
	if userID != nil {
		task.UserID = userID.(uint)
	}

	applyTaskTypeDefaults(task)
	setCompletionTimestamp(task, false)

	if err := database.DB.Create(&task).Error; err != nil {
		println("Database error:", err.Error())
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create task: " + err.Error()})
	}

	return c.Status(201).JSON(task)
}

// UpdateTask updates an existing task
func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task

	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	previousStatus := task.Status
	preservedUserID := task.UserID  // Preserve the original user_id

	if err := c.BodyParser(&task); err != nil {
		println("Update body parser error:", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Restore the user_id if it was cleared by body parser
	if task.UserID == 0 {
		task.UserID = preservedUserID
	}

	applyTaskTypeDefaults(&task)
	setCompletionTimestamp(&task, previousStatus == "done")

	if err := database.DB.Save(&task).Error; err != nil {
		println("Update database error:", err.Error())
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task: " + err.Error()})
	}

	return c.JSON(task)
}

// DeleteTask deletes a task
func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task

	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	if err := database.DB.Delete(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	return c.Status(204).SendString("")
}

// GetTasksByStatus retrieves tasks filtered by status
func GetTasksByStatus(c *fiber.Ctx) error {
	status := c.Query("status")
	var tasks []models.Task

	query := database.DB.Where("is_archived = ?", false)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Find(&tasks)

	return c.JSON(tasks)
}

// GetTasksByQuadrant retrieves tasks filtered by Eisenhower Matrix quadrant
func GetTasksByQuadrant(c *fiber.Ctx) error {
	quadrant := c.Params("quadrant")
	var tasks []models.Task

	database.DB.Where("quadrant = ? AND is_archived = ?", quadrant, false).Find(&tasks)

	return c.JSON(tasks)
}

func GetArchivedTasks(c *fiber.Ctx) error {
	var tasks []models.Task

	if err := database.DB.Unscoped().
		Where("is_archived = ?", true).
		Or("completed_at IS NOT NULL").
		Or("deleted_at IS NOT NULL").
		Order("COALESCE(deleted_at, completed_at, updated_at) DESC").
		Find(&tasks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load archived tasks"})
	}

	return c.JSON(tasks)
}

func ArchiveTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task

	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	task.IsArchived = true

	if err := database.DB.Save(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to archive task"})
	}

	return c.JSON(task)
}

func RestoreTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task

	if err := database.DB.Unscoped().First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	task.IsArchived = false
	if task.DeletedAt.Valid {
		task.DeletedAt = gorm.DeletedAt{}
	}
	if task.CompletedAt != nil {
		task.CompletedAt = nil
		task.Status = "todo"
	}

	if err := database.DB.Unscoped().Save(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to restore task"})
	}

	return c.JSON(task)
}

func PermanentDeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := database.DB.Unscoped().Delete(&models.Task{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to permanently delete task"})
	}

	return c.Status(204).SendString("")
}

func setCompletionTimestamp(task *models.Task, previouslyCompleted bool) {
	if task == nil {
		return
	}

	if task.Status == "done" {
		if task.CompletedAt == nil {
			now := time.Now()
			task.CompletedAt = &now
		}
		return
	}

	if previouslyCompleted {
		task.CompletedAt = nil
	}
}
