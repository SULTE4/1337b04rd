<div align="center">

# 1337b04rd

**Post a thread. No account. It's gone in 10 minutes unless people reply.**

An anonymous imageboard in Go with hexagonal architecture. Every visitor gets a Rick and Morty character as their avatar,<br>
and threads that go quiet move to a read-only archive on their own.

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Architecture](https://img.shields.io/badge/architecture-hexagonal-blueviolet)
![Rick and Morty API](https://img.shields.io/badge/avatars-Rick%20and%20Morty%20API-97ce4c)
![License: MIT](https://img.shields.io/badge/license-MIT-blue)

</div>

---

## Why

Classic imageboards are anonymous and ephemeral: no sign-up, and threads live only as long as people keep talking. 1337b04rd rebuilds that idea as a clean Go web app, with the business rules in a core that knows nothing about HTTP, PostgreSQL or where images are stored.

## How it works

```
 browser ──► net/http router ──► TokenMiddleware ──► handlers ──► html/template (embedded)
                                    │  no "token" cookie?
                                    ▼
                         new user: random Rick and Morty character
                         (rickandmortyapi.com, avoids ones already taken),
                         signed into a JWT cookie for 7 days
                                    │
              ┌─────────────────────┼──────────────────────┐
              ▼                     ▼                      ▼
        PostService           CommentService           UserService       core (ports)
              │                     │                      │
              ▼                     ▼                      ▼
        PostgreSQL: Users, Post, Comment          ./s3_storage/postsImg, commentsImg

 thread lifecycle:
   new post           expires = now + 10 min
   each new comment   expires = now + 15 min
   expires passed     leaves /catalog, readable at /archive, comments rejected
```

## Features

- **Anonymous.** No registration. Your identity is a signed cookie with a Rick and Morty name and picture.
- **Unique avatars.** Characters already taken by other users are skipped when a new one is picked.
- **Threads that expire.** A new thread lives 10 minutes. Every reply pushes that to 15 minutes from the reply.
- **Archive.** Expired threads stay readable at `/archive` but are locked for comments.
- **Images.** JPEG, PNG or GIF up to 5 MB on posts and comments, stored in S3-style buckets on disk.
- **Replies.** Comments can answer other comments.
- **Hexagonal architecture.** Services depend on ports, while PostgreSQL, the image store, the avatar API and HTTP are adapters.
- **Single binary.** Templates and CSS are embedded with `go:embed`. The only dependency is `lib/pq`.

## Quick start

**1. Start PostgreSQL** (Docker)

```bash
git clone https://github.com/SULTE4/1337b04rd.git
cd 1337b04rd
docker compose up -d
```

**2. Create the tables**

```bash
docker exec -i mydb_container psql -U postgres -d mydb < init.sql
```

**3. Run**

```bash
go run ./cmd/web --port 8080
```

Open <http://localhost:8080>. You are redirected to `/catalog`.

The connection string `postgres://postgres:pass@localhost:5432/mydb` is set in `cmd/web/main.go` and matches `docker-compose.yaml`. Uploaded images go to `./s3_storage/` next to where you start the server.

```
$ go run ./cmd/web --help
hacker board

Usage:
  1337b04rd [--port <N>]
  1337b04rd --help

Options:
  -help
    	Show this screen.
  -port string
    	Port number. (default "8080")
```

## Routes

| Method | Path | Page / action |
|---|---|---|
| `GET` | `/` | Redirects to `/catalog` |
| `GET` | `/catalog` | Live threads, soonest to expire first |
| `GET` | `/post/create` | New thread form |
| `GET` | `/post/{id}` | A thread with its comments |
| `GET` | `/archive` | Expired threads |
| `GET` | `/archive/post/{id}` | An archived thread (read-only) |
| `POST` | `/submit-post` | Create a thread (name, subject, text, image) |
| `POST` | `/submit-comment/{id}` | Comment on a thread, optionally as a reply |
| `GET` | `/postsImg/...`, `/commentsImg/...` | Uploaded images |

## Data

```
Users    (userID, username, userURL, userToken)
Post     (id, title, content, imageURL, created, expires, author, userID -> Users)
Comment  (comment_id, content, imageURL, created, userID -> Users,
          postID -> Post, parentID -> Comment)
```

## Development

```bash
go build ./...
go vet ./...
```

```
1337b04rd/
├── cmd/web/main.go                 # flags, DSN, HTTP server
├── internal/
│   ├── core/
│   │   ├── domain/                 # Post, Comment, User, File
│   │   ├── ports/                  # repository, service and API interfaces
│   │   ├── service/                # post, comment and user rules
│   │   └── util/                   # JWT, validation
│   ├── adapters/
│   │   ├── interface/http/         # handlers, token middleware, templates
│   │   ├── interface/router/       # routes
│   │   ├── postgres/               # repositories
│   │   ├── s3/                     # on-disk image buckets
│   │   └── external/               # Rick and Morty API client
│   └── app/                        # wiring
├── web/                            # embedded templates and CSS
├── init.sql                        # schema
└── docker-compose.yaml             # PostgreSQL 16
```

## Roadmap

- [ ] Read the database URL, port and JWT secret from environment variables
- [ ] Fix the tests in `test/` (they don't build against the current service constructors)
- [ ] Real S3 / MinIO adapter behind the existing port
- [ ] Run the app itself in Docker Compose (a `Dockerfile` is already there)

## Author

Sultanbek ([@SULTE4](https://github.com/SULTE4)), Astana, Kazakhstan.

## License

[MIT](LICENSE)
