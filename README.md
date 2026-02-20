# Tonish - Modern Task Management Application

A modern, full-stack task management and note-taking application built with Go (Backend) and SvelteKit (Frontend). Tonish offers a beautiful, space-optimized interface for organizing your tasks, notes, and workflows with real-time synchronization across all devices.

## ✨ Key Features

- 📋 **Dual Task Systems**: Kanban board and Eisenhower Matrix for flexible task management
- 📅 **Calendar View**: Visual task scheduling and tracking
- 📓 **MyFlowBook**: Rich notebook and note-taking system
- 🔄 **Real-time Sync**: WebSocket-based live updates across all connected devices
- 📱 **Responsive Design**: Space-optimized for mobile, tablet, and desktop
- 🎨 **Modern UI**: Clean interface with official Tonish branding
- 🔐 **Secure**: JWT authentication and password hashing
- 🚀 **Fast**: Lightweight Go backend with efficient SvelteKit frontend

## 🚀 Quick Start

Once the application is running:

- **Frontend (Web App)**: http://localhost:50001
- **Backend API**: http://localhost:50002
- **WebSocket Endpoint**: ws://localhost:50002/ws
- **API Health Check**: http://localhost:50002/

## 🔐 Default Login

User self-registration is **disabled by default**. Use the pre-configured account:

- **Email**: `klist@gmail.com`
- **Password**: `Klist123`

To change credentials, update `DEFAULT_USER_EMAIL`, `DEFAULT_USER_PASSWORD`, and `DEFAULT_USER_NAME` in the `.env` file and restart the containers.

## 📋 Prerequisites

Before setting up Tonish, ensure you have:

- **Docker** (version 20.10 or higher)
- **Docker Compose** (version 2.0 or higher)

Verify your installation:
```bash
docker --version
docker-compose --version
```

## 🛠️ Installation & Setup

### Step 1: Clone the Repository
```bash
git clone https://github.com/Sachitasap/Tonish.git
cd Tonish
```

### Step 2: Configure Environment Variables

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:

```env
# Backend Configuration
BACKEND_PORT=50002
JWT_SECRET=your-secure-random-string-here-change-in-production
DATABASE_PATH=/data/tonish.db
DEFAULT_USER_EMAIL=klist@gmail.com
DEFAULT_USER_PASSWORD=Klist123
DEFAULT_USER_NAME=Klist

# Frontend Configuration
FRONTEND_PORT=50001
```

**Important**: Change the `JWT_SECRET` to a secure random string for production use.

### Step 3: Build and Start

Build containers from scratch:
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

Or use the quick start:
```bash
docker-compose up -d --build
```

### Step 4: Access the Application

Open your browser and navigate to:
```
http://localhost:50001
```

## 🎯 Application Features

### 1. Dashboard
- Quick overview of all tasks and notebooks
- Task statistics by status and quadrant
- Recent activity summary

### 2. MyFlow - Task Management
Switch between two powerful task management systems:

**Kanban Board**
- To Do, In Progress, Done columns
- Drag-and-drop task organization
- Quick task creation
- Visual task cards with priority indicators

**Eisenhower Matrix**
- Urgent & Important (Do First)
- Not Urgent & Important (Schedule)
- Urgent & Not Important (Delegate)
- Not Urgent & Not Important (Eliminate)
- Drag-and-drop between quadrants
- Strategic task prioritization

### 3. Calendar
- Monthly calendar view
- Task scheduling with due dates
- Visual task indicators
- Create and edit tasks directly from calendar
- Quick stats for the month

### 4. MyFlowBook - Note Taking
- Create and organize notebooks
- Rich text editing with Markdown support
- Page management within notebooks
- Search functionality
- Pin important notebooks
- Export capabilities

### 5. LookBack - Archive
- Review completed and archived tasks
- Restore or permanently delete tasks
- Filter by status and date
- Task history tracking

### 6. Real-time Synchronization
- Instant updates across all connected devices
- WebSocket-based live sync
- Automatic reconnection
- No manual refresh needed

## 📁 Project Structure

