# 🚀 Tonish - Personal Productivity System

Tonish is a **full-stack web application** combining task management (MyFlow) and digital note-taking (MyFlowBook) in one unified platform. Built with **SvelteKit** (frontend) and **Go with Fiber** (backend), it provides a responsive, touch-friendly experience across all devices.

## ✨ Features

### 🗂️ MyFlow (Task Management) - **COMPLETE**
- ✅ Full CRUD operations (create, edit, delete tasks)
- 📋 **Dual task management**:
  - **Kanban board** for personal projects (To Do, In Progress, Done columns, horizontally scrollable on mobile)
  - **Eisenhower Matrix** for daily tasks & billing (Do First, Schedule, Delegate, Eliminate quadrants)
- 🎯 Priority levels (Low, Medium, High) with visual indicators
- 🏷️ Tags and categorization for easy organization
- 📅 **Always-visible due dates** with overdue task highlighting and calendar icons
- ✏️ **Inline editing** – Edit any task's title, description, and due date directly in Kanban or Matrix view
- 🎲 **Eisenhower Matrix enhancements**:
  - Drag-and-drop between quadrants for flexibility
  - **+ quick-add buttons** on each quadrant for rapid task entry
  - Edit buttons on every matrix task
  - **Auto-quadrant assignment** based on priority and due date
- ⚡ Quick tasks for rapid entry without full form
- 🔄 Smart status transitions between columns
- 🎨 **Professional icon system** using lucide-svelte (100% emoji-free UI)
- 📱 **Fully responsive** – touch-optimized controls on mobile, desktop-optimized on larger screens

