# Todoer - Project State & Instructions

**Last Updated:** February 16, 2026

## Overview
A Go + HTMX todo list application with Postgres backend.

## Tech Stack
- **Backend:** Go 1.21+
- **Router:** alexedwards/flow
- **Database:** Postgres (via pgxpool + scany)
- **Query Builder:** Masterminds/squirrel
- **Migrations:** Atlas (HCL schema-as-code)
- **Logging:** logrus
- **Frontend:** HTMX 2.0.8 + Pico CSS v2
- **Templates:** Go html/template (embedded FS)

## Project Structure
```
cmd/todoer/          # main.go - entry point
internal/
  server/            # HTTP server, handlers, middleware
    templates/       # embedded HTML templates
      components/    # reusable UI components
      pages/         # full page templates
      layout.html    # base layout
  store/             # database repository layer
pkg/
  types/             # shared types (models, config)
migrations/          # Atlas schema definitions
```

## Architecture Decisions

### Layout & Conventions
- **Main binary:** `cmd/todoer/main.go`
- **Business logic:** `internal/`
- **Shared types:** `pkg/types/`
- **Config:** Environment variables via `kelseyhightower/envconfig`

### Database
- **Local dev:** Postgres in Docker Compose
- **Connection:** pgxpool for connection pooling
- **Queries:** squirrel for query building, scany for row scanning
- **Migrations:** Atlas with HCL schema files in `migrations/`

### Templates
- **Pattern:** Component-based (Pattern B)
- **Naming:** 
  - `layout.base` - base HTML shell
  - `component.*` - reusable components
  - `page.*` - full pages
- **Embedded:** Templates are embedded in binary via `embed.FS`
- **Loading:** All templates parsed recursively at startup in `server.New()`

### Template Components
- `layout.base` - HTML shell with flexbox layout (nav + sidebar + main)
- `component.nav` - Navigation header
- `component.sidebar` - Sidebar with filters and categories
- `component.todo.form` - Quick add todo form (simple)
- `component.todo.detail-form` - Full todo form with all metadata
- `component.todo.list` - Todo list container
- `component.todo.item` - Enhanced todo row with priority, category, due date
- `component.todo.empty` - Empty state message
- `component.todo.empty-oob` - Empty state with OOB swap (restore on delete)
- `component.todo.empty-delete` - OOB delete command (remove on create)
- `component.modal` - Modal wrapper for HTMX
- `component.modal.todo-detail` - Todo detail modal content
- `component.modal.category-form` - Category form modal
- `page.home` - Home page with sidebar and filtered list
- `page.todo-detail` - Dedicated todo detail/edit page

### Handlers
- **Page data:** Typed structs per page (`HomePageData`, `TodoDetailPageData`)
- **Handler organization:** Separate files for todos (`todo.go`), categories (`category.go`), pages (`page.go`)
- **HTMX responses:** Render component fragments for dynamic updates
- **Modal handlers:** Separate endpoints for modal content vs full pages
- **Error handling:** Log with logrus, return HTTP errors
- **Content-Type:** Always set `text/html` for template responses

### UI Layout Architecture
- **Flexbox layout:** Container with fixed nav, flexible sidebar + main content
- **Sidebar:** 250px fixed width, filters + categories list
- **Main content:** Flexible width, scrollable
- **Responsive:** Sidebar and content adapt to screen size
- **Modals:** HTMX-driven overlays for quick actions

### HTMX Patterns
- **Create (simple):** POST `/todos` → append `component.todo.item` + OOB delete empty state
- **Create (full):** POST `/todos` from detail form → redirect or close modal
- **Toggle:** PATCH `/todos/{id}` → swap `component.todo.item` in place
- **Update:** PUT `/todos/{id}` → swap updated `component.todo.item`
- **Delete:** DELETE `/todos/{id}` → remove row, restore empty state if last item
- **Modals:** GET `/todos/{id}/modal` → load `component.modal.todo-detail` into modal container
- **Filters:** GET `/?filter=today` → full page reload with filtered todos
- **OOB swaps:** Used for empty state management (adding/removing)

### Logging
- **Library:** logrus (switched from slog for cleaner syntax)
- **Format:** JSON for production, could add text for local dev
- **Fatal errors:** `logger.WithError(err).Fatal()` for startup failures
- **Request logging:** Middleware logs method, path, status, duration

