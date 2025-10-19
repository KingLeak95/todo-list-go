# Todo List API

A modern, containerized Todo List API built with Go, Gin, and PostgreSQL. Features comprehensive testing, CI/CD pipeline, and Docker support.

## 🚀 Features

- **RESTful API** with Gin framework
- **User Management** - Create, list, and delete users
- **Task Management** - Create, complete, and delete tasks
- **Database Integration** - PostgreSQL with GORM ORM
- **Comprehensive Testing** - Unit and integration tests
- **Docker Support** - Multi-stage builds for production and testing
- **CI/CD Pipeline** - Automated testing and building
- **Modern Go** - Built with Go 1.23

## 📋 API Endpoints

### Users
- `POST /createUser` - Create a new user
- `GET /allUsers` - List all users
- `DELETE /deleteUser/:id` - Delete a user

### Tasks
- `POST /tasks` - Create a new task
- `PUT /tasks/:id/complete` - Mark task as complete
- `DELETE /tasks/:id` - Delete a task

## 🛠️ Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL (for production)
- Make (optional, for convenience commands)

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/KingLeak95/todo-list-go.git
   cd todo-list-go
   ```

2. **Start with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Test the API**
   ```bash
   curl http://localhost:8080/
   ```

### Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   # Using Docker
   docker run -d --name postgres-todolist \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=todolist \
     -p 5432:5432 \
     postgres:15
   ```

3. **Set environment variables**
   ```bash
   export DB_HOST=localhost
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=todolist
   export DB_PORT=5432
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

## 🧪 Testing

### Run Tests Locally
```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run Dockerized tests
make docker-test
```

### Test API Endpoints

**Create a user:**
```bash
curl -X POST http://localhost:8080/createUser \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

**Create a task:**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"task":"Buy groceries","userId":1}'
```

**Complete a task:**
```bash
curl -X PUT http://localhost:8080/tasks/1/complete
```

## 🐳 Docker

### Build and Run

```bash
# Build the application image
make docker-build

# Run with Docker
make docker-start

# Run tests in Docker
make docker-test
```

### Docker Compose

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=todolist
    depends_on:
      - postgres

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=todolist
    ports:
      - "5432:5432"
```

## 🏗️ Development

### Project Structure
```
.
├── .github/workflows/    # CI/CD pipelines
├── models/               # Data models and handlers
│   ├── setup.go         # Database connection
│   ├── user.go          # User model and handlers
│   ├── tasks.go         # Task model and handlers
│   └── models_test.go   # Model tests
├── main.go              # Application entry point
├── main_test.go         # Integration tests
├── Dockerfile           # Production image
├── Dockerfile.test      # Testing image
├── Makefile            # Build automation
└── README.md           # This file
```

### Available Make Commands
```bash
make help              # Show available commands
make test              # Run tests
make build             # Build binary
make run               # Build and run
make docker-build      # Build Docker image
make docker-test       # Run tests in Docker
make docker-compose-up # Start with Docker Compose
make docker-compose-down # Stop Docker Compose
make docker-compose-logs # View Docker Compose logs
make clean             # Clean up
```

## 🔧 Configuration

### Environment Variables
- `DB_HOST` - Database host (default: localhost)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: todolist)
- `DB_PORT` - Database port (default: 5432)
- `GIN_MODE` - Gin mode (default: debug, set to release for production)

## 🚀 Deployment

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todo-list-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: todo-list-api
  template:
    metadata:
      labels:
        app: todo-list-api
    spec:
      containers:
      - name: todo-list-api
        image: todo-list:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
```

### Production Considerations
- Use environment-specific configuration
- Set up proper logging
- Configure health checks
- Set up monitoring and alerting
- Use secrets management for database credentials

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🐛 Known Issues

- SQLite driver requires CGO for testing (resolved with Debian-based test image)
- Database migrations are handled automatically by GORM

## 🔄 CI/CD Pipeline

The project includes a comprehensive CI/CD pipeline that:
- Runs unit tests with coverage
- Executes Dockerized tests
- Builds application artifacts
- Tests Docker images
- Supports multiple trigger events (push, PR, release)

## 📊 Project Status

- ✅ User CRUD operations
- ✅ Task CRUD operations  
- ✅ Database integration
- ✅ Comprehensive testing
- ✅ Docker support
- ✅ CI/CD pipeline
- ✅ Documentation 
