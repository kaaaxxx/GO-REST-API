# Students API

A lightweight RESTful API for managing student records, built with Go and SQLite. The API provides endpoints to create, retrieve, and list student information with built-in validation and CORS support.

## Features

- **Create Students** - Add new student records with automatic validation
- **Retrieve Students** - Fetch individual students by ID or get a complete list
- **SQLite Storage** - Lightweight, file-based database for persistence
- **Input Validation** - Automatic validation of required fields (name, email, age)
- **CORS Enabled** - Ready for frontend integration
- **Graceful Shutdown** - Proper server cleanup and signal handling
- **Structured Logging** - Built-in logging with slog for better debugging

## Prerequisites

- Go 1.26.3 or higher
- SQLite3 (included with Go driver)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/kaaaxxx/students-api.git
cd students-api
```

2. Install dependencies:
```bash
go mod download
```

3. Configure the application (optional):
   - Edit `config/config.yaml` to customize the server address and database path
   - By default, the server runs on `localhost:8082`

## Running the Server

```bash
go run cmd/students-api/main.go
```

You should see:
```
Server started Address=localhost:8082
```

The server will gracefully shutdown on `CTRL+C`.

## API Endpoints

### Create a Student
**POST** `/api/students`

Request body:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

Response (201 Created):
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

### Get Student by ID
**GET** `/api/students/{id}`

Response (200 OK):
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

### List All Students
**GET** `/api/students`

Response (200 OK):
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 20
  },
  {
    "id": 2,
    "name": "Jane Smith",
    "email": "jane@example.com",
    "age": 21
  }
]
```

## Project Structure

```
student-api/
├── cmd/
│   └── students-api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loader
│   ├── http/
│   │   └── handlers/
│   │       └── student/
│   │           └── student.go   # Student request handlers
│   ├── storage/
│   │   ├── storage.go           # Storage interface
│   │   └── sqlite/
│   │       └── sqlite.go        # SQLite implementation
│   ├── types/
│   │   └── types.go             # Data structures (Student)
│   └── utils/
│       └── response/
│           └── response.go      # Response formatting utilities
├── config/
│   └── config.yaml              # Configuration file
├── frontend/
│   └── index.html               # Simple web interface
└── storage/                     # Database storage directory
```

## Configuration

Edit `config/config.yaml` to customize behavior:

```yaml
env: "dev"                              # Environment (dev/prod)
storage_path: "storage/storage.db"      # SQLite database path
http_server:
  address: "localhost:8082"             # Server address and port
```

## Development

### Dependencies

The project uses:
- **go-playground/validator** - Request validation
- **ilyakaznacheev/cleanenv** - Configuration management
- **mattn/go-sqlite3** - SQLite driver

### Key Concepts

- **CORS Middleware** - Allows requests from any origin (all GET, POST methods)
- **Structured Logging** - Uses Go's slog for better observability
- **Graceful Shutdown** - Responds to interrupt signals and cleans up resources

## Future Enhancements

- Update and delete endpoints
- Database migrations
- Authentication/Authorization
- Pagination for list endpoint
- More comprehensive error handling
- Unit tests

## License

This project is open source and available under the MIT License.
