# Go Microservices Backend

A portfolio/project management backend built with Go microservices architecture, Clean Architecture, and gRPC.

[View OpenAPI Specification](bff-gateway/api/openapi.yaml)

## Services

| Service | Port | Description |
|---------|------|-------------|
| **BFF Gateway** | 8080 | REST API Gateway |
| **Auth Service** | 50051 | User authentication, roles, project access |
| **Project Service** | 50052 | Project management, skills, tech stack, images, links |
| **Task Service** | 50053 | Task management, subtasks, comments, attachments, tags |
| **Analytics Service** | 50054 | Project views, task activities, statistics |
| **Media Service** | 50055 | File upload and management |

## Project Structure

```
goBackend/
├── bff-gateway/                # REST API Gateway
│   ├── cmd/main.go
│   └── internal/
│       ├── config/
│       ├── grpc/
│       ├── handler/
│       ├── middleware/
│       └── router/
├── proto/                      # Protobuf definitions
│   ├── auth/
│   ├── project/
│   ├── task/
│   ├── analytics/
│   └── media/
├── shared/                     # Shared libraries
│   ├── database/
│   ├── middleware/
│   └── jwt/
├── services/                   # Microservices
│   ├── auth-service/
│   ├── project-service/
│   ├── task-service/
│   ├── analytics-service/
│   └── media-service/
├── migrations/                 # Database migrations
├── docker-compose.yml
├── Makefile
└── go.work
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Protocol Buffers compiler (protoc)

### With Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

---

## REST API Endpoints (BFF Gateway - Port 8080)

### 🔐 Auth (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login |
| POST | `/api/auth/validate` | Validate token |

**Request Examples:**

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"pass123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass123"}'
```

---

### 👤 Auth (Protected - Requires Token)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/profile` | Get current user profile |

---

### 👥 Users (Admin Only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users` | List all users |
| GET | `/api/users/:id` | Get user by ID |
| PUT | `/api/users/:id` | Update user |
| DELETE | `/api/users/:id` | Delete user |

---

### 📂 Projects

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects` | Create project |
| GET | `/api/projects` | List projects |
| GET | `/api/projects/:id` | Get project |
| PUT | `/api/projects/:id` | Update project |
| DELETE | `/api/projects/:id` | Delete project |
| POST | `/api/projects/:id/skills` | Add skill to project |
| POST | `/api/projects/:id/tech` | Add tech stack |
| POST | `/api/projects/:id/images` | Add image |
| POST | `/api/projects/:id/links` | Add link |

**Query Parameters (GET /api/projects):**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)
- `status` - Filter by status (active/completed/archived)

---

### 🏷️ Skills

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/skills` | List all skills |
| POST | `/api/skills` | Create skill |

---

### ✅ Tasks

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/tasks` | Create task |
| GET | `/api/tasks` | List tasks |
| GET | `/api/tasks/:id` | Get task |
| PUT | `/api/tasks/:id` | Update task |
| DELETE | `/api/tasks/:id` | Delete task |

**Query Parameters (GET /api/tasks):**
- `project_id` - Filter by project
- `page` - Page number
- `limit` - Items per page
- `status` - Filter by status (Todo/InProgress/Done)
- `assigned_to` - Filter by assigned user ID

---

### 📝 Subtasks

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/tasks/:id/subtasks` | Create subtask |
| GET | `/api/tasks/:id/subtasks` | List subtasks |

---

### 💬 Comments

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/tasks/:id/comments` | Add comment |
| GET | `/api/tasks/:id/comments` | List comments |

---

### 📎 Attachments

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/tasks/:id/attachments` | Add attachment |
| GET | `/api/tasks/:id/attachments` | List attachments |

---

### 🏷️ Tags

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/tags` | List all tags |
| POST | `/api/tags` | Create tag |
| POST | `/api/tasks/:id/tags` | Add tag to task |

---

### 📊 Analytics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/analytics/dashboard` | Get dashboard stats |
| POST | `/api/analytics/projects/:id/view` | Record project view |
| GET | `/api/analytics/projects/:id/views` | Get project views |
| GET | `/api/analytics/projects/:id/stats` | Get project stats |
| POST | `/api/analytics/tasks/:id/activity` | Record task activity |
| GET | `/api/analytics/tasks/:id/activities` | Get task activities |

---

### 📁 Media

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/media/upload` | Upload file (multipart/form-data) |
| GET | `/api/media` | List all files |
| GET | `/api/media/my-files` | List current user's files |
| GET | `/api/media/:id` | Get file |
| DELETE | `/api/media/:id` | Delete file |

**Upload Example:**

```bash
curl -X POST http://localhost:8080/api/media/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/file.jpg" \
  -F "file_type=image"
```

---

## Authentication

All protected endpoints require JWT token in header:

```
Authorization: Bearer <token>
```

---

## gRPC Services (Internal)

| Service | Port | Proto File |
|---------|------|------------|
| Auth | 50051 | `proto/auth/auth.proto` |
| Project | 50052 | `proto/project/project.proto` |
| Task | 50053 | `proto/task/task.proto` |
| Analytics | 50054 | `proto/analytics/analytics.proto` |
| Media | 50055 | `proto/media/media.proto` |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | 8080 | BFF Gateway port |
| `GRPC_PORT` | varies | gRPC server port |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | postgres | Database password |
| `DB_NAME` | portfolio | Database name |
| `JWT_SECRET` | (required) | JWT signing key |
| `STORAGE_PATH` | ./uploads | Media storage path |

---

## Development

```bash
# Generate protobuf files
make proto

# Build all services
make build

# Run tests
make test

# Install protobuf tools
make install-proto-tools
```

---

## API Summary

| Category | Endpoints |
|----------|-----------|
| Auth | 4 |
| Users | 4 |
| Projects | 9 |
| Skills | 2 |
| Tasks | 5 |
| Subtasks | 2 |
| Comments | 2 |
| Attachments | 2 |
| Tags | 3 |
| Analytics | 6 |
| Media | 5 |
| **Total** | **44 endpoints** |

---

## License

MIT
