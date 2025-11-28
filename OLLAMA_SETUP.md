# Tonish AI Integration - Using Existing Ollama

## 🎯 Overview

Your Tonish application is now configured to use your existing Ollama instance at **192.168.5.10:11434**.

## 📡 Ollama Connection Details

- **External URL**: `http://192.168.5.10:11434`
- **Local URL**: `http://localhost:11434` (same machine)
- **Base API**: `http://192.168.5.10:11434/api`

### Available Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/generate` | POST | Generate text completions |
| `/api/chat` | POST | Chat completions |
| `/api/tags` | GET | List available models |
| `/api/pull` | POST | Download models |

## 🚀 Quick Start

### 1. Test Your Ollama Connection

```bash
chmod +x test-ollama.sh
./test-ollama.sh
```

Expected output:
```
✓ Connection successful!
✓ AI generation works!
✓ Chat endpoint works!
```

### 2. Verify/Install Model

Check if qwen2.5-coder is available:
```bash
curl http://192.168.5.10:11434/api/tags
```

If not present, pull it:
```bash
curl -X POST http://192.168.5.10:11434/api/pull \
  -H "Content-Type: application/json" \
  -d '{"name": "qwen2.5-coder:latest"}'
```

Or run the automated setup:
```bash
chmod +x setup-ai.sh
./setup-ai.sh
```

### 3. Start Tonish

**Option A - Docker Compose (Production)**:
```bash
docker-compose up -d
```

**Option B - Development Mode**:
```bash
./start-dev.sh
```

### 4. Test AI Integration

```bash
curl http://192.168.5.10:8080/api/ai/health
```

Expected response:
```json
{
  "status": "healthy",
  "model": "qwen2.5-coder:latest",
  "url": "http://192.168.5.10:11434"
}
```

## 🏗️ Architecture

```
┌──────────────┐
│   Browser    │
└──────┬───────┘
       │ http://192.168.5.10:5173
       ▼
┌──────────────┐
│   Frontend   │ (SvelteKit)
│   :5173      │
└──────┬───────┘
       │ JSON API
       ▼
┌──────────────┐
│   Backend    │ (Go/Fiber)
│   :8080      │
└──────┬───────┘
       │ TOON Format
       ▼
┌──────────────────────────┐
│   Ollama (External)      │
│   192.168.5.10:11434     │
│   Model: qwen2.5-coder   │
└──────────────────────────┘
```

## 🔧 Configuration Files Updated

### docker-compose.yml
- **Removed**: Ollama container (using external)
- **Added**: `OLLAMA_URL=http://192.168.5.10:11434`
- **Added**: `extra_hosts` for Docker network access

### backend/main.go
- Initializes AI client on startup
- Connects to external Ollama
- Health check on boot

### backend/ai/ollama.go
- Environment variable: `OLLAMA_URL` (defaults to 192.168.5.10:11434)
- Environment variable: `OLLAMA_MODEL` (defaults to qwen2.5-coder:latest)

## 🧪 Testing

### Test Ollama Directly

**List models**:
```bash
curl http://192.168.5.10:11434/api/tags
```

**Test generation**:
```bash
curl -X POST http://192.168.5.10:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen2.5-coder:latest",
    "prompt": "Write a hello world in Go",
    "stream": false
  }'
```

**Test chat**:
```bash
curl -X POST http://192.168.5.10:11434/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen2.5-coder:latest",
    "messages": [
      {"role": "user", "content": "What is 2+2?"}
    ],
    "stream": false
  }'
```

### Test Tonish AI Endpoints

**Health check**:
```bash
curl http://192.168.5.10:8080/api/ai/health
```

**Enhance a task**:
```bash
curl -X POST http://192.168.5.10:8080/api/ai/enhance-task \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Fix bug",
    "description": "Something broke",
    "priority": "medium"
  }'
```

