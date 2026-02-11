# 1337b04rd - Anonymous Imageboard

A modern anonymous imageboard backend built with Go, implementing hexagonal architecture, RESTful API design, and real-time content management.

## Overview

1337b04rd is an anonymous textboard/imageboard that allows users to create posts with images, comment on threads, and interact without registration. The application features automatic thread archival, unique user avatars from Rick and Morty API, and S3-compatible image storage.

## Features

- **Anonymous Posting**: No user registration required - users identified via secure session cookies
- **Image Support**: Upload images to posts stored in S3-compatible storage
- **Thread Management**: Automatic deletion of inactive threads (10-15 minutes)
- **Archive System**: Preserved threads remain viewable but locked from new comments
- **Unique Avatars**: Each user gets a unique Rick and Morty character avatar
- **Real-time Expiration**: Background workers manage thread lifecycle
- **Reply System**: Thread-based commenting with reply-to functionality

## Architecture

The project follows **Hexagonal Architecture** (Ports and Adapters pattern) to maintain clean separation of concerns:

```
┌─────────────────────────────────────┐
│     Presentation Layer (HTTP)       │
│   - Handlers                        │
│   - Middleware (Auth, Logging)      │
│   - Template Rendering              │
└───────────────┬─────────────────────┘
                │
┌───────────────▼─────────────────────┐
│        Domain Layer (Core)          │
│   - Business Logic                  │
│   - Ports (Interfaces)              │
│   - Models & Entities               │
└───────────────┬─────────────────────┘
                │
┌───────────────▼─────────────────────┐
│     Infrastructure Layer            │
│   - PostgreSQL Adapter              │
│   - S3 Storage Adapter              │
│   - Rick & Morty API Client         │
│   - Session Management              │
└─────────────────────────────────────┘
```

This architecture ensures:
- **Maintainability**: Easy to modify or replace infrastructure components
- **Testability**: Core business logic can be tested independently
- **Flexibility**: New features can be added without affecting existing code

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL
- **Storage**: S3-compatible (MinIO/custom implementation)
- **External API**: Rick and Morty API
- **Logging**: log/slog
- **Testing**: Go testing package with comprehensive coverage

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- S3-compatible storage (MinIO or custom implementation)
- Internet connection (for Rick and Morty API)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/SULTE4/1337b04rd.git
cd 1337b04rd
```

2. Set up environment variables:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/1337b04rd"
export S3_ENDPOINT="http://localhost:9000"
export S3_ACCESS_KEY="your-access-key"
export S3_SECRET_KEY="your-secret-key"
export S3_BUCKET_POSTS="posts-bucket"
export S3_BUCKET_COMMENTS="comments-bucket"
```

3. Initialize the database:
```bash
psql -U postgres -d 1337b04rd -f schema.sql
```

4. Build the application:
```bash
go build -o 1337b04rd ./cmd/1337b04rd
```

## Usage

Start the server:
```bash
./1337b04rd --port 8080
```

View help:
```bash
./1337b04rd --help
```

Output:
```
hacker board

Usage:
  1337b04rd [--port <N>]  
  1337b04rd --help

Options:
  --help       Show this screen.
  --port N     Port number.
```

Access the application at `http://localhost:8080`

## Project Structure

```
1337b04rd/
├── cmd/
│   └── 1337b04rd/
│       └── main.go           # Application entry point
├── internal/
│   ├── domain/
│   │   ├── models/           # Domain entities
│   │   ├── ports/            # Interface definitions
│   │   └── services/         # Business logic
│   ├── adapters/
│   │   ├── http/             # HTTP handlers & middleware
│   │   ├── postgres/         # Database adapter
│   │   ├── s3/               # S3 storage adapter
│   │   └── rickmorty/        # External API client
│   └── infrastructure/
│       ├── config/           # Configuration
│       └── logger/           # Logging setup
├── web/
│   ├── templates/            # HTML templates
│   └── static/               # CSS, JS, images
├── migrations/               # Database migrations
└── tests/                    # Integration & unit tests
```

## API Endpoints

### Pages
- `GET /` - Main catalog page
- `GET /archive` - Archive page
- `GET /post/:id` - View post and comments
- `GET /archive/post/:id` - View archived post
- `GET /create` - Create new post page

### API
- `POST /api/post` - Create new post
- `POST /api/comment` - Add comment to post
- `POST /api/user/name` - Update user name
- `GET /api/avatars` - Get available avatars

## Key Implementation Details

### Session Management
Sessions are managed using secure HTTP-only cookies with 1-week expiration. Each session is identified by a UUID and stored in PostgreSQL, providing persistent session storage across server restarts.

### Avatar Assignment
On a user's first visit, the application fetches a random character from the Rick and Morty API. Avatars are unique per user per thread, and the system intelligently reuses avatars when the pool is exhausted. Character data is cached locally to minimize external API calls.

### Thread Lifecycle
The application uses background goroutines with tickers to manage thread expiration:
- New posts without comments are automatically deleted after 10 minutes
- Posts with comments are deleted 15 minutes after the last comment
- Deleted threads are moved to the archive where they remain viewable but locked

### Image Handling
Images are validated for type (JPEG, PNG, GIF) and stored in separate S3 buckets for posts and comments. Each image gets a unique UUID-based key and is served via signed URLs for security.

## Testing

Run all tests:
```bash
go test ./...
```

Run with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

The project maintains comprehensive test coverage across all layers, with unit tests for business logic and integration tests for adapters.

## Configuration

### Environment Variables
- `DATABASE_URL` - PostgreSQL connection string
- `S3_ENDPOINT` - S3 server endpoint
- `S3_ACCESS_KEY` - S3 access key
- `S3_SECRET_KEY` - S3 secret key
- `S3_BUCKET_POSTS` - Bucket for post images
- `S3_BUCKET_COMMENTS` - Bucket for comment images
- `RICKMORTY_API_URL` - Rick and Morty API URL (default: https://rickandmortyapi.com/api)
- `SESSION_SECRET` - Secret key for session cookies
- `LOG_LEVEL` - Logging level (debug, info, warn, error)

## Development

### Code Standards
The codebase follows strict Go conventions:
- Formatted with [gofumpt](https://github.com/mvdan/gofumpt)
- Zero external dependencies except PostgreSQL driver
- No panics or unexpected runtime errors
- Comprehensive error handling

### Logging
Structured logging using `log/slog`:
```go
slog.Info("post created", "post_id", postID, "user_session", sessionID)
slog.Error("failed to store image", "error", err, "image_key", key)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Resources

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Go net/http Package](https://pkg.go.dev/net/http)
- [Rick and Morty API](https://rickandmortyapi.com/)
- [HTTP Cookies RFC](https://datatracker.ietf.org/doc/html/rfc6265)

## Author

Sultanbek
- GitHub: [@SULTE4](https://github.com/SULTE4)
- Location: Astana, Kazakhstan

---

Built with ❤️ using Go