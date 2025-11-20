package handlers

import (
	"tonish/backend/database"
	"tonish/backend/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllTasks retrieves all tasks for a user
func GetAllTasks(c *fiber.Ctx) error {
	var tasks []models.Task
	
	// For now, get all tasks (will add user filtering with auth later)
	database.DB.Find(&tasks)
	
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
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	database.DB.Create(&task)
	
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
	
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	database.DB.Save(&task)
	
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
	
	database.DB.Delete(&task)
	
	return c.Status(204).SendString("")
}

// GetTasksByStatus retrieves tasks filtered by status
func GetTasksByStatus(c *fiber.Ctx) error {
	status := c.Query("status")
	var tasks []models.Task
	
	if status != "" {
		database.DB.Where("status = ?", status).Find(&tasks)
	} else {
		database.DB.Find(&tasks)
	}
	
	return c.JSON(tasks)
}

// GetTasksByQuadrant retrieves tasks filtered by Eisenhower Matrix quadrant
func GetTasksByQuadrant(c *fiber.Ctx) error {
	quadrant := c.Params("quadrant")
	var tasks []models.Task
	
	database.DB.Where("quadrant = ?", quadrant).Find(&tasks)
	
	return c.JSON(tasks)
}
