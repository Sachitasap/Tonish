# Tonish — Personal Productivity Suite

A full-stack productivity application built with **Go** (backend) and **SvelteKit** (frontend). Tonish provides a unified workspace for task management, calendar scheduling, and rich note-taking — all synchronized in real-time across devices.

---

## ✨ Features

### Task Management — MyFlow
- **Kanban Board** — To Do / In Progress / Done columns
- **Eisenhower Matrix** — Prioritize tasks across four urgency/importance quadrants
- Drag-and-drop task cards
- Due dates, priority levels, and status tracking

### Calendar
- Monthly calendar view with per-day task indicators
- Create tasks directly from any date cell
- **Recurring tasks** — daily, weekly, monthly, or yearly with configurable occurrence count
- Month-at-a-glance stats

### MyFlowBook — Notebooks
- **Color-coded notebooks** — choose from 8 accent colors; notebooks are grouped by color on the list page
- Pin important notebooks to the top
- **Page creation modal** with 4 quick-start templates (Blank, Meeting Notes, Daily Log, Code Snippet) and a mini markdown toolbar
- Markdown-aware editor with live preview (bold, italic, code, headings)
- Per-page copy, duplicate, and inline export (`.md` / `.txt`)
- **Export entire notebook** as PDF, Markdown, Plain Text, CSV, or JSON
- Full-text search within pages

### LookBack — Archive
- Review completed and archived tasks
- Restore or permanently delete
- Filter by status and date

### General
- **Real-time sync** via WebSocket — changes propagate instantly across all browser tabs/devices
- **Responsive layout** — optimized for mobile, tablet, and desktop
- Sticky nav bar with responsive date pill (adapts label length to screen width)
- JWT authentication with bcrypt password hashing

---