**Breakdown a task**:
```bash
curl -X POST http://192.168.5.10:8080/api/ai/breakdown-task \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Build new feature",
    "description": "Create user dashboard with analytics"
  }'
```

## 📋 Available Models

You can use any model available in your Ollama instance. To change the model, update the environment variable:

```bash
# In docker-compose.yml
environment:
  - OLLAMA_MODEL=codellama:latest  # or any other model

# Or in development
export OLLAMA_MODEL=mistral:latest
./start-dev.sh
```

**Recommended models for Tonish**:
- `qwen2.5-coder:latest` - Best for code & structured data (default)
- `qwen2.5-coder:7b` - Balanced performance
- `qwen2.5-coder:1.5b` - Fastest, lower resource usage
- `codellama:latest` - Good alternative for code tasks
- `mistral:latest` - General purpose

## 🛠️ Troubleshooting

### Backend can't connect to Ollama

**From host machine**:
```bash
curl http://192.168.5.10:11434/api/tags
```

**From Docker container**:
```bash
docker-compose exec backend sh -c "wget -qO- http://192.168.5.10:11434/api/tags"
```

If Docker can't reach it, ensure:
1. Ollama is listening on `0.0.0.0:11434` not just `127.0.0.1`
2. Firewall allows connections to port 11434
3. The IP `192.168.5.10` is correct for your network

### Slow responses

First request is always slower (model loading). To keep model in memory:

```bash
# Pre-warm the model
curl -X POST http://192.168.5.10:11434/api/generate \
  -d '{"model":"qwen2.5-coder:latest","prompt":"warm up","stream":false}'
```

### Model not found

Pull the model:
```bash
curl -X POST http://192.168.5.10:11434/api/pull \
  -d '{"name": "qwen2.5-coder:latest"}'
```

Check pull progress:
```bash
# If you have access to Ollama's Docker logs
docker logs -f <ollama-container-name>
```

## 📚 Documentation

- **AI Integration Guide**: `AI_INTEGRATION_GUIDE.md`
- **Configuration Summary**: `AI_CONFIGURATION_SUMMARY.md`
- **Ollama API Docs**: https://github.com/ollama/ollama/blob/main/docs/api.md

## 🔐 Security Notes

✅ **Your setup is secure**:
- Ollama runs locally on your network
- No external API calls
- No data leaves your infrastructure
- Models are open-source

⚠️ **For production**:
- Consider adding rate limiting to AI endpoints
- Add authentication if exposing publicly
- Monitor resource usage (Ollama can be CPU/memory intensive)

## 🎨 Using Magic Wand in Frontend

Import the component in your Svelte pages:

```svelte
<script>
  import MagicWand from '$lib/components/MagicWand.svelte';
  
  let currentTask = {
    id: 1,
    title: "Write documentation",
    description: "Need to document the new feature"
  };
  
  function handleUpdate(updatedTask) {
    currentTask = updatedTask;
    // Refresh your task list
  }
</script>

<MagicWand 
  task={currentTask} 
  onTaskUpdated={handleUpdate}
/>
```

## 📊 Performance

- **First AI request**: 5-15 seconds (model loading)
- **Subsequent requests**: 2-5 seconds
- **Memory**: Ollama handles this (check your setup)
- **CPU**: Depends on your Ollama configuration

## 🎯 Next Steps

1. ✅ Test Ollama connection: `./test-ollama.sh`
2. ✅ Ensure model is available: `curl http://192.168.5.10:11434/api/tags`
3. ✅ Start Tonish: `docker-compose up -d` or `./start-dev.sh`
4. ✅ Test AI health: `curl http://192.168.5.10:8080/api/ai/health`
5. ✅ Try Magic Wand in the UI
6. 🎨 Customize TOON prompts in `backend/ai/toon.go`
7. 🚀 Add more AI features as needed

---

**Status**: ✅ Configured for external Ollama at 192.168.5.10:11434  
**Model**: qwen2.5-coder:latest  
**Architecture**: Hybrid AI (Local Ollama)
