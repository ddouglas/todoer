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
- `layout.base` - HTML shell with nav + main container
- `component.nav` - Navigation header
- `component.todo.form` - Add todo form
- `component.todo.list` - Todo list container
- `component.todo.item` - Single todo row
- `component.todo.empty` - Empty state message
- `component.todo.empty-oob` - Empty state with OOB swap (restore on delete)
- `component.todo.empty-delete` - OOB delete command (remove on create)
- `page.home` - Home page composition

### Handlers
- **Page data:** Single struct per page (`HomePageData`)
- **HTMX responses:** Render component fragments
- **Error handling:** Log with logrus, return HTTP errors
- **Content-Type:** Always set `text/html` for template responses

### HTMX Patterns
- **Create:** POST `/todos` → append `component.todo.item` + OOB delete empty state
- **Toggle:** PATCH `/todos/{id}` → swap `component.todo.item` in place
- **Delete:** DELETE `/todos/{id}` → remove row, restore empty state if last item
- **OOB swaps:** Used for empty state management (adding/removing)

### Logging
- **Library:** logrus (switched from slog for cleaner syntax)
- **Format:** JSON for production, could add text for local dev
- **Fatal errors:** `logger.WithError(err).Fatal()` for startup failures
- **Request logging:** Middleware logs method, path, status, duration

## Current State - What's Working ✓

### Features Implemented
- ✅ Basic CRUD operations (Create, Read, Update, Delete todos)
- ✅ Real-time updates via HTMX (no page reloads)
- ✅ Toggle todo completion (checkbox)
- ✅ Delete individual todos
- ✅ Empty state management (shows/hides based on todo count)
- ✅ Responsive layout with Pico CSS
- ✅ Database migrations with Atlas
- ✅ Docker Compose for local Postgres

### Database Schema
**Table:** `todos`
- `id` - serial (primary key)
- `title` - text (not null)
- `completed` - boolean (default false)
- `created_at` - timestamp (default now())

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

## Known Issues / Quirks
1. **HTMX OOB swaps:** Must wrap elements in a container div when using `hx-swap-oob` to include the outer element itself
2. **Checkbox state:** Unchecked checkboxes send no value; checked sends `"on"`
3. **Template hot reload:** Templates are embedded, so require app restart to see changes

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
- [ ] Deployment preparation (platform TBD)
  - Containerization strategy
  - Database hosting decisions
  - Environment configuration
  - CI/CD setup

### Features to Add
- [ ] User authentication (sign in/sign up)
- [ ] User-specific todos (todos belong to users)
- [ ] Edit todo title (inline editing)
- [ ] Todo priority/categories
- [ ] Due dates
- [ ] Search/filter todos
- [ ] Sort todos (by date, priority, completion)

### Improvements
- [ ] Input validation (max length, sanitization)
- [ ] Error messages to user (not just server errors)
- [ ] Loading states for HTMX requests
- [ ] Undo delete (with toast notification)
- [ ] Keyboard shortcuts
- [ ] Better empty state (CTA button to focus input)
- [ ] Add timestamps display (created/updated)

### Code Quality
- [ ] Add unit tests (handlers, store layer)
- [ ] Integration tests (full CRUD flow)
- [ ] Add request timeouts
- [ ] Rate limiting
- [ ] CSRF protection
- [ ] Input sanitization
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
