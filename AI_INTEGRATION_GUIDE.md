# Tonish AI Integration - Complete Guide

## Architecture Overview

Tonish uses a **Hybrid AI** architecture with local Ollama integration:

```
┌─────────────┐      JSON Request       ┌──────────────┐
│  Frontend   │─────────────────────────▶│   Backend    │
│  (Svelte)   │                          │     (Go)     │
└─────────────┘                          └──────────────┘
      ▲                                          │
      │                                          ▼
      │                                   ┌──────────────┐
      │         JSON Response             │ TOON Converter│
      │◀──────────────────────────────────│  (ai/toon.go) │
                                          └──────────────┘
                                                 │
                                                 ▼
                                          ┌──────────────┐
                                          │    Ollama    │
                                          │ Qwen2.5-Coder│
                                          └──────────────┘
                                                 │
                                                 ▼
                                          ┌──────────────┐
                                          │   SQLite     │
                                          │   Database   │
                                          └──────────────┘
```

## Components

### 1. **TOON (Token-Oriented Object Notation)** - `backend/ai/toon.go`
Converts structured data (Tasks, Notebooks) into AI-friendly format:

```go
TASK_OBJECT {
  title: "Complete project documentation"
  description: "Write comprehensive docs"
  priority: high
  status: todo
  quadrant: urgent-important
}
```

### 2. **Ollama Client** - `backend/ai/ollama.go`
Handles communication with the Ollama API:
- Sends TOON-formatted prompts
- Receives and parses JSON responses
- Health checks and error handling

### 3. **AI Handlers** - `backend/handlers/ai.go`
REST API endpoints for AI features:
- `POST /api/ai/enhance-task` - Improve task details
- `POST /api/ai/breakdown-task` - Split into subtasks
- `GET /api/ai/analyze-notebook/:id` - Notebook insights
- `GET /api/ai/generate-page-ideas/:id` - Page suggestions
- `GET /api/ai/health` - Service status

### 4. **Frontend API** - `frontend/src/lib/api.ts`
TypeScript interface for AI calls from Svelte components

## Configuration

### Network & Ports
All services bind to `192.168.5.10`:
- **Frontend**: `192.168.5.10:5173`
- **Backend**: `192.168.5.10:8080`
- **Ollama**: `192.168.5.10:11434`

### Environment Variables

**Backend (`docker-compose.yml`)**:
```yaml
environment:
  - PORT=8080
  - OLLAMA_URL=http://ollama:11434
  - OLLAMA_MODEL=qwen2.5-coder:latest
```

**Development (`start-dev.sh`)**:
```bash
OLLAMA_URL=http://192.168.5.10:11434
PORT=8080
```

## Setup Instructions

### Prerequisites

**Ollama is already running** at `http://192.168.5.10:11434`

Verify it's accessible:
```bash
curl http://192.168.5.10:11434/api/tags
```

### Option 1: Docker Compose (Recommended)

1. **Start Tonish services** (Backend + Frontend):
```bash
docker-compose up -d
```

2. **Ensure AI model is available**:
```bash
# Check available models
curl http://192.168.5.10:11434/api/tags

# Pull qwen2.5-coder if not present
curl -X POST http://192.168.5.10:11434/api/pull \
  -d '{"name": "qwen2.5-coder:latest"}'
```

3. **Verify AI is working**:
```bash
curl http://192.168.5.10:8080/api/ai/health
```

### Option 2: Development Mode

1. **Run the test script**:
```bash
chmod +x test-ollama.sh
./test-ollama.sh
```

2. **Start development servers**:
```bash
chmod +x start-dev.sh
./start-dev.sh
```

### Option 3: Quick Setup (Automated)

```bash
chmod +x setup-ai.sh
./setup-ai.sh
```

This will:
- Check Ollama connection
- Verify/pull qwen2.5-coder model
- Test AI functionality

## API Usage Examples

### 1. Enhance a Task (Magic Wand)

**Frontend**:
```typescript
import { aiAPI } from '$lib/api';

async function enhanceTask(task) {
  try {
    const result = await aiAPI.enhanceTask(task);
    console.log('Enhanced task:', result.enhanced);
  } catch (error) {
    console.error('AI enhancement failed:', error);
  }
}
```