## 🚀 Quick Start

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/) 20.10+
- [Docker Compose](https://docs.docker.com/compose/) v2+

```bash
docker --version        # Docker version 20.x+
docker compose version  # Docker Compose version v2.x+
```

### 1 — Clone
```bash
git clone https://github.com/Sachitasap/Tonish.git
cd Tonish
```

### 2 — Configure environment
```bash
cp .env.example .env
```

Open `.env` and update at minimum the `JWT_SECRET`:

```env
# Backend
BACKEND_PORT=50002
JWT_SECRET=change-this-to-a-secure-random-string-in-production
DATABASE_PATH=/data/tonish.db
DEFAULT_USER_EMAIL=klist@gmail.com
DEFAULT_USER_PASSWORD=Klist123
DEFAULT_USER_NAME=Klist

# Frontend
FRONTEND_PORT=50001
BACKEND_URL=http://<your-server-ip>:50002
```

> Replace `<your-server-ip>` with `localhost` for local use, or your machine's LAN/public IP for network access.

### 3 — Build & run
```bash
docker compose up -d --build
```

### 4 — Open the app

| Service | URL |
|---|---|
| Web App | http://localhost:50001 |
| Backend API | http://localhost:50002 |
| API Health | http://localhost:50002/ |
| WebSocket | ws://localhost:50002/ws |

---

## 🔐 Login

User self-registration is **disabled by default**. Use the pre-configured account (set in `.env`):

| Field | Default value |
|---|---|
| Email | `klist@gmail.com` |
| Password | `Klist123` |

To change credentials, update `DEFAULT_USER_EMAIL`, `DEFAULT_USER_PASSWORD`, and `DEFAULT_USER_NAME` in `.env`, then rebuild:

```bash
docker compose down -v   # ⚠️ wipes the existing database
docker compose up -d --build
```

---

## 📁 Project Structure

```
Tonish/
├── backend/
│   ├── database/        # SQLite connection & auto-migration
│   ├── handlers/        # HTTP route handlers (auth, tasks, notebooks, pages)
│   ├── middleware/      # JWT auth & CORS middleware
│   ├── models/          # GORM data models (User, Task, Notebook, Page)
│   ├── routes/          # Route registration
│   ├── websocket/       # WebSocket hub & broadcast
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
│
├── frontend/
│   └── src/
│       ├── lib/
│       │   ├── api.ts          # Typed REST API client
│       │   ├── websocket.ts    # WS client with auto-reconnect
│       │   ├── utils.ts
│       │   └── components/
│       └── routes/
│           ├── +layout.svelte  # App shell, nav bar, bottom mobile nav
│           ├── +page.svelte    # Dashboard
│           ├── login/
│           ├── myflow/         # Kanban + Eisenhower matrix
│           ├── calendar/       # Calendar + recurring task creation
│           ├── myflowbook/     # Notebook list (color-coded groups)
│           │   └── [id]/       # Notebook detail — pages + export
│           └── lookback/       # Archived tasks
│
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## 🗄️ Database

SQLite is stored in Docker volume `tonish_tonish-db` at `/data/tonish.db` inside the backend container. The volume persists across `docker compose` restarts as long as you **do not** pass the `-v` flag.

### Back up the database
```bash
docker compose cp backend:/data/tonish.db ./tonish-backup.db
```

### Reset the database
```bash
docker compose down -v   # removes the volume → fresh DB on next start
docker compose up -d
```

---

## 🛠️ Development (without Docker)

**Backend**
```bash
cd backend
go mod download
PORT=50002 JWT_SECRET=dev-secret DATABASE_PATH=./tonish.db go run main.go
```

**Frontend**
```bash
cd frontend
npm install
npm run dev     # serves on :5173
```

Set `BACKEND_URL_INTERNAL` or adjust the proxy target in `src/hooks.server.ts` when running the backend on a non-default port.

---

## 📝 API Reference

### Auth
| Method | Path | Description |
|---|---|---|
| POST | `/api/auth/login` | Login → returns JWT |
| POST | `/api/auth/register` | Register (disabled by default) |
| GET | `/api/user/me` | Current user profile |

### Tasks
| Method | Path | Description |
|---|---|---|
| GET | `/api/tasks` | All tasks |
| GET | `/api/tasks/archived` | Archived tasks |
| GET | `/api/tasks/status?status=todo` | Filter by status |
| GET | `/api/tasks/quadrant/:q` | Filter by Eisenhower quadrant |
| POST | `/api/tasks` | Create task |
| PUT | `/api/tasks/:id` | Update task |
| DELETE | `/api/tasks/:id` | Soft delete |
| POST | `/api/tasks/:id/archive` | Archive |
| POST | `/api/tasks/:id/restore` | Restore from archive |
| DELETE | `/api/tasks/:id/permanent` | Permanent delete |

### Notebooks & Pages
| Method | Path | Description |
|---|---|---|
| GET | `/api/notebooks` | All notebooks (includes pages) |
| GET | `/api/notebooks/:id` | Single notebook |
| POST | `/api/notebooks` | Create notebook |
| PUT | `/api/notebooks/:id` | Update notebook |
| DELETE | `/api/notebooks/:id` | Delete notebook + all pages |
| GET | `/api/pages/search?q=` | Full-text page search |
| POST | `/api/pages` | Create page |
| PUT | `/api/pages/:id` | Update page |
| DELETE | `/api/pages/:id` | Delete page |

### WebSocket
| | |
|---|---|
| Endpoint | `WS /ws?user_id=<id>` |
| Events | `task_create` · `task_update` · `task_delete` · `notebook_create` · `notebook_update` · `notebook_delete` |

---

## 🐛 Troubleshooting

### App not loading
```bash
docker compose ps                       # both containers should show "Up"
docker compose logs backend --tail=40   # check for startup errors
docker compose logs frontend --tail=20
```

### Verify backend is alive
```bash
curl http://localhost:50002/
# {"message":"Tonish API is running","version":"1.0"}
```

### Test login
```bash
curl -X POST http://localhost:50002/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"klist@gmail.com","password":"Klist123"}'
# Should return {"token":"..."}
```

### Login fails despite correct password
Clear browser localStorage (remove the `authToken` key) — you may have a stale token from a previous session. Open DevTools → Application → Local Storage → delete `authToken`, then refresh.

### Data disappeared after restart
You likely ran `docker compose down -v`. Always stop without `-v` to keep the database volume:
```bash
docker compose down     # stops containers, keeps data
docker compose up -d    # restarts with existing DB
```

### Port conflicts
Edit `FRONTEND_PORT` / `BACKEND_PORT` in `.env` and restart.

### Clean rebuild
```bash
docker compose down
docker compose build --no-cache
docker compose up -d
```

### Full reset ⚠️ (deletes all data)
```bash
docker compose down -v
docker system prune -af --volumes
docker compose up -d --build
```

---

## 🛡️ Security Notes

- Set `JWT_SECRET` to a long random string before exposing the app on a network
- For HTTPS, place a reverse proxy (nginx / Caddy) in front and forward `Upgrade` headers for WebSocket (`wss://`)
- Restrict `CORS_ORIGINS` in `.env` to specific domains in production

---

## 🧰 Tech Stack

| Layer | Technology |
|---|---|
| Backend language | Go 1.23 |
| Web framework | Fiber v2 |
| ORM | GORM |
| Database | SQLite |
| Auth | JWT + Bcrypt |
| Real-time | WebSocket (gorilla) |
| Frontend | SvelteKit + Svelte 5 (runes) |
| Language | TypeScript |
| Styling | Tailwind CSS |
| Icons | Lucide Svelte |
| Build tool | Vite |
| Containers | Docker + Docker Compose |
| Base images | Alpine Linux |

---

> **Stay organized, stay productive. ✅**
