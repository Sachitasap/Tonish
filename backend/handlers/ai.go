package handlers

import (
	"fmt"
	"tonish/backend/ai"
	"tonish/backend/database"
	"tonish/backend/models"

	"github.com/gofiber/fiber/v2"
)

// EnhanceTask uses AI to enhance a task with better descriptions, priorities, and quadrant placement
func EnhanceTask(c *fiber.Ctx) error {
	var task models.Task
	
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	// Convert task to TOON format
	toonData := ai.TaskToTOON(task)
	
	// Create prompt for AI
	context := map[string]interface{}{
		"request": "Analyze this task and suggest improvements for better productivity",
	}
	prompt := ai.TOONPrompt("enhance_task", toonData, context)
	
	// Call Ollama
	response, err := ai.Client.Generate(prompt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("AI processing failed: %v", err),
		})
	}
	
	// Parse AI response
	result, err := ai.ParseAIResponse(response)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse AI response",
			"raw":   response,
		})
	}
	
	return c.JSON(fiber.Map{
		"enhanced": result,
		"original": task,
	})
}

// BreakdownTask uses AI to break down a complex task into subtasks
func BreakdownTask(c *fiber.Ctx) error {
	var task models.Task
	
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	// Convert task to TOON format
	toonData := ai.TaskToTOON(task)
	
	// Create prompt for AI
	context := map[string]interface{}{
		"request": "Break down this task into smaller, actionable subtasks",
	}
	prompt := ai.TOONPrompt("suggest_breakdown", toonData, context)
	
	// Call Ollama
	response, err := ai.Client.Generate(prompt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("AI processing failed: %v", err),
		})
	}
	
	// Parse AI response
	result, err := ai.ParseAIResponse(response)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse AI response",
			"raw":   response,
		})
	}
	
	return c.JSON(result)
}

// AnalyzeNotebook uses AI to analyze a notebook and provide insights
func AnalyzeNotebook(c *fiber.Ctx) error {
	id := c.Params("id")
	var notebook models.Notebook
	
	if err := database.DB.Preload("Pages").First(&notebook, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Notebook not found",
		})
	}
	
	// Convert notebook to TOON format
	toonData := ai.NotebookToTOON(notebook)
	
	// Create prompt for AI
	context := map[string]interface{}{
		"request": "Analyze this notebook and provide insights, topics, and suggestions",
	}
	prompt := ai.TOONPrompt("analyze_notebook", toonData, context)
	
	// Call Ollama
	response, err := ai.Client.Generate(prompt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("AI processing failed: %v", err),
		})
	}
	
	// Parse AI response
	result, err := ai.ParseAIResponse(response)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse AI response",
			"raw":   response,
		})
	}
	
	return c.JSON(result)
}

// GeneratePageIdeas uses AI to suggest new page ideas for a notebook
func GeneratePageIdeas(c *fiber.Ctx) error {
	id := c.Params("id")
	var notebook models.Notebook
	
	if err := database.DB.Preload("Pages").First(&notebook, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Notebook not found",
		})
	}
	
	// Convert notebook to TOON format
	toonData := ai.NotebookToTOON(notebook)
	
	// Create prompt for AI
	context := map[string]interface{}{
		"request": "Based on the existing pages, suggest new page ideas that would complement this notebook",
	}
	prompt := ai.TOONPrompt("generate_page_ideas", toonData, context)
	
	// Call Ollama
	response, err := ai.Client.Generate(prompt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("AI processing failed: %v", err),
		})
	}
	
	// Parse AI response
	result, err := ai.ParseAIResponse(response)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse AI response",
			"raw":   response,
		})
	}
	
	return c.JSON(result)
}

// AIHealthCheck checks if the AI service is running
func AIHealthCheck(c *fiber.Ctx) error {
	if err := ai.Client.HealthCheck(); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "unavailable",
			"error":  err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"status": "healthy",
		"model":  ai.Client.Model,
		"url":    ai.Client.BaseURL,
	})
}
