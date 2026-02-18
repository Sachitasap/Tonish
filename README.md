# Tonish - Modern Note-Taking Application

A modern, full-stack note-taking application built with Go (Backend) and SvelteKit (Frontend). Tonish offers a beautiful, intuitive interface for organizing your notes, notebooks, and tasks with real-time search and seamless user experience.

## 🚀 Quick Access

Once the application is running:

- **Frontend (Web App)**: http://192.168.4.213:50001
- **Backend API**: http://192.168.4.213:50002
- **API Health Check**: http://192.168.4.213:50002/

**Note**: The application is configured to run on IP address 192.168.4.213 with frontend on port 50001 and backend on port 50002. Access from any device on the same network using these URLs.

## 🔐 Login Information

User self-registration is **temporarily disabled**. Use the pre-provisioned default account to sign in:

- **Email**: `klist@gmail.com`
- **Password**: `Klist123`

To change these credentials, update `DEFAULT_USER_EMAIL`, `DEFAULT_USER_PASSWORD`, and `DEFAULT_USER_NAME` in the `.env` file (or supply them through your infrastructure secrets) and restart the containers. The backend seeds the specified account automatically if it does not already exist.

**Note**: All user data, including the default account, is stored locally in the SQLite database within the Docker volume.

## 📋 Prerequisites

Before setting up Tonish, ensure you have the following installed:

- **Docker** (version 20.10 or higher)
- **Docker Compose** (version 2.0 or higher)

To verify your installation:
```bash
docker --version
docker-compose --version
```

## 🛠️ Initial Setup

### Step 1: Clone the Repository
```bash
git clone <repository-url>
cd Tonish
```

### Step 2: Configure Environment Variables

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit the `.env` file with your preferred configuration:

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
BACKEND_URL=http://192.168.4.213:50002

# Production Settings (optional)
# CORS_ORIGINS=https://yourdomain.com
```

**Important**: Change the `JWT_SECRET` to a secure random string for production use. Update the default user credentials (or remove them) before deploying to production.

### Step 3: Build and Start the Application

```bash
docker-compose up --build
```

This command will:
- Build the backend Go application
- Build the frontend SvelteKit application
- Start both services with networking configured
- Create persistent database volume

### Step 4: Access the Application

Open your browser and navigate to:
```
http://192.168.4.213:50001
```

You can also access it from other devices on the same network (192.168.4.0/22) using this URL.

## 🎯 Usage Guide

### Features

1. **Notebook Management**
   - Create, edit, and delete notebooks
   - Organize notes in collections
   - Visual cards with gradient backgrounds

2. **Note Taking**
   - Rich text editing
   - Markdown support
   - Quick search functionality

3. **Task Management**
   - Create and track tasks
   - Task organization within notebooks

4. **MyFlow & Lookback**
   - Daily workflow tracking
   - Reflection and review features

### Navigation

- **Home**: Dashboard overview
- **MyFlowBook**: Main notebook interface
- **MyFlow**: Daily workflow tracker
- **Lookback**: Review and reflection area

## 📁 Project Structure

```
Tonish/
├── backend/                 # Go backend application
│   ├── database/           # Database connection and migrations
│   ├── handlers/           # API request handlers
│   ├── middleware/         # Authentication and CORS
│   ├── models/             # Data models
│   ├── routes/             # API routes
│   ├── Dockerfile          # Backend container config
│   ├── go.mod              # Go dependencies
│   └── main.go             # Application entry point
│
├── frontend/               # SvelteKit frontend application
│   ├── src/
│   │   ├── lib/           # Shared utilities and components
│   │   │   ├── api.ts     # API client
│   │   │   └── components/ # Reusable UI components
│   │   └── routes/        # Application pages
│   │       ├── login/     # Login page
│   │       ├── register/  # Registration page
│   │       ├── myflowbook/ # Notebook pages
│   │       ├── myflow/    # Workflow tracker
│   │       └── lookback/  # Review page
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
PORT=8080 JWT_SECRET=dev-secret go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

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

## 🗄️ Database

The application uses **SQLite** for data persistence. The database file is stored in a Docker volume named `tonish-db` to ensure data persists across container restarts.

### Models
- **Users**: Authentication and profile information
- **Notebooks**: Note collections
- **Tasks**: Task items with status tracking

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

2. **HTTPS**:
   - Use a reverse proxy (nginx, traefik)
   - Configure SSL certificates

3. **Database Backup**:
   - Regular backups of the `tonish-db` volume
   - Consider using a more robust database for production

4. **Monitoring**:
   - Set up logging and monitoring
   - Configure health checks

## 🛠️ Technologies Used

### Backend
- **Go** (Golang) - Programming language
- **Fiber** - Web framework
- **GORM** - ORM for database operations
- **SQLite** - Database
- **JWT** - Authentication
- **Bcrypt** - Password hashing

### Frontend
- **SvelteKit** - Frontend framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **TailwindCSS** - Styling
- **Lucide** - Icons

### DevOps
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration

## 📝 API Endpoints

### Authentication
- `POST /api/auth/register` - Create new account
- `POST /api/auth/login` - Login and get JWT token
- `GET /api/auth/me` - Get current user (requires auth)

### Notebooks
- `GET /api/notebooks` - List all notebooks
- `POST /api/notebooks` - Create notebook
- `GET /api/notebooks/:id` - Get specific notebook
- `PUT /api/notebooks/:id` - Update notebook
- `DELETE /api/notebooks/:id` - Delete notebook

### Tasks
- `GET /api/tasks` - List all tasks
- `POST /api/tasks` - Create task
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task

## 🐛 Troubleshooting

### Testing Registration
If you're having issues with registration, use the included test page:
```bash
# Open test-registration.html in your browser
open test-registration.html  # macOS
xdg-open test-registration.html  # Linux
# Or navigate to file:///path/to/Tonish/test-registration.html
```

This standalone test page will help you:
- Verify the backend API is accessible
- Test registration without the frontend
- See detailed error messages
- Check for CORS issues

### Registration Not Working
If registration fails:

1. **Check both services are running:**
   ```bash
   docker-compose ps
   # Both backend and frontend should show "Up"
   ```

2. **Verify backend API:**
   ```bash
   curl http://localhost:8080/
   # Should return: {"message":"Tonish API is running","version":"1.0"}
   ```

3. **Test registration directly:**
   ```bash
   curl -X POST http://localhost:8080/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@test.com","password":"test123456","name":"Test User"}'
   ```

4. **Rebuild containers if needed:**
   ```bash
   docker-compose down
   docker-compose up -d --build
   ```

### Port Already in Use
If ports 5173 or 8080 are already in use:
1. Edit `.env` file and change `FRONTEND_PORT` or `BACKEND_PORT`
2. Restart with `docker-compose up`

### Database Issues
```bash
# Reset database (⚠️ deletes all data)
docker-compose down -v
docker-compose up --build
```

### Container Build Issues
```bash
# Clean rebuild
docker-compose down
docker system prune -a
docker-compose up --build
```

### "ContainerConfig" Error
If you see a "ContainerConfig" KeyError:
```bash
# Complete reset
docker-compose down -v
docker-compose up -d --force-recreate
```

## 📧 Support

For issues, questions, or contributions, please refer to the project repository.

## 📄 License

[Add your license information here]

---

**Happy Note-Taking! 📝**
