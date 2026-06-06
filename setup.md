# Student API - Complete Setup & Project Documentation

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Project Overview](#project-overview)
3. [Project Structure](#project-structure)
4. [Setup Instructions](#setup-instructions)
5. [Core Components](#core-components)
   - [Datatypes](#datatypes)
   - [Functions & Handlers](#functions--handlers)
   - [Storage Layer](#storage-layer)
   - [Configuration](#configuration)
   - [Response Handling](#response-handling)
6. [Packages & Dependencies](#packages--dependencies)
7. [Git Commit History](#git-commit-history)

---

## Prerequisites

Before setting up the Student API backend, ensure you have the following installed:

- **Go** (version 1.26.3 or higher)
- **Git** (for version control)
- **SQLite3** (the project uses SQLite as the database)
- **Any code editor** (VS Code, GoLand, etc.)

---

## Project Overview

The Student API is a RESTful backend service built with Go that manages student records. It provides endpoints for:
- Creating new students
- Retrieving a specific student by ID
- Retrieving all students

The project uses:
- **Go's standard `net/http`** for HTTP server and routing
- **SQLite** for data persistence
- **YAML** configuration files for application settings
- **Validation** using the `go-playground/validator` package

---

## Project Structure

```
student-api/
├── cmd/
│   └── students-api/
│       └── main.go                 # Application entry point
├── config/
│   └── config.yaml                 # Configuration file (env, port, db path)
├── frontend/
│   ├── index.html                  # Frontend UI
│   └── plan.md                     # Frontend documentation
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration loading logic
│   ├── http/
│   │   └── handlers/
│   │       └── student/
│   │           └── student.go     # HTTP handlers for student endpoints
│   ├── storage/
│   │   ├── storage.go             # Storage interface definition
│   │   └── sqlite/
│   │       └── sqlite.go          # SQLite implementation
│   ├── types/
│   │   └── types.go               # Data structures/models
│   └── utils/
│       └── response/
│           └── response.go        # Response formatting utilities
├── storage/                        # Directory for SQLite database file
├── go.mod                          # Go module definition
├── go.sum                          # Go dependencies lock file
├── README.md                       # Project README
└── folder-structure.md             # Folder structure documentation
```

---

## Setup Instructions

### Step 1: Create the Project Directory Structure

```bash
# Create the main project directory
mkdir student-api
cd student-api

# Initialize as a git repository
git init

# Create directory structure
mkdir -p cmd/students-api
mkdir -p config
mkdir -p frontend
mkdir -p internal/config
mkdir -p internal/http/handlers/student
mkdir -p internal/storage/sqlite
mkdir -p internal/types
mkdir -p internal/utils/response
mkdir -p storage
```

### Step 2: Initialize Go Module

```bash
# Initialize Go module
go mod init github.com/kaaaxxx/students-api
```

This creates `go.mod` with initial content:
```
module github.com/kaaaxxx/students-api

go 1.26.3
```

### Step 3: Create Configuration File

Create `config/config.yaml`:
```yaml
env: "dev"
storage_path: "storage/storage.db"
http_server:
  address: "localhost:8082"
```

Also create `.gitignore` to exclude config from version control:
```
config/config.yaml
```

### Step 4: Install Dependencies

```bash
# The dependencies will be added as you progress through the commits
# Key packages to install:
go get github.com/ilyakaznacheev/cleanenv        # For config parsing
go get github.com/mattn/go-sqlite3                # For SQLite
go get github.com/go-playground/validator/v10    # For validation
```

### Step 5: Create Core Files (in order)

Follow the git commit history section below to implement files in the correct order.

---

## Core Components

### Datatypes

#### Student Struct
**File:** `internal/types/types.go`

```go
type Student struct {
    Id    int64  `json:"id"`                           // Unique identifier
    Name  string `json:"name" validate:"required"`     // Student name (required)
    Email string `json:"email" validate:"required"`    // Student email (required)
    Age   int    `json:"age" validate:"required"`      // Student age (required)
}
```

**Purpose:** Represents a student record with validation tags for request validation.

#### HTTPServer Struct
**File:** `internal/config/config.go`

```go
type HTTPServer struct {
    Addr string `yaml:"address" env-required:"true"`  // Server address (e.g., "localhost:8082")
}
```

**Purpose:** Holds HTTP server configuration.

#### Config Struct
**File:** `internal/config/config.go`

```go
type Config struct {
    Env         string `yaml:"env" env:"ENV" env-required:"true"`           // Environment (dev/prod)
    StoragePath string `yaml:"storage_path" env-required:"true"`            // Database file path
    HTTPServer  `yaml:"http_server"`                                         // Embedded HTTPServer config
}
```

**Purpose:** Main configuration structure for the application.

#### Sqlite Struct
**File:** `internal/storage/sqlite/sqlite.go`

```go
type Sqlite struct {
    Db *sql.DB  // Database connection pointer
}
```

**Purpose:** Implements the Storage interface and holds the database connection.

#### Response Struct
**File:** `internal/utils/response/response.go`

```go
type Response struct {
    Status string  // Response status ("OK" or "Error")
    Error  string  // Error message (if any)
}
```

**Purpose:** Standard response structure for API responses.

### Functions & Handlers

#### Handler Functions

**File:** `internal/http/handlers/student/student.go`

##### 1. **New(storage Storage.Storage) http.HandlerFunc**
- **Purpose:** Creates a new student
- **Method:** POST
- **Route:** `/api/students`
- **Input:** JSON body with name, email, age (all required)
- **Output:** JSON with created student ID
- **Validation:** Uses `go-playground/validator` to validate required fields
- **Process:**
  1. Decodes JSON request body
  2. Validates all fields are present
  3. Calls storage layer to insert student
  4. Returns created student ID

##### 2. **GetById(storage Storage.Storage) http.HandlerFunc**
- **Purpose:** Retrieves a student by ID
- **Method:** GET
- **Route:** `/api/students/{id}`
- **Input:** Student ID in URL path
- **Output:** JSON with student details
- **Process:**
  1. Extracts ID from URL path
  2. Parses ID to int64
  3. Calls storage layer to fetch student
  4. Returns student data

##### 3. **GetList(storage Storage.Storage) http.HandlerFunc**
- **Purpose:** Retrieves all students
- **Method:** GET
- **Route:** `/api/students`
- **Input:** None
- **Output:** JSON array of students
- **Process:**
  1. Calls storage layer to fetch all students
  2. Returns list of students

#### Configuration Functions

**File:** `internal/config/config.go`

##### **MustLoad() *Config**
- **Purpose:** Loads configuration from YAML file
- **Process:**
  1. Reads CONFIG_PATH environment variable or command-line flag
  2. Validates config file exists
  3. Parses YAML using `cleanenv`
  4. Returns Config struct or panics if load fails
- **Error Handling:** Uses `log.Fatal()` to stop execution on missing config

#### Storage Interface Functions

**File:** `internal/storage/storage.go`

```go
type Storage interface {
    CreateStudent(name string, email string, age int) (int64, error)
    GetStudentById(id int64) (types.Student, error)
    GetStudents() ([]types.Student, error)
}
```

**Purpose:** Defines storage contract that SQLite implements.

#### SQLite Implementation Functions

**File:** `internal/storage/sqlite/sqlite.go`

##### 1. **New(cfg *config.Config) (*Sqlite, error)**
- **Purpose:** Initializes SQLite database connection and creates schema
- **Process:**
  1. Creates storage directory if it doesn't exist
  2. Opens SQLite database connection
  3. Verifies connection with Ping()
  4. Creates `students` table if it doesn't exist
  5. Returns Sqlite struct with database connection

##### 2. **CreateStudent(name, email string, age int) (int64, error)**
- **Purpose:** Inserts a new student into database
- **Process:**
  1. Prepares INSERT statement
  2. Executes statement with parameters (name, email, age)
  3. Retrieves and returns the last inserted ID
- **SQL:** `INSERT INTO students(name, email, age) VALUES (?, ?, ?)`

##### 3. **GetStudentById(id int64) (types.Student, error)**
- **Purpose:** Retrieves a specific student by ID
- **Process:**
  1. Prepares SELECT statement with WHERE clause
  2. Queries single row
  3. Scans result into Student struct
  4. Handles `sql.ErrNoRows` with custom error message
- **SQL:** `SELECT * FROM students WHERE id = ? LIMIT 1`

##### 4. **GetStudents() ([]types.Student, error)**
- **Purpose:** Retrieves all students from database
- **Process:**
  1. Prepares SELECT statement
  2. Queries all rows
  3. Iterates through rows and scans into Student structs
  4. Returns slice of students
- **SQL:** `SELECT id, name, email, age FROM students`

#### Response Utility Functions

**File:** `internal/utils/response/response.go`

##### 1. **WriteJson(w http.ResponseWriter, status int, data interface{}) error**
- **Purpose:** Writes JSON response with appropriate headers and status code
- **Process:**
  1. Sets `Content-Type: application/json` header
  2. Writes HTTP status code
  3. Encodes data to JSON and writes to response writer

##### 2. **GeneralError(err error) Response**
- **Purpose:** Creates a standard error response
- **Returns:** Response struct with Status="Error" and Error message

##### 3. **ValidationError(errs validator.ValidationErrors) Response**
- **Purpose:** Formats validation errors into readable message
- **Process:**
  1. Iterates through validation errors
  2. Generates human-readable error messages
  3. Joins all messages with "; " separator
- **Example:** "field name is required field; field email is invalid"

#### Main Function

**File:** `cmd/students-api/main.go`

**Purpose:** Application entry point

**Key Components:**

1. **corsMiddleware(next http.Handler) http.Handler**
   - Adds CORS headers to allow cross-origin requests
   - Handles OPTIONS requests

2. **main() function**
   - Loads configuration
   - Initializes SQLite storage
   - Sets up HTTP router with three routes:
     - POST /api/students
     - GET /api/students/{id}
     - GET /api/students
   - Applies CORS middleware
   - Starts HTTP server on configured address
   - Implements graceful shutdown with signal handling

### Storage Layer

The storage layer uses an interface-based design for flexibility:

**File:** `internal/storage/storage.go`

```go
type Storage interface {
    CreateStudent(name string, email string, age int) (int64, error)
    GetStudentById(id int64) (types.Student, error)
    GetStudents() ([]types.Student, error)
}
```

**Benefits:**
- Easy to swap implementations (SQLite, PostgreSQL, etc.)
- Easier to test with mock implementations
- Loose coupling between HTTP handlers and database

**SQLite Database Schema:**

```sql
CREATE TABLE IF NOT EXISTS students(
    Id    INTEGER PRIMARY KEY AUTOINCREMENT,
    name  TEXT,
    email TEXT,
    age   INTEGER
)
```

### Configuration

**File:** `config/config.yaml`

```yaml
env: "dev"                           # Environment: dev or prod
storage_path: "storage/storage.db"  # Path to SQLite database file
http_server:
  address: "localhost:8082"          # Server listening address
```

**Loading Process:**
1. Check `CONFIG_PATH` environment variable
2. If not set, check command-line `-config` flag
3. Parse YAML file using `cleanenv`
4. Validate all required fields are present

### Response Handling

The API uses a consistent response format for both success and error responses:

**Success Response Example:**
```json
{
    "id": 1
}
```

**Error Response Example:**
```json
{
    "Status": "Error",
    "Error": "field name is required field"
}
```

**HTTP Status Codes Used:**
- `200 OK` - Successful GET request
- `201 Created` - Student successfully created
- `400 Bad Request` - Invalid request (validation error, malformed JSON)
- `500 Internal Server Error` - Database or server error

---

## Packages & Dependencies

### External Packages

1. **github.com/ilyakaznacheev/cleanenv**
   - **Purpose:** Parse and validate YAML configuration files
   - **Usage:** Reads config.yaml and maps to Config struct
   - **Version:** v1.5.0

2. **github.com/mattn/go-sqlite3**
   - **Purpose:** SQLite driver for Go
   - **Usage:** Enables database/sql package to work with SQLite
   - **Version:** v1.14.44
   - **Note:** Requires CGO enabled

3. **github.com/go-playground/validator/v10**
   - **Purpose:** Struct and field validation
   - **Usage:** Validates Student struct fields (required validation)
   - **Version:** v10.30.3
   - **Supports:** Tags like `validate:"required"`

4. **github.com/joho/godotenv**
   - **Purpose:** Load environment variables from .env files
   - **Status:** Included as indirect dependency
   - **Version:** v1.5.1

### Standard Library Packages

1. **net/http** - HTTP server and routing
2. **database/sql** - Database abstraction layer
3. **encoding/json** - JSON encoding/decoding
4. **context** - Context for request/response handling
5. **log/slog** - Structured logging
6. **os** - OS operations (environment variables, signals)
7. **strconv** - String conversions
8. **flag** - Command-line flag parsing
9. **errors** - Error handling
10. **fmt** - String formatting

---

## Git Commit History

This section details each commit in the project, explaining what was added, why, and the overall impact on the codebase.

### Commit 1: `7ebbee3` - Initial Commit
**Date:** Mon Jun 1 18:03:13 2026 +0530

**Changes:**
- Created `.gitignore` with 32 lines (standard Go exclusions)
- Created `cmd/students-api/main.go` with basic skeleton:
  ```go
  package main
  
  import (
      "log"
  )
  
  func main() {
      log.Println("Hello, World!")
  }
  ```
- Created `go.mod` with module declaration:
  ```
  module github.com/kaaaxxx/students-api
  go 1.26.3
  ```

**Impact:** 
- Project initialized with basic Go project structure
- Sets foundation for version control and dependency management
- Total: 42 lines added

---

### Commit 2: `b416e7e` - Add Config File
**Date:** Mon Jun 1 18:09:26 2026 +0530

**Changes:**
- Created `config/config.yaml`:
  ```yaml
  env: "dev"
  storage_path: "storage/storage.db"
  http_server:
    address: "localhost:8082"
  ```
- Updated `.gitignore` to exclude config file

**Impact:** 
- Externalized configuration management
- Allows different configurations for dev/production
- Total: 6 lines added to config.yaml

**Rationale:** Configuration files shouldn't be in version control due to sensitive data and environment-specific values.

---

### Commit 3: `392119d` - Add Config to Gitignore
**Date:** Mon Jun 1 18:13:10 2026 +0530

**Changes:**
- Updated `.gitignore` to include `config/config.yaml`

**Impact:** 
- Ensures configuration files won't accidentally be committed
- Protects sensitive information in version control
- Total: 1 line added

**Rationale:** Configuration with database paths and server addresses shouldn't be in public repositories.

---

### Commit 4: `c052740` - Parse Config File
**Date:** Mon Jun 1 19:16:10 2026 +0530

**Changes:**
- Created `internal/config/config.go` (51 lines):
  - Defined `HTTPServer` struct with YAML tags
  - Defined `Config` struct with nested HTTPServer
  - Implemented `MustLoad()` function to:
    - Check CONFIG_PATH environment variable
    - Support -config command-line flag
    - Parse YAML using `cleanenv` package
    - Validate config file exists

- Updated `go.mod` and `go.sum`:
  - Added dependency: `github.com/ilyakaznacheev/cleanenv v1.5.0`
  - Added dependencies: `BurntSushi/toml`, `joho/godotenv` (indirect)

**Impact:** 
- Application can now load and parse configuration
- Supports both environment variables and CLI flags
- 71 lines added to project
- Foundation for externalized configuration

**Codebase Changes:**
```
Files modified: 3
Insertions: 71
go.mod now has cleanenv dependency
```

---

### Commit 5: `8c0bfa7` - Add Request Validation
**Date:** Tue Jun 2 08:21:44 2026 +0530

**Changes:**

1. **Created `internal/types/types.go`** - Defined data models:
   ```go
   type Student struct {
       Id    int64  `json:"id"`
       Name  string `json:"name" validate:"required"`
       Email string `json:"email" validate:"required"`
       Age   int    `json:"age" validate:"required"`
   }
   ```

2. **Created `internal/utils/response/response.go`** (52 lines):
   - `Response` struct for standardized responses
   - `WriteJson()` function to write JSON with headers
   - `GeneralError()` for error responses
   - `ValidationError()` to format validation errors

3. **Enhanced `cmd/students-api/main.go`** (58 lines added):
   - Imported packages: net/http, log, slog
   - Created basic HTTP routes structure
   - Added corsMiddleware function
   - Set up basic router

4. **Created `internal/http/handlers/student/student.go`** (43 lines):
   - `New()` handler for POST /api/students
   - Implemented JSON decoding
   - Added request validation
   - Returns error responses

5. **Updated `config/config.yaml`**:
   - Added storage_path and http_server configuration

6. **Updated `go.mod`**:
   - Added: `github.com/go-playground/validator/v10`
   - Added multiple indirect dependencies for validation

**Impact:**
- 192 lines added across 8 files
- Project now has HTTP handlers and response structure
- Request validation framework in place
- Students model defined with JSON serialization
- Validation pipeline established

**Codebase Structure After Commit:**
```
HTTP Handlers → Validation → Response Formatting → Response Sent
```

**Key Features Added:**
- CORS middleware for cross-origin requests
- Student data structure with required field validation
- Error response formatting with field-specific messages

---

### Commit 6: `68791e5` - Implement SQLite Storage and Update Student Creation Logic
**Date:** Tue Jun 2 09:52:59 2026 +0530

**Changes:**

1. **Created `internal/storage/storage.go`** - Storage interface:
   ```go
   type Storage interface {
       CreateStudent(name string, email string, age int) (int64, error)
       GetStudentById(id int64) (types.Student, error)
       GetStudents() ([]types.Student, error)
   }
   ```

2. **Created `internal/storage/sqlite/sqlite.go`** (58 lines):
   - `Sqlite` struct wrapping `*sql.DB`
   - `New()` function:
     - Creates storage directory
     - Opens SQLite connection
     - Creates `students` table with schema:
       ```sql
       CREATE TABLE IF NOT EXISTS students(
           Id INTEGER PRIMARY KEY AUTOINCREMENT,
           name TEXT,
           email TEXT,
           age INTEGER
       )
       ```
   - `CreateStudent()` - INSERT operation
   - `GetStudentById()` - SELECT by ID
   - `GetStudents()` - SELECT all (prepared for future use)

3. **Updated `cmd/students-api/main.go`**:
   - Added storage initialization
   - Setup database connection
   - Added logging for storage init
   - Cleanup with defer

4. **Updated handler** - Modified student creation handler

5. **Updated `go.mod` and `go.sum`**:
   - Added: `github.com/mattn/go-sqlite3 v1.14.44`
   - Added: `github.com/ncruces/go-sqlite3` (optional alternative)

6. **Updated `.gitignore`**:
   - Added storage database file exclusion

**Impact:**
- 104 lines added
- Database persistence layer implemented
- Three storage operations available
- Auto-increment primary key for students
- Foundation for data persistence

**Database Schema:**
```sql
students table:
├── Id (INTEGER, PRIMARY KEY, AUTOINCREMENT)
├── name (TEXT)
├── email (TEXT)
└── age (INTEGER)
```

**Codebase Architecture:**
```
HTTP Handlers
    ↓
Validation
    ↓
Storage Interface
    ↓
SQLite Implementation
    ↓
Database
```

---

### Commit 7: `c006ce5` - Add Get Student by ID
**Date:** Fri Jun 5 09:20:57 2026 +0530

**Changes:**

1. **Enhanced `internal/http/handlers/student/student.go`**:
   - Added `GetById()` handler function:
     - Extracts ID from URL path parameter
     - Parses string ID to int64
     - Calls storage.GetStudentById()
     - Returns student JSON or error

2. **Enhanced `internal/storage/sqlite/sqlite.go`**:
   - Implemented `GetStudentById()` method:
     - Prepared SELECT statement with WHERE clause
     - Handles `sql.ErrNoRows` with user-friendly error
     - Returns full Student struct

3. **Updated `cmd/students-api/main.go`**:
   - Added route: `GET /api/students/{id}`
   - Connected GetById handler

4. **Updated `internal/types/types.go`**:
   - Minor adjustments to Student struct

5. **Updated `folder-structure.md`**:
   - Added documentation (50 lines)

**Impact:**
- 131 lines added across 6 files
- API now supports retrieving individual students
- Error handling for non-existent students
- URL path parameter handling demonstrated

**New Endpoint:**
```
GET /api/students/{id}
Response: 200 OK
{
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 20
}
```

**SQL Query Used:**
```sql
SELECT * FROM students WHERE id = ? LIMIT 1
```

---

### Commit 8: `f7c571c` - Add Students List
**Date:** Fri Jun 5 09:49:18 2026 +0530

**Changes:**

1. **Enhanced `internal/http/handlers/student/student.go`**:
   - Added `GetList()` handler function:
     - Calls storage.GetStudents()
     - Returns array of all students
     - Error handling for storage failures

2. **Enhanced `internal/storage/sqlite/sqlite.go`**:
   - Implemented `GetStudents()` method:
     - Prepares SELECT statement
     - Iterates through result rows
     - Scans each row into Student struct
     - Returns slice of students
     - Handles row scanning errors

3. **Updated `cmd/students-api/main.go`**:
   - Added route: `GET /api/students`
   - Connected GetList handler

**Impact:**
- 54 lines added across 4 files
- API now supports listing all students
- Iterator pattern for multiple records
- Foundation for pagination (future enhancement)

**New Endpoint:**
```
GET /api/students
Response: 200 OK
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

**SQL Query Used:**
```sql
SELECT id, name, email, age FROM students
```

---

### Commit 9: `54ccece` - Initial Project Structure and Student API Implementation (Latest)
**Date:** Fri Jun 5 10:58:18 2026 +0530

**Changes:**

1. **Updated `README.md`** (175 lines):
   - Comprehensive project documentation
   - API endpoint descriptions
   - Setup instructions
   - Usage examples

2. **Updated `cmd/students-api/main.go`** (18 insertions, 1 deletion):
   - Minor refinements and clean-up
   - Final version of main function

3. **Created `folder-structure.md`** (2 lines):
   - Project structure overview

4. **Created `frontend/index.html`** (275 lines):
   - Interactive HTML interface
   - Form for creating students
   - Table for displaying students
   - JavaScript for API integration

5. **Created `frontend/plan.md`** (486 lines):
   - Detailed frontend documentation
   - UI/UX planning
   - Feature descriptions
   - Implementation notes

**Impact:**
- 955 lines added
- Project now has comprehensive documentation
- Frontend interface created
- Full project completion documented
- Ready for production use

**Overall Project Statistics:**
```
Total Commits: 9
Total Lines Added: 2,500+
Key Components:
├── Backend API (Go)
├── Database (SQLite)
├── Frontend (HTML/JS)
└── Documentation
```

---

## Project Evolution Summary

### Phase 1: Foundation (Commits 1-3)
- Project initialization
- Configuration setup
- Git management

### Phase 2: API Framework (Commits 4-5)
- Configuration loading
- HTTP handler structure
- Request validation
- Response formatting

### Phase 3: Data Persistence (Commit 6)
- SQLite integration
- Storage interface design
- Database schema creation

### Phase 4: Core API Operations (Commits 7-8)
- Student retrieval endpoints
- List all students
- GET operations

### Phase 5: Polish & Documentation (Commit 9)
- Frontend implementation
- Comprehensive documentation
- Project completion

---

## Running the Application

### 1. Build the Project
```bash
cd student-api
go build -o students-api ./cmd/students-api
```

### 2. Run the Application
```bash
./students-api -config config/config.yaml
```

Or with environment variable:
```bash
export CONFIG_PATH=config/config.yaml
./students-api
```

### 3. Test the API

**Create a Student:**
```bash
curl -X POST http://localhost:8082/api/students \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 20
  }'
```

**Get All Students:**
```bash
curl http://localhost:8082/api/students
```

**Get Student by ID:**
```bash
curl http://localhost:8082/api/students/1
```

---

## Key Architectural Decisions

1. **Interface-Based Storage:** Using a Storage interface allows easy swapping of implementations.

2. **Structured Logging:** Using `log/slog` for structured logging aids debugging and monitoring.

3. **CORS Middleware:** Allows frontend and backend to communicate across origins.

4. **Configuration File:** Externalizes configuration for different environments.

5. **Validation Framework:** Using go-playground/validator for declarative validation.

6. **Prepared Statements:** SQL prepared statements prevent SQL injection vulnerabilities.

7. **Graceful Shutdown:** Proper server shutdown handling with context timeouts.

---

## Future Enhancement Opportunities

- Pagination for large student lists
- Filtering and sorting capabilities
- Update and delete endpoints
- JWT authentication
- Database migration tool
- Unit tests and integration tests
- API versioning
- Rate limiting
- Comprehensive error codes
- Database connection pooling configuration