**Request**:
```bash
curl -X POST http://192.168.5.10:8080/api/ai/enhance-task \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Fix bug",
    "description": "Something is broken",
    "priority": "medium"
  }'
```

**Response**:
```json
{
  "enhanced": {
    "title": "Fix critical authentication bug in user login",
    "description": "Investigate and resolve authentication failure affecting user login flow. Check token validation and session management.",
    "priority": "high",
    "quadrant": "urgent-important",
    "tags": ["bug", "authentication", "security"],
    "suggestions": "Consider adding error logging and user feedback"
  },
  "original": { ... }
}
```

### 2. Break Down Complex Task

**Frontend**:
```typescript
const breakdown = await aiAPI.breakdownTask({
  title: "Launch new feature",
  description: "Deploy the new analytics dashboard"
});

// Create subtasks automatically
breakdown.subtasks.forEach(subtask => {
  taskAPI.create(subtask);
});
```

### 3. Analyze Notebook

**Frontend**:
```typescript
const analysis = await aiAPI.analyzeNotebook(notebookId);
console.log('Topics:', analysis.topics);
console.log('Insights:', analysis.key_insights);
```

### 4. Generate Page Ideas

**Frontend**:
```typescript
const ideas = await aiAPI.generatePageIdeas(notebookId);
ideas.ideas.forEach(idea => {
  console.log(`${idea.title}: ${idea.description}`);
});
```

## Data Flow

1. **User clicks "Magic Wand"** in Svelte component
2. **Frontend sends JSON** to Go backend via `aiAPI`
3. **Backend converts to TOON** format (`ai.TaskToTOON`)
4. **TOON + prompt sent to Ollama** (`ai.Client.Generate`)
5. **Ollama processes with Qwen 2.5-Coder** model
6. **Returns clean JSON** response
7. **Backend parses & validates** (`ai.ParseAIResponse`)
8. **Updates database** if needed
9. **Returns result to frontend** for display

## Model Information

**Qwen 2.5-Coder**:
- Specialized for code and technical tasks
- Excellent at structured data processing
- Fast inference on CPU
- ~7B parameters (lightweight)

**Alternative models**:
```bash
# Smaller, faster
docker exec ollama ollama pull qwen2.5-coder:1.5b

# Larger, more capable
docker exec ollama ollama pull qwen2.5-coder:14b
```

## Troubleshooting

### AI service unavailable
```bash
# Check Ollama is running
curl http://192.168.5.10:11434/api/tags

# Check from localhost (if on same machine)
curl http://localhost:11434/api/tags

# Test specific endpoints
curl http://192.168.5.10:11434/api/generate \
  -d '{"model":"qwen2.5-coder:latest","prompt":"test","stream":false}'
```

### Model not found
```bash
# List available models
curl http://192.168.5.10:11434/api/tags

# Pull the model via API
curl -X POST http://192.168.5.10:11434/api/pull \
  -d '{"name": "qwen2.5-coder:latest"}'

# Or if you have ollama CLI access
ollama pull qwen2.5-coder:latest
```

### Connection refused
Check network configuration:
```bash
# Verify Ollama is listening on correct interface
netstat -tuln | grep 11434

# Test from Tonish backend container
docker-compose exec backend sh -c "wget -qO- http://192.168.5.10:11434/api/tags"
```

## Performance Notes

- First request may be slow (model loading)
- Subsequent requests are faster (model cached)
- Response time: 2-10 seconds depending on complexity
- Consider adding loading states in UI
- Backend has 120s timeout for AI requests

## Security Considerations

- Ollama runs locally, no external API calls
- No data leaves your network
- Models are open-source
- Consider rate limiting for production
- Add authentication to AI endpoints if needed

## Future Enhancements

- [ ] Streaming responses for real-time feedback
- [ ] Fine-tuned model for Tonish-specific tasks
- [ ] Multi-turn conversations for complex queries
- [ ] AI-powered search across notebooks
- [ ] Automated task categorization
- [ ] Smart scheduling suggestions
