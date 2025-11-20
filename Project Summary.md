📝 Project Summary: Tonish
Tonish is a personal productivity system designed for web and mobile platforms(In future), combining task management, scheduling, prioritization,Kanabn board,Eisenhower Matrix and note-taking to store for future in one application. The application is structured into two main sections:
Section 1: MyFlow
Purpose: Personal productivity dashboard and task management.
Features:
Daily Tasks: Create, edit, delete, mark complete, set priority, due dates, and tags.
Quick Tasks: Rapid task addition with the option to convert into full tasks later.
Calendar View: Month/Week/Day view with drag-and-drop task scheduling.
Kanban Board: Columns (To Do, In Progress, Done) with draggable task cards.
Eisenhower Matrix: Four quadrants for prioritization (Urgent/Important, etc.) with automatic sorting.
Task Export: Export tasks to CSV or PDF.
Section 2: MyFlowBook
Purpose: Digital notebook for work logs, notes, and documentation.
Features:
Multiple notebooks, each containing multiple pages.
Rich-text editor (TipTap) supporting headings, lists, checklists, and formatting.
Tagging system for easy organization.
Search notes by title or content.
Export single notes or entire notebooks to PDF or CSV.
Pin important notes and maintain workspace clarity.
🏗 Tech Stack
Frontend
Web: SvelteKit + TailwindCSS + shadcn-svelte components
Future Mobile: Kotlin + Jetpack Compose (Android)
FullCalendar for calendar views
dnd-kit for Kanban drag/drop
TipTap editor for rich-text notes
Backend
Golang (Fiber framework)
REST API (JSON) for tasks, notes, calendar, kanban, CSV/Google Cloud sync
JWT or OAuth2 authentication
Database
Primary: SQLite (single-user)
JSONB fields for rich-text notes
Arrays for tags
Full-text search support
Cloud & Sync
Optional Google Cloud integration for CSV export and storage
Service account or OAuth2 authentication
Go backend handles parsing, uploading, and syncing
Deployment
Docker Compose for web app and backend
Optional reverse proxy with Caddy or Nginx
Lightweight resource footprint: ~300–400MB RAM for single-user deployment
📱 Future Mobile App Support
Android app connects to the same Golang backend API
Features mirror MyFlow + MyFlowBook functionality
Push notifications, offline cache, and incremental sync supported
Multiple platforms possible in the future (iOS, tablet, web)
Architecture:
[Android App] <--HTTPS/JSON--> [Golang API] <---> [SQLite]
                                         |
                                         --> [Google Cloud Sync / CSV / PDF Export]
⚡ Key Advantages
Unified backend for web and mobile (future-proof)
Lightweight, efficient stack (Golang + SvelteKit)
Offline capabilities and optional cloud sync
Scalable design for future multi-user support
Modern, minimal, and responsive UI
🏁 Conclusion
Tonish combines MyFlow (task management) and MyFlowBook (digital notebook) into a single productivity platform with:
Efficient backend in Golang
Lightweight and reactive frontend in SvelteKit
Future mobile support through API-first architecture
CSV/PDF export and Google Cloud sync
This structure ensures flexibility, scalability, and ease of future development.