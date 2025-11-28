# Tonish Environment Configuration Guide

## Environment Files Created

### 1. Root `.env` file
Main configuration file for the entire application.

**Location**: `/home/kubuntu1/Tonish/.env`

**Variables**:
- `HOST_IP` - Network IP address (192.168.5.10)
- `BACKEND_PORT` - Backend server port (8080)
- `FRONTEND_PORT` - Frontend server port (5173)
- `OLLAMA_URL` - Ollama AI service URL
- `OLLAMA_MODEL` - AI model name (qwen2.5-coder:3b)

### 2. Backend `.env` file
Backend-specific configuration.

**Location**: `/home/kubuntu1/Tonish/backend/.env`

**Variables**:
- `PORT` - Backend port
- `OLLAMA_URL` - AI service connection
- `OLLAMA_MODEL` - AI model to use
- `DATABASE_PATH` - SQLite database file path

### 3. Frontend `.env` file
Frontend-specific configuration.

**Location**: `/home/kubuntu1/Tonish/frontend/.env`

**Variables**:
- `PUBLIC_API_URL` - Backend API URL
- `PORT` - Frontend port
- `HOST` - Bind address

## Docker Compose Updates

✅ **Updated to use environment variables from .env file**
✅ **Added persistent volume for database** (`tonish-db`)
✅ **CORS configuration updated** to allow 192.168.5.10

## Database Integration

✅ **Database Status**: Working correctly
✅ **Location**: Docker volume `tonish-db:/root/tonish.db`
✅ **Type**: SQLite with GORM ORM
✅ **Auto-migration**: Enabled for all models

### Models Migrated:
- User
- Task
- Notebook
- Page

## Testing Database

### Create a task:
```bash
curl -X POST http://192.168.5.10:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Task",
    "description": "Task description",
    "priority": "high",
    "status": "todo",
    "quadrant": "urgent-important"
  }'
```

### List all tasks:
```bash
curl http://192.168.5.10:8080/api/tasks
```

### Update a task:
```bash
curl -X PUT http://192.168.5.10:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Task",
    "status": "done"
  }'
```

### Delete a task:
```bash
curl -X DELETE http://192.168.5.10:8080/api/tasks/1
```

## CORS Configuration

Updated `/backend/middleware/cors.go` to allow:
- `http://localhost:5173` (dev)
- `http://localhost:4173` (preview)
- `http://192.168.5.10:5173` (production)
- `http://192.168.5.10:4173` (production preview)

## Restarting Services

```bash
# Stop and remove everything
docker-compose down

# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# View specific service
docker-compose logs backend
docker-compose logs frontend
```

## Modifying Environment

1. Edit `.env` file in project root
2. Restart services:
   ```bash
   docker-compose down
   docker-compose up -d
   ```

## Troubleshooting

### Database issues
```bash
# Check database file exists
docker-compose exec backend ls -la /root/*.db

# View database logs
docker-compose logs backend | grep -i database
```

### CORS errors
- Ensure frontend URL is in CORS AllowOrigins
- Check browser console for specific error
- Verify frontend is using correct API URL

### Frontend can't connect
```bash
# Test backend from host
curl http://192.168.5.10:8080/api/tasks

# Check frontend environment
docker-compose exec frontend env | grep API
```

## Security Notes

⚠️ `.env` files are in `.gitignore` - never commit them!
⚠️ Change default values for production
⚠️ Use strong secrets for production deployment