```
Tonish/
├── backend/                 # Go backend application
│   ├── database/           # Database connection and migrations
│   ├── handlers/           # API request handlers
│   │   ├── auth.go        # Authentication endpoints
│   │   ├── task.go        # Task management
│   │   └── notebook.go    # Notebook and page management
│   ├── middleware/         # Authentication and CORS
│   ├── models/             # Data models
│   ├── routes/             # API routes configuration
│   ├── websocket/          # WebSocket hub and handlers
│   │   ├── hub.go         # WebSocket connection manager
│   │   └── handler.go     # WebSocket message handling
│   ├── Dockerfile          # Backend container config
│   ├── go.mod              # Go dependencies
│   └── main.go             # Application entry point
│
├── frontend/               # SvelteKit frontend application
│   ├── src/
│   │   ├── lib/           # Shared utilities and components
│   │   │   ├── api.ts     # API client
│   │   │   ├── websocket.ts # WebSocket client
│   │   │   ├── utils.ts   # Utility functions
│   │   │   ├── assets/    # Images and logos
│   │   │   └── components/
│   │   │       └── TonishLogo.svelte
│   │   ├── routes/        # Application pages
│   │   │   ├── +layout.svelte  # Main layout with navigation
│   │   │   ├── +page.svelte    # Dashboard
│   │   │   ├── login/     # Login page
│   │   │   ├── register/  # Registration (disabled)
│   │   │   ├── myflow/    # Task management (Kanban/Matrix)
│   │   │   ├── calendar/  # Calendar view
│   │   │   ├── myflowbook/ # Notebook management
│   │   │   │   └── [id]/  # Notebook detail pages
│   │   │   └── lookback/  # Archive and history
│   │   ├── app.css        # Global styles and utilities
│   │   └── app.html       # HTML template
│   ├── static/            # Static assets
│   │   ├── tonish-logo.svg # Official logo
│   │   └── manifest.json  # PWA manifest
│   ├── Dockerfile         # Frontend container config
│   └── package.json       # Node dependencies
│
├── docker-compose.yml      # Docker orchestration
├── .env.example           # Environment template
└── README.md              # This file
```

## 🔧 Development

### Running in Development Mode

For local development without Docker:

**Backend:**
```bash
cd backend
go mod download
PORT=50002 JWT_SECRET=dev-secret go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
# Frontend runs on port 5173 by default
# Update API_URL in src/lib/api.ts if needed
```

