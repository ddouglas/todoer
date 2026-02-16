# List available commands
default:
    @just --list

# Run database migrations
migrate:
    atlas -c file://migrations/atlas.hcl --env local schema apply --auto-approve

# Check migration plan (dry run)
migrate-plan:
    atlas -c file://migrations/atlas.hcl --env local schema apply --dry-run

# Start the HTTP server
serve:
    go run cmd/todoer/*.go -c config.yml serve

# Start the notification worker
worker:
    go run cmd/todoer/*.go -c config.yml worker

# Build the binary
build:
    go build -o bin/todoer cmd/todoer/*.go

# Run both server and worker (requires tmux or run in separate terminals)
dev:
    @echo "Start server: just serve"
    @echo "Start worker: just worker"

# Clean build artifacts
clean:
    rm -rf bin/

# Show Atlas schema status
migrate-status:
    atlas -c file://migrations/atlas.hcl --env local schema inspect

# Format Go code
fmt:
    go fmt ./...

# Run tests
test:
    go test ./...

# Install dependencies
deps:
    go mod download
    go mod tidy