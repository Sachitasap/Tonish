# Tonish - Quick Start Guide

## 🚀 Quick Start (Development)

### Option 1: Using the start script (Recommended)
```bash
./start-dev.sh
```

### Option 2: Manual start

**Terminal 1 - Backend:**
```bash
cd backend
go run main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```

Then visit:
- Frontend: http://localhost:5173
- Backend API: http://localhost:3000

## 📝 First Steps

1. **Register a new account**
   - Visit http://localhost:5173
   - You'll be redirected to the login page
   - Click "Register" and create your account

2. **Start using MyFlow**
   - Click on "MyFlow" in the navigation
   - Add your first task with the "+ Add Task" button
   - Organize tasks in the Kanban board (To Do, In Progress, Done)
   - Set priorities (Low, Medium, High)

3. **Create your first notebook in MyFlowBook**
   - Click on "MyFlowBook" in the navigation
   - Create a new notebook with "+ New Notebook"
   - Add pages to your notebook
   - Use tags to organize your notes

## 🐳 Docker Deployment

### Build and run with Docker Compose:
```bash
docker-compose up -d
```

Access the application at:
- Frontend: http://localhost:4173
- Backend API: http://localhost:3000

### Stop the application:
```bash
docker-compose down
```

### View logs:
```bash
docker-compose logs -f
```

## 🔧 Configuration

### Backend Configuration

The backend runs on port 3000 by default. Key files to configure:

1. **JWT Secret** (IMPORTANT for production)
   - File: `backend/handlers/auth.go`
   - File: `backend/middleware/auth.go`
   - Change `jwtSecret` to a secure random string

2. **CORS Settings**
   - File: `backend/middleware/cors.go`
   - Update `AllowOrigins` for your production domain

### Frontend Configuration

The frontend connects to the backend at `http://localhost:3000` by default.

- File: `frontend/src/lib/api.ts`
- Update `API_BASE_URL` for production deployment

## 📊 Database

- Database file: `backend/tonish.db` (SQLite)
- Automatic migrations run on startup
- Backup the `.db` file regularly for data safety

## 🛡️ Security Checklist for Production

- [ ] Change JWT secret to a strong random value
- [ ] Use HTTPS/TLS encryption
- [ ] Update CORS to allow only your domain
- [ ] Set up database backups
- [ ] Use environment variables for secrets
- [ ] Enable firewall rules
- [ ] Regular security updates

## 📱 Testing the API

### Register a user
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```

### Login
```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Create a task (use token from login)
```bash
curl -X POST http://localhost:3000/api/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"title":"My First Task","priority":"high","status":"todo"}'
```

## 🐛 Troubleshooting

### Backend won't start
- Check if Go is installed: `go version`
- Check if port 3000 is available
- Review logs for database errors

### Frontend won't start
- Check if Node.js is installed: `node --version`
- Run `npm install` in the frontend directory
- Check if port 5173 is available

### Can't login
- Make sure backend is running
- Check browser console for errors
- Verify CORS settings in backend

### Database errors
- Delete `tonish.db` file to reset (loses all data)
- Check file permissions
- Ensure SQLite is available

## 💡 Tips

- Use the Kanban board in MyFlow for visual task management
- Tag your notebooks and pages for easy searching
- Use Quick Tasks for rapid capture of ideas
- Regularly backup your `tonish.db` file
- Check the API documentation in README.md

## 🎯 Next Steps

1. Explore the Eisenhower Matrix for task prioritization
2. Create multiple notebooks for different projects
3. Set up Docker deployment for production
4. Customize the UI with Tailwind classes
5. Add calendar integration (future feature)

## 📞 Support

If you encounter issues:
1. Check the logs (backend terminal or `docker-compose logs`)
2. Review the README.md for detailed documentation
3. Check the browser console for frontend errors

---

Happy productivity! 🚀
