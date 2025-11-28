package ai

import (
	"encoding/json"
	"fmt"
	"strings"
	"tonish/backend/models"
)

// TOON - Token-Oriented Object Notation
// Converts structured data to AI-friendly format and back

// TaskToTOON converts a Task to TOON format for AI processing
func TaskToTOON(task models.Task) string {
	var sb strings.Builder
	
	sb.WriteString("TASK_OBJECT {\n")
	sb.WriteString(fmt.Sprintf("  title: \"%s\"\n", task.Title))
	sb.WriteString(fmt.Sprintf("  description: \"%s\"\n", task.Description))
	sb.WriteString(fmt.Sprintf("  priority: %s\n", task.Priority))
	sb.WriteString(fmt.Sprintf("  status: %s\n", task.Status))
	sb.WriteString(fmt.Sprintf("  quadrant: %s\n", task.Quadrant))
	
	if task.DueDate != nil {
		sb.WriteString(fmt.Sprintf("  due_date: \"%s\"\n", task.DueDate.Format("2006-01-02")))
	}
	
	if task.Tags != "" {
		sb.WriteString(fmt.Sprintf("  tags: %s\n", task.Tags))
	}
	
	sb.WriteString("}")
	
	return sb.String()
}

// NotebookToTOON converts a Notebook to TOON format for AI processing
func NotebookToTOON(notebook models.Notebook) string {
	var sb strings.Builder
	
	sb.WriteString("NOTEBOOK_OBJECT {\n")
	sb.WriteString(fmt.Sprintf("  name: \"%s\"\n", notebook.Name))
	sb.WriteString(fmt.Sprintf("  tags: %s\n", notebook.Tags))
	sb.WriteString(fmt.Sprintf("  is_pinned: %t\n", notebook.IsPinned))
	sb.WriteString(fmt.Sprintf("  page_count: %d\n", len(notebook.Pages)))
	
	if len(notebook.Pages) > 0 {
		sb.WriteString("  pages: [\n")
		for i, page := range notebook.Pages {
			sb.WriteString(fmt.Sprintf("    PAGE_%d {\n", i+1))
			sb.WriteString(fmt.Sprintf("      title: \"%s\"\n", page.Title))
			contentPreview := page.Content
			if len(contentPreview) > 200 {
				contentPreview = contentPreview[:200] + "..."
			}
			sb.WriteString(fmt.Sprintf("      content_preview: \"%s\"\n", contentPreview))
			sb.WriteString("    }\n")
		}
		sb.WriteString("  ]\n")
	}
	
	sb.WriteString("}")
	
	return sb.String()
}

// TOONPrompt creates a complete prompt for the AI with TOON data
func TOONPrompt(action string, toonData string, context map[string]interface{}) string {
	var sb strings.Builder
	
	sb.WriteString("You are Tonish AI Assistant, specialized in task management and productivity.\n")
	sb.WriteString(fmt.Sprintf("Action requested: %s\n\n", action))
	sb.WriteString("Input data in TOON format:\n")
	sb.WriteString(toonData)
	sb.WriteString("\n\n")
	
	if len(context) > 0 {
		sb.WriteString("Additional context:\n")
		for key, value := range context {
			sb.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
		sb.WriteString("\n")
	}
	
	sb.WriteString("Please respond with valid JSON only. No markdown, no explanation.\n")
	sb.WriteString("Response format: ")
	
	switch action {
	case "enhance_task":
		sb.WriteString(`{"title": "...", "description": "...", "priority": "...", "quadrant": "...", "tags": [...], "suggestions": "..."}`)
	case "suggest_breakdown":
		sb.WriteString(`{"subtasks": [{"title": "...", "description": "...", "priority": "..."}], "reasoning": "..."}`)
	case "analyze_notebook":
		sb.WriteString(`{"summary": "...", "topics": [...], "suggestions": [...], "key_insights": [...]}`)
	case "generate_page_ideas":
		sb.WriteString(`{"ideas": [{"title": "...", "description": "...", "relevance": "..."}]}`)
	default:
		sb.WriteString(`{"result": "...", "message": "..."}`)
	}
	
	return sb.String()
}

// ParseAIResponse parses the AI JSON response
func ParseAIResponse(response string) (map[string]interface{}, error) {
	// Clean up the response - remove markdown code blocks if present
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)
	
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}
	
	return result, nil
}
