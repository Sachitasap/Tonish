package handlers

import (
	"tonish/backend/database"
	"tonish/backend/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllNotebooks retrieves all notebooks
func GetAllNotebooks(c *fiber.Ctx) error {
	var notebooks []models.Notebook
	
	database.DB.Preload("Pages").Find(&notebooks)
	
	return c.JSON(notebooks)
}

// GetNotebook retrieves a single notebook by ID with its pages
func GetNotebook(c *fiber.Ctx) error {
	id := c.Params("id")
	var notebook models.Notebook
	
	if err := database.DB.Preload("Pages").First(&notebook, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Notebook not found",
		})
	}
	
	return c.JSON(notebook)
}

// CreateNotebook creates a new notebook
func CreateNotebook(c *fiber.Ctx) error {
	notebook := new(models.Notebook)
	
	if err := c.BodyParser(notebook); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	database.DB.Create(&notebook)
	
	return c.Status(201).JSON(notebook)
}

// UpdateNotebook updates an existing notebook
func UpdateNotebook(c *fiber.Ctx) error {
	id := c.Params("id")
	var notebook models.Notebook
	
	if err := database.DB.First(&notebook, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Notebook not found",
		})
	}
	
	if err := c.BodyParser(&notebook); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	database.DB.Save(&notebook)
	
	return c.JSON(notebook)
}

// DeleteNotebook deletes a notebook and its pages
func DeleteNotebook(c *fiber.Ctx) error {
	id := c.Params("id")
	var notebook models.Notebook
	
	if err := database.DB.First(&notebook, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Notebook not found",
		})
	}
	
	// Delete all pages first
	database.DB.Where("notebook_id = ?", id).Delete(&models.Page{})
	
	// Delete the notebook
	database.DB.Delete(&notebook)
	
	return c.Status(204).SendString("")
}

// GetPage retrieves a single page by ID
func GetPage(c *fiber.Ctx) error {
	id := c.Params("id")
	var page models.Page
	
	if err := database.DB.First(&page, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Page not found",
		})
	}
	
	return c.JSON(page)
}

// CreatePage creates a new page in a notebook
func CreatePage(c *fiber.Ctx) error {
	page := new(models.Page)
	
	if err := c.BodyParser(page); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	database.DB.Create(&page)
	
	return c.Status(201).JSON(page)
}

// UpdatePage updates an existing page
func UpdatePage(c *fiber.Ctx) error {
	id := c.Params("id")
	var page models.Page
	
	if err := database.DB.First(&page, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Page not found",
		})
	}
	
	if err := c.BodyParser(&page); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	database.DB.Save(&page)
	
	return c.JSON(page)
}

// DeletePage deletes a page
func DeletePage(c *fiber.Ctx) error {
	id := c.Params("id")
	var page models.Page
	
	if err := database.DB.First(&page, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Page not found",
		})
	}
	
	database.DB.Delete(&page)
	
	return c.Status(204).SendString("")
}

// SearchPages searches for pages by title or content
func SearchPages(c *fiber.Ctx) error {
	query := c.Query("q")
	var pages []models.Page
	
	if query != "" {
		database.DB.Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%").Find(&pages)
	} else {
		database.DB.Find(&pages)
	}
	
	return c.JSON(pages)
}