### 📓 MyFlowBook (Digital Notebook) - **COMPLETE**
- 📓 **Multiple notebooks** with unlimited pages per notebook
- ✍️ **Rich-text editing** with markdown support (**bold**, *italic*, `code`, # headings)
- 🏷️ **Tagging system** for both notebooks and individual pages
- 🔍 **Full-text search** across all pages (searches titles and content)
- 📌 **Pin/unpin** important notebooks and pages for quick access
- 📤 **Export functionality** in 3 formats:
  - 📄 **Text files** (formatted plaintext)
  - 📊 **CSV files** (structured data)
  - 📝 **Markdown files** (portable format)
- 🎨 **Glass-morphism sidebar** with 50% opacity backdrop blur and professional indigo/blue palette
- 📋 **Copy & duplicate** functionality for quick content reuse
- 📊 **Organized views** by pins, recent items, and tags (with counts)
- 🎨 **Professional icon system** using lucide-svelte (100% emoji-free UI)
- 🔐 **Content security** – password hashing, no sensitive data exposure

## 💻 Tech Stack

### Frontend
- **Framework**: [SvelteKit](https://kit.svelte.dev/) v2.47.1 with TypeScript
- **Language**: TypeScript 5.9.3 with strict type checking
- **Styling**: [TailwindCSS](https://tailwindcss.com/) 3.4.9 with custom animations and glass-morphism effects
- **Icons**: [lucide-svelte](https://lucide.dev/) for professional, consistent icon system
- **State Management**: **Svelte 5 runes** ($state, $derived for reactive data)
- **Build Tool**: [Vite](https://vitejs.dev/) 7.1.10 (fast hot module replacement)
- **Brand Logo**: Custom TonishLogo component with handwriting font and animated three-dot progression

### Backend
- **Language**: Go 1.21+ (compiled, high-performance)
- **Framework**: [Fiber](https://gofiber.io/) v2 (lightweight, fast HTTP web framework)
- **Database**: SQLite with [GORM](https://gorm.io/) ORM (zero-config, file-based)
- **Authentication**: JWT tokens (7-day expiration)
- **Password Security**: bcrypt hashing
- **API Design**: RESTful JSON API with proper HTTP status codes
- **Middleware**: CORS enabled, request logging

### Deployment
- **Containerization**: Multi-stage Docker builds (optimized image sizes)
- **Orchestration**: Docker Compose for local & development environments
- **Resource Usage**: 
  - Backend: ~50MB RAM (Go binary)
  - Frontend: ~200MB RAM (Node dev server)
  - Database: SQLite (single file, <5MB typical)
  - **Total**: ~300-400MB for full stack (single-user, development)
- **Scalability**: Ready for reverse proxy (nginx), load balancing, or cloud deployment

## Project Structure

```
Tonish/
├── backend/
│   ├── config/          # Configuration files
│   ├── database/        # Database connection and migrations
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Authentication and CORS middleware
│   ├── models/          # Data models (User, Task, Notebook, Page)
│   ├── routes/          # API route definitions
│   ├── main.go          # Application entry point
│   ├── Dockerfile       # Backend Docker configuration
│   └── go.mod           # Go dependencies
├── frontend/
│   ├── src/
│   │   ├── lib/         # Shared utilities and API client
│   │   └── routes/      # SvelteKit pages and layouts
│   ├── Dockerfile       # Frontend Docker configuration
│   └── package.json     # Node dependencies
└── docker-compose.yml   # Multi-container orchestration
```

## 🚀 Getting Started

### Prerequisites
- **Go** 1.21 or later ([download](https://go.dev/dl))
- **Node.js** 20 or later ([download](https://nodejs.org/))
- **Docker** & **Docker Compose** (optional, for containerized deployment)
- **Git** (for cloning the repository)

### Development Setup

Quick start (both services):
```bash
# Clone and navigate to project
git clone <repo-url>
cd Tonish

# Start both backend and frontend
./start-dev.sh
```

Or start individually:

#### Backend Setup
```bash
cd backend
go mod download      # Download dependencies
go run main.go       # Start server
```
✅ Backend: http://localhost:3000  
🗄️ Database: Auto-created `tonish.db` on first run

#### Frontend Setup
```bash
cd frontend
npm install          # Install dependencies
npm run dev          # Start dev server with hot reload
```
✅ Frontend: http://localhost:5173  
🔄 Hot Module Replacement: Changes reload instantly

### Docker Deployment
```bash
# Build and start all services in the background
docker-compose up -d

# View real-time logs
docker-compose logs -f

# View logs for specific service
docker-compose logs -f backend   # Backend logs
docker-compose logs -f frontend  # Frontend logs

# Stop all services (keep data)
docker-compose stop

# Stop and remove containers (data persists in volumes)
docker-compose down

# Full cleanup (remove volumes too)
docker-compose down -v
```
✅ Frontend: http://localhost:5001 (mapped from internal 4173)  
✅ Backend: http://localhost:3000

## 🔌 API Endpoints

All endpoints use **JSON** request/response format. Base URL: `http://localhost:3000/api`

### 🔐 Authentication
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/auth/register` | Create new user account | No |
| POST | `/auth/login` | Authenticate and get JWT token | No |
| GET | `/user/me` | Get current authenticated user | Yes |

**Token Format**: `Authorization: Bearer <JWT_TOKEN>`

### ✅ Tasks (MyFlow)
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/tasks` | Get all tasks | Yes |
| GET | `/tasks/:id` | Get single task | Yes |
| POST | `/tasks` | Create new task | Yes |
| PUT | `/tasks/:id` | Update task | Yes |
| DELETE | `/tasks/:id` | Delete task | Yes |
| GET | `/tasks/status?status=todo` | Filter tasks by status (todo, in-progress, done) | Yes |
| GET | `/tasks/quadrant/:quadrant` | Filter by Eisenhower quadrant | Yes |
| GET | `/tasks/type?type=kanban` | Filter by task type (kanban for projects, matrix for daily tasks) | Yes |

### 📓 Notebooks (MyFlowBook)
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/notebooks` | Get all notebooks (with pages preloaded) | Yes |
| GET | `/notebooks/:id` | Get single notebook with all pages | Yes |
| POST | `/notebooks` | Create new notebook | Yes |
| PUT | `/notebooks/:id` | Update notebook (name, tags, pin status) | Yes |
| DELETE | `/notebooks/:id` | Delete notebook and all pages | Yes |

### 📄 Pages
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/pages/:id` | Get single page | Yes |
| POST | `/pages` | Create new page in notebook | Yes |
| PUT | `/pages/:id` | Update page (title, content, tags) | Yes |
| DELETE | `/pages/:id` | Delete page | Yes |
| GET | `/pages/search?q=query` | Full-text search across all pages | Yes |

## ⚙️ Configuration

### Backend Environment Variables

Create `.env` in the `backend` directory (optional, defaults provided):

```dotenv
PORT=3000                                    # HTTP server port
JWT_SECRET=your-secret-key-here             # CHANGE THIS IN PRODUCTION!
DATABASE_PATH=tonish.db                     # SQLite database file
```

**Note**: A default `.env.example` is provided in the repository.

### Frontend Environment Variables

Create `.env` in the `frontend` directory (optional):

```dotenv
VITE_API_URL=http://localhost:3000/api      # Backend API URL (dev)
```

## 🔒 Security & Production Checklist

⚠️ **Before deploying to production:**

- [ ] **JWT Secret**: Change the default JWT secret (currently: `your-secret-key-change-in-production`)
  - Generate: `openssl rand -base64 32`
  - Set in environment: `JWT_SECRET=<generated-secret>`

- [ ] **HTTPS/TLS**: Enable HTTPS on your domain
  - Use Let's Encrypt or similar for free certificates
  - Configure reverse proxy (nginx/Apache) to handle TLS

- [ ] **CORS Origins**: Update allowed CORS origins in `backend/middleware/cors.go`
  - Default allows all origins (fine for local dev only)
  - Set to your actual domain: `https://yourdomain.com`

- [ ] **Database**: Set up automated backups
  - Back up `tonish.db` regularly (daily recommended)
  - Consider migrating to PostgreSQL for multi-user scenarios

- [ ] **Rate Limiting**: Implement rate limiting for auth endpoints
  - Prevent brute-force attacks on login/register

- [ ] **Password Policy**: Consider enforcing minimum password strength
  - Current: No minimum (bcrypt only)

- [ ] **Environment Variables**: Store sensitive values in `.env` (never commit)
  - Add `.env` to `.gitignore` ✅ (already done)

- [ ] **Logging**: Monitor and log authentication events
  - Set up log aggregation (ELK, Datadog, etc.)

- [ ] **API Rate Limiting**: Protect against abuse
  - Rate limit per IP or user

⚠️ **Current State**: No multi-user isolation per user_id (comment in code indicates future work)

## 🗄️ Database Schema

SQLite with GORM ORM. Auto-migrated on startup.

### Users Table
```sql
CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,        -- bcrypt hash
  name TEXT,
  created_at DATETIME,
  updated_at DATETIME
);
```

### Tasks Table
```sql
CREATE TABLE tasks (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  priority TEXT DEFAULT 'medium',  -- low, medium, high
  status TEXT DEFAULT 'todo',       -- todo, in-progress, done (for kanban tasks)
  tags TEXT,                        -- JSON array as string
  due_date DATETIME,
  is_quick_task BOOLEAN DEFAULT 0,
  quadrant TEXT,                    -- Eisenhower: urgent-important, not-urgent-important, urgent-not-important, not-urgent-not-important
  task_type TEXT DEFAULT 'kanban',  -- kanban (personal projects) or matrix (daily tasks/billing)
  user_id INTEGER,                  -- Foreign key to users
  created_at DATETIME,
  updated_at DATETIME
);
```

### Notebooks Table
```sql
CREATE TABLE notebooks (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  tags TEXT,                        -- JSON array as string
  is_pinned BOOLEAN DEFAULT 0,
  user_id INTEGER,                  -- Foreign key to users
  created_at DATETIME,
  updated_at DATETIME
);
```

### Pages Table
```sql
CREATE TABLE pages (
  id INTEGER PRIMARY KEY,
  notebook_id INTEGER NOT NULL,     -- Foreign key to notebooks
  title TEXT NOT NULL,
  content TEXT,                     -- Rich-text JSON content
  tags TEXT,                        -- JSON array as string
  is_pinned BOOLEAN DEFAULT 0,
  created_at DATETIME,
  updated_at DATETIME
);
```

**Total Tables**: 4  
**Primary Key**: UUID (uint in Go)  
**Relationships**: One-to-many (User → Tasks/Notebooks, Notebook → Pages)

## 🗓️ Roadmap & Future Enhancements

### Phase 2 (Next Priority)
- 🔐 **Per-user data isolation**: Enforce user_id filtering in all queries
- 📱 **Mobile app**: SvelteKit Capacitor wrapper (iOS/Android)
- 📊 **Calendar view**: Visualize tasks and due dates
- 📈 **Analytics dashboard**: Productivity insights (completed tasks, trends)
- 🎨 **Themes**: Light/dark mode, custom color schemes

### Phase 3 (Future)
- ☁️ **Cloud sync**: Optional backup to S3 or Google Drive
- 🤖 **AI features**: Smart task suggestions, content summarization
- 🔔 **Notifications**: Email digests, task reminders
- 📝 **Collaborative notebooks**: Real-time sharing & editing
- 🔄 **Offline mode**: Service worker support, sync on reconnect
- 📄 **PDF export**: Professional report generation
- 🌐 **Multi-user deployment**: PostgreSQL migration, user management UI
- 🔍 **Advanced search**: Filters, regex support, saved searches
- ⏱️ **Time tracking**: Pomodoro timer, task duration tracking
- 🏷️ **Custom tags & colors**: User-defined tag system with icons

## 📋 Implementation Status

### ✅ Completed Features
- [x] Task management (MyFlow) – full CRUD with Kanban & Eisenhower Matrix
- [x] **Task type separation** – Kanban for personal projects, Eisenhower Matrix for daily tasks & billing
- [x] **Inline task editing** – Edit title, description, and due dates directly in list/matrix views
- [x] **Always-visible due dates** – Calendar icons and formatted dates on all tasks
- [x] **Eisenhower quick-add** – + buttons on each matrix quadrant for rapid task creation
- [x] **Drag-and-drop** – Move tasks between Eisenhower Matrix quadrants
- [x] Note-taking (MyFlowBook) – notebooks, pages, rich-text editing
- [x] Search functionality – full-text across pages
- [x] Export (Text, CSV, Markdown)
- [x] Responsive UI – desktop & mobile optimized
- [x] Authentication (JWT)
- [x] Type safety (TypeScript frontend, Go backend)
- [x] Markdown support
- [x] Pin/tag functionality
- [x] **Professional icon system** – 100% emoji-free UI with lucide-svelte
- [x] **Glass-morphism sidebar** – Modern design with backdrop blur effects
- [x] **Custom Tonish logo** – Handwriting-style text with animated three-dot animation
- [x] **Code cleanup & refactoring** – Removed duplication, consolidated helper functions for priority badges

### 🚧 In Progress / Known Limitations
- [ ] Per-user data isolation (all users see all data currently)
- [ ] Real-time collaboration
- [ ] Advanced search filters
- [ ] Mobile app

### ❌ Not Yet Implemented
- [ ] Offline mode / sync
- [ ] Calendar views
- [ ] Analytics
- [ ] Email notifications
- [ ] Multi-user editing

## 📄 License

MIT License – Free to use, modify, and distribute. See [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 💬 Support & Issues

- 📝 **Bug reports**: [Create an issue](../../issues)
- 💡 **Feature requests**: [Start a discussion](../../discussions)
- 📧 **Email**: [contact@example.com]

## 🙏 Acknowledgments

- [SvelteKit](https://kit.svelte.dev/) – Amazing framework
- [Fiber](https://gofiber.io/) – Fast Go web framework
- [TailwindCSS](https://tailwindcss.com/) – Utility-first CSS
- [GORM](https://gorm.io/) – Go ORM

---

**Built with ❤️ using SvelteKit, TypeScript, Go, and Fiber**

⭐ If you find this useful, please star the repository!