**Note:** When running locally, ensure WebSocket URL in frontend matches backend port (ws://localhost:50002/ws)

### Stopping the Application

```bash
# Stop containers but keep data
docker-compose down

# Stop containers and remove volumes (⚠️ deletes all data)
docker-compose down -v
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Backend only
docker-compose logs -f backend

# Frontend only
docker-compose logs -f frontend
```

### WebSocket Development

Monitor WebSocket connections in development:

```bash
# Backend logs show WebSocket activity
docker-compose logs -f backend | grep -i websocket

# Expected output:
# WebSocket hub started
# New WebSocket connection: user_id=1
# Broadcasting to user 1: task_create
```

Browser console will show:
- Connection status
- Reconnection attempts
- Incoming real-time messages

## 🗄️ Database

The application uses **SQLite** for data persistence. The database file is stored in a Docker volume named `tonish-db` to ensure data persists across container restarts.

### Models
- **Users**: Authentication and profile information
- **Tasks**: Task items with status, priority, quadrant, and archival support
- **Notebooks**: Note collections with pages
- **Pages**: Individual notebook pages with rich content

### Real-Time Synchronization
- Changes to tasks and notebooks are broadcast via WebSocket
- All connected clients for the same user receive updates instantly
- No polling required - updates are pushed in real-time
- Automatic reconnection ensures consistent sync across devices

## 🛡️ Security

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: Bcrypt for password security
- **CORS Protection**: Configurable CORS middleware
- **Environment Variables**: Sensitive data in environment configuration

## 🚀 Production Deployment

For production deployment, consider:

1. **Environment Variables**:
   - Use strong, unique `JWT_SECRET`
   - Configure appropriate `CORS_ORIGINS`
   - Set production ports if needed
   - Configure WebSocket URL for production domain

2. **HTTPS & WSS**:
   - Use a reverse proxy (nginx, caddy, traefik)
   - Configure SSL certificates (Let's Encrypt recommended)
   - Ensure WebSocket upgrade headers are properly proxied
   - Update frontend to use `wss://` instead of `ws://` for secure WebSocket

3. **Reverse Proxy Configuration** (nginx example):
   ```nginx
   location /ws {
       proxy_pass http://backend:50002;
       proxy_http_version 1.1;
       proxy_set_header Upgrade $http_upgrade;
       proxy_set_header Connection "upgrade";
       proxy_set_header Host $host;
   }
   ```

4. **Database**:
   - Regular backups of the `tonish-db` volume
   - Consider PostgreSQL or MySQL for production scale
   - Set up automated backup schedules

5. **Monitoring & Logging**:
   - Set up centralized logging
   - Monitor WebSocket connection health
   - Configure health check endpoints
   - Track active WebSocket connections

6. **Performance**:
   - Enable HTTP/2 for better WebSocket performance
   - Configure connection pooling
   - Set appropriate timeouts for WebSocket connections
   - Consider load balancing for high traffic

## 🛠️ Technologies Used

### Backend
- **Go 1.23** - Programming language
- **Fiber v2** - High-performance web framework
- **GORM** - ORM for database operations
- **SQLite** - Embedded database
- **JWT** - Token-based authentication
- **Bcrypt** - Password hashing
- **WebSocket** - Real-time bidirectional communication

### Frontend
- **SvelteKit** - Modern frontend framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Fast build tool
- **TailwindCSS** - Utility-first CSS framework
- **Lucide Svelte** - Beautiful icon library
- **WebSocket API** - Real-time client

### DevOps
- **Docker** - Containerization platform
- **Docker Compose** - Multi-container orchestration
- **Alpine Linux** - Minimal container base images

## 📝 API Endpoints

### Authentication
- `POST /api/auth/register` - Create new account (disabled by default)
- `POST /api/auth/login` - Login and receive JWT token
- `GET /api/user/me` - Get current user profile

### Tasks
- `GET /api/tasks` - List all tasks
- `GET /api/tasks/archived` - Get archived tasks
- `GET /api/tasks/status?status=todo` - Filter tasks by status
- `GET /api/tasks/quadrant/:quadrant` - Filter by matrix quadrant
- `GET /api/tasks/:id` - Get specific task
- `POST /api/tasks` - Create new task
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task (soft delete)
- `POST /api/tasks/:id/archive` - Archive task
- `POST /api/tasks/:id/restore` - Restore archived task
- `DELETE /api/tasks/:id/permanent` - Permanently delete task

### Notebooks
- `GET /api/notebooks` - List all notebooks
- `GET /api/notebooks/:id` - Get specific notebook with pages
- `POST /api/notebooks` - Create new notebook
- `PUT /api/notebooks/:id` - Update notebook
- `DELETE /api/notebooks/:id` - Delete notebook and its pages

### Pages
- `GET /api/pages/search?q=keyword` - Search pages
- `GET /api/pages/:id` - Get specific page
- `POST /api/pages` - Create new page
- `PUT /api/pages/:id` - Update page content
- `DELETE /api/pages/:id` - Delete page

### WebSocket
- `WS /ws?user_id=<id>` - WebSocket connection for real-time updates

### WebSocket Message Types
- `task_create` - New task created
- `task_update` - Task updated
- `task_delete` - Task deleted
- `notebook_create` - New notebook created
- `notebook_update` - Notebook/page updated
- `notebook_delete` - Notebook deleted

## 🐛 Troubleshooting

### Services Not Running
1. **Check both services are running:**
   ```bash
   docker-compose ps
   # Both backend and frontend should show "Up"
   ```

2. **Verify backend API:**
   ```bash
   curl http://localhost:50002/
   # Should return: {"message":"Tonish API is running","version":"1.0"}
   ```

3. **Check container logs:**
   ```bash
   docker-compose logs -f backend
   docker-compose logs -f frontend
   ```

### Authentication Issues
1. **Login with default credentials:**
   - Email: `klist@gmail.com`
   - Password: `Klist123`

2. **Test authentication API:**
   ```bash
   curl -X POST http://localhost:50002/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"klist@gmail.com","password":"Klist123"}'
   ```

3. **Registration is disabled by default** - Use the existing user account

### WebSocket Not Connecting
1. **Check WebSocket endpoint:**
   - URL should be: `ws://localhost:50002/ws?user_id=<user_id>`
   - Frontend auto-reconnects with exponential backoff
   - Max 5 retry attempts with increasing delays

2. **Verify in browser console:**
   ```javascript
   // Should see WebSocket connection messages
   // Check for "WebSocket connected" or reconnection attempts
   ```

3. **Backend logs should show:**
   ```
   WebSocket hub started
   New WebSocket connection: user_id=<id>
   ```

### Port Already in Use
If ports 50001 or 50002 are in use:
1. Edit `docker-compose.yml` to change port mappings
2. Update frontend API URL in `src/lib/api.ts`
3. Restart: `docker-compose down && docker-compose up -d`

### Database Issues
```bash
# Reset database (⚠️ deletes all data)
docker-compose down -v
docker-compose up --build
```

### Real-Time Updates Not Working
1. Check browser console for WebSocket errors
2. Verify user_id is correctly passed to WebSocket endpoint
3. Check backend logs for broadcast messages
4. Ensure multiple devices use the same user account for sync

### Container Build Issues
```bash
# Clean rebuild
docker-compose down
docker system prune -a
docker-compose up --build --no-cache
```

### Complete Reset
```bash
# Nuclear option - removes everything
docker-compose down -v
docker system prune -af --volumes
docker-compose up -d --build
```

## 📧 Support

For issues, questions, or contributions, please refer to the project repository.

## 📄 License

[Add your license information here]

---

**Stay Organized, Stay Productive! ✅**