## Current State - What's Working ✓

### Features Implemented
- ✅ Full CRUD operations with enhanced metadata (Create, Read, Update, Delete)
- ✅ Real-time updates via HTMX (no page reloads)
- ✅ Sidebar navigation with filters (All, Today, This Week)
- ✅ Category management (create, view, filter by category)
- ✅ Todo priorities (low, medium, high) with visual indicators
- ✅ Due dates with date picker
- ✅ Todo descriptions for additional context
- ✅ Reminders (timestamp field)
- ✅ Toggle todo completion (checkbox)
- ✅ Delete todos with confirmation
- ✅ Empty state management (shows/hides based on todo count)
- ✅ Modal support via HTMX for quick add/edit
- ✅ Dedicated detail pages for full todo editing
- ✅ Responsive sidebar layout with Pico CSS
- ✅ Category color coding and emoji icons
- ✅ Database migrations with Atlas
- ✅ Docker Compose for local Postgres

### Database Schema
**Table:** `categories`
- `id` - text (primary key)
- `name` - text (not null)
- `color` - text (not null, default: #3b82f6)
- `icon` - text (nullable, for emoji)
- `sort_order` - integer (not null, default: 0)
- `created_at` - timestamp (default now())

**Table:** `todos`
- `id` - text (primary key)
- `title` - text (not null)
- `description` - text (nullable)
- `completed` - boolean (default false)
- `priority` - text (not null, default: medium) [low, medium, high]
- `category_id` - text (nullable, foreign key to categories)
- `due_date` - timestamp (nullable)
- `reminder_at` - timestamp (nullable)
- `created_at` - timestamp (default now())
- `updated_at` - timestamp (default now())

**Indexes:**
- `todos.category_id` for fast category filtering
- `todos.due_date` for date-based queries
- `todos.completed` for filtering active/completed

### Environment Variables
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=todoer
DB_PASSWORD=password
DB_NAME=todoer
DB_SSLMODE=disable
PORT=8080
```

## Routes / API Endpoints

### Pages
- `GET /` - Home page with sidebar, shows all todos or filtered view
  - Query params: `?filter=all|today|week|category&category={id}`
- `GET /todos/new` - Full page for creating new todo
- `GET /todos/{id}` - Full page for viewing/editing todo

### Todo Operations (HTMX)
- `POST /todos` - Create new todo, returns `component.todo.item`
- `PUT /todos/{id}` - Update todo (full update), returns `component.todo.item`
- `PATCH /todos/{id}` - Partial update (toggle complete), returns `component.todo.item`
- `DELETE /todos/{id}` - Delete todo, returns empty or OOB empty state

### Modal Endpoints (HTMX)
- `GET /todos/new/modal` - Returns modal with todo creation form
- `GET /todos/{id}/modal` - Returns modal with todo edit form

### Category Operations
- `GET /categories` - List categories (as HTMX fragments)
- `POST /categories` - Create category, redirect to home
- `GET /categories/new` - Returns modal with category creation form
- `PUT /categories/{id}` - Update category, redirect to home
- `DELETE /categories/{id}` - Delete category, redirect to home

## Known Issues / Quirks
1. **HTMX OOB swaps:** Must wrap elements in a container div when using `hx-swap-oob` to include the outer element itself
2. **Checkbox state:** Unchecked checkboxes send no value; checked sends `"on"`
3. **Template hot reload:** Templates are embedded, so require app restart to see changes
4. **Modal close:** Currently uses JavaScript + page reload; could be improved with HTMX events

## Development Commands

### Start local database
```bash
docker compose up -d
```

### Run migrations
```bash
atlas -c file://migrations/atlas.hcl --env local schema apply
```

### Generate new migration (after editing schema)
```bash
atlas migrate diff \
  --dev-url "postgres://todoer:password@localhost:5432/todoer?sslmode=disable" \
  --to "file://migrations/todos.pg.hcl" \
  --dir "file://migrations"
```

### Run application
```bash
go run cmd/todoer/*.go
```

### Install dependencies
```bash
go mod tidy
```

## Next Steps / TODOs

### High Priority
- [ ] Test redesigned UI and fix any bugs
- [ ] Run migrations on fresh database
- [ ] Verify all CRUD operations work with new schema
- [ ] Add sample categories for demo
- [ ] Deployment preparation (platform TBD)
  - Containerization strategy
  - Database hosting decisions
  - Environment configuration
  - CI/CD setup

### Features to Add
- [ ] User authentication (sign in/sign up)
- [ ] User-specific todos and categories (multi-user support)
- [ ] Inline editing for todo title
- [ ] Search/filter todos by text
- [ ] Sort todos (by date, priority, completion, manual drag-drop)
- [ ] Recurring tasks
- [ ] Todo tags/labels (in addition to categories)
- [ ] Bulk operations (mark multiple as complete, delete, move category)
- [ ] Subtasks / checklists within todos
- [ ] File attachments

### UI/UX Improvements
- [ ] Better modal close behavior (use HTMX events instead of reload)
- [ ] Loading states for HTMX requests
- [ ] Toast notifications for actions (created, updated, deleted)
- [ ] Undo delete (with toast notification)
- [ ] Keyboard shortcuts (n for new, / for search, etc.)
- [ ] Drag and drop reordering
- [ ] Dark mode toggle
- [ ] Mobile-responsive sidebar (hamburger menu)
- [ ] Overdue indicator styling (red/bold)
- [ ] Progress indicators (X of Y todos complete per category)

### Code Quality
- [ ] Input validation (max length, sanitization)
- [ ] Error messages to user (not just server errors)
- [ ] Add unit tests (handlers, store layer)
- [ ] Integration tests (full CRUD flow)
- [ ] Add request timeouts
- [ ] Rate limiting
- [ ] CSRF protection
- [ ] SQL injection prevention verification

### DevOps
- [ ] Health check endpoint (`/health`)
- [ ] Dockerfile for production builds
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Structured logging with request IDs
- [ ] Metrics/monitoring
- [ ] Database connection pool tuning

## Deployment Prep

### Build Command
```bash
go build -o bin/todoer cmd/todoer/main.go
```

### Start Command
```bash
./bin/todoer
```

### Required Environment Variables
- `DB_HOST` - Database host
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSLMODE` - SSL mode (disable for local, require for production)
- `PORT` - HTTP server port (default: 8080)

### Pre-deployment Checklist
- [ ] Ensure proper `PORT` env var handling
- [ ] Test with `sslmode=require` for production database
- [ ] Set up migration strategy (manual, init container, or startup script)
- [ ] Set logrus to JSON output for production
- [ ] Configure database connection pooling for production load

## Tips & Patterns

### Adding a New Page
1. Create page data struct in `pkg/types/`
2. Create page template in `internal/server/templates/pages/`
3. Create handler in `internal/server/page.go` or new file
4. Register route in `buildRouter()`
5. Update nav component if needed

### Adding a New Component
1. Create template in `internal/server/templates/components/`
2. Use naming: `component.<feature>.<name>`
3. Define clear data requirements (document in comments)
4. Test with both full page and HTMX contexts

### Database Query Pattern
```go
// 1. Define method on repository
func (r *TodoRepository) MethodName(ctx context.Context, ...) (*types.Todo, error) {
    query := squirrel.Select("*").From("todos").Where(...)
    sql, args, _ := query.PlaceholderFormat(squirrel.Dollar).ToSql()
    
    var todo types.Todo
    err := pgxscan.Get(ctx, r.pool, &todo, sql, args...)
    return &todo, err
}

// 2. Call from handler
todo, err := s.todos.MethodName(r.Context(), ...)
if err != nil {
    s.logger.WithError(err).Error("failed to...")
    s.internalServerError(w)
    return
}
```

### Template Rendering Pattern
```go
// Full page
data := &types.PageData{...}
w.Header().Set("Content-Type", "text/html")
s.templates.ExecuteTemplate(w, "page.name", data)

// HTMX fragment
w.Header().Set("Content-Type", "text/html")
s.templates.ExecuteTemplate(w, "component.name", item)
```

## Resources
- [HTMX Documentation](https://htmx.org/docs/)
- [Pico CSS](https://picocss.com/)
- [Flow Router](https://github.com/alexedwards/flow)
- [Atlas CLI](https://atlasgo.io/)
- [pgx](https://github.com/jackc/pgx)

---

**Remember:** This is advice mode unless you explicitly ask for code generation. Keep momentum—ship features, polish later.
