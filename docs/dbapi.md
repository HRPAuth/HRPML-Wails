# Database API Documentation

Base URL: `http://localhost:34501/api/v1/db`

## Table of Contents

- [Overview](#overview)
- [Database Statistics](#database-statistics)
- [Users API](#users-api)
- [Settings API](#settings-api)
- [Logs API](#logs-api)
- [Tasks API](tasks-api)
- [Files API](#files-api)
- [Error Responses](#error-responses)

## Overview

The Database API provides CRUD operations for the following entities:

| Entity | Description |
|--------|-------------|
| Users | User management with roles |
| Settings | Key-value configuration storage |
| Logs | Application logging |
| Tasks | Task/job tracking with progress |
| Files | File metadata tracking |

Database file location: `~/.hrpml/data.db`

---

## Database Statistics

### Get Database Statistics

```
GET /api/v1/db/stats
```

Returns record counts for all tables.

**Response:**

```json
{
  "users": 5,
  "settings": 12,
  "logs": 150,
  "tasks": 8,
  "files": 23
}
```

### Cleanup Database

```
POST /api/v1/db/cleanup
```

Runs VACUUM on the database to reclaim space.

**Response:**

```json
{
  "success": true,
  "message": "database cleaned up"
}
```

---

## Users API

### List All Users

```
GET /api/v1/db/users
```

**Response:**

```json
[
  {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

### Create User

```
POST /api/v1/db/users
```

**Request Body:**

```json
{
  "username": "john",
  "email": "john@example.com",
  "password_hash": "hashed_password",
  "role": "user"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| username | string | Yes | Unique username |
| email | string | No | User email |
| password_hash | string | No | Hashed password |
| role | string | No | User role (`user` or `admin`, default: `user`) |

**Response:** `201 Created`

```json
{
  "id": 2,
  "username": "john",
  "email": "john@example.com",
  "role": "user",
  "created_at": "2024-01-15T11:00:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

### Get User by ID

```
GET /api/v1/db/users/:id
```

**Response:**

```json
{
  "id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "role": "admin",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Update User

```
PUT /api/v1/db/users/:id
```

**Request Body:**

```json
{
  "username": "john_updated",
  "email": "john.new@example.com",
  "role": "admin"
}
```

**Response:** Updated user object

### Delete User

```
DELETE /api/v1/db/users/:id
```

**Response:**

```json
{
  "success": true,
  "message": "user deleted"
}
```

---

## Settings API

### List All Settings

```
GET /api/v1/db/settings
```

**Response:**

```json
[
  {
    "id": 1,
    "key": "app_name",
    "value": "HRPML",
    "type": "string",
    "description": "Application name",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

### Set Setting (Create or Update)

```
POST /api/v1/db/settings
```

**Request Body:**

```json
{
  "key": "theme",
  "value": "dark",
  "type": "string",
  "description": "UI theme preference"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| key | string | Yes | Unique setting key |
| value | string | No | Setting value |
| type | string | No | Value type (`string`, `int`, `bool`, `json`, default: `string`) |
| description | string | No | Setting description |

**Response:**

```json
{
  "id": 2,
  "key": "theme",
  "value": "dark",
  "type": "string",
  "description": "UI theme preference",
  "created_at": "2024-01-15T11:00:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

### Get Setting by Key

```
GET /api/v1/db/settings/:key
```

**Response:** Setting object

### Delete Setting

```
DELETE /api/v1/db/settings/:key
```

**Response:**

```json
{
  "success": true,
  "message": "setting deleted"
}
```

---

## Logs API

### List Logs

```
GET /api/v1/db/logs
```

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| level | string | - | Filter by level (`debug`, `info`, `warn`, `error`) |
| limit | int | 100 | Max records to return |
| offset | int | 0 | Pagination offset |

**Example:**

```
GET /api/v1/db/logs?level=error&limit=50&offset=0
```

**Response:**

```json
[
  {
    "id": 1,
    "level": "error",
    "message": "Connection failed",
    "source": "network",
    "metadata": "{\"retry_count\": 3}",
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

### Create Log Entry

```
POST /api/v1/db/logs
```

**Request Body:**

```json
{
  "level": "info",
  "message": "User logged in",
  "source": "auth",
  "metadata": "{\"user_id\": 1}"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| level | string | No | Log level (`debug`, `info`, `warn`, `error`, default: `info`) |
| message | string | Yes | Log message |
| source | string | No | Source module/component |
| metadata | string | No | JSON metadata |

**Response:** `201 Created`

### Delete Old Logs

```
DELETE /api/v1/db/logs/old?days=30
```

Deletes logs older than specified days.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| days | int | 30 | Delete logs older than N days |

**Response:**

```json
{
  "success": true,
  "message": "old logs deleted"
}
```

---

## Tasks API

### List Tasks

```
GET /api/v1/db/tasks
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| status | string | Filter by status (`pending`, `running`, `completed`, `failed`, `cancelled`) |

**Response:**

```json
[
  {
    "id": 1,
    "name": "Data Import",
    "description": "Import user data from CSV",
    "status": "completed",
    "priority": 1,
    "progress": 100,
    "result": "Imported 1500 records",
    "error": "",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "completed_at": "2024-01-15T10:30:00Z"
  }
]
```

### Create Task

```
POST /api/v1/db/tasks
```

**Request Body:**

```json
{
  "name": "Backup Database",
  "description": "Create daily backup",
  "status": "pending",
  "priority": 2
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Task name |
| description | string | No | Task description |
| status | string | No | Task status (default: `pending`) |
| priority | int | No | Priority (higher = more important) |

**Response:** `201 Created`

### Get Task by ID

```
GET /api/v1/db/tasks/:id
```

### Update Task

```
PUT /api/v1/db/tasks/:id
```

**Request Body:** Full task object

### Update Task Status

```
PATCH /api/v1/db/tasks/:id/status
```

**Request Body:**

```json
{
  "status": "running"
}
```

Valid statuses: `pending`, `running`, `completed`, `failed`, `cancelled`

### Update Task Progress

```
PATCH /api/v1/db/tasks/:id/progress
```

**Request Body:**

```json
{
  "progress": 75
}
```

Progress value: 0-100

### Delete Task

```
DELETE /api/v1/db/tasks/:id
```

---

## Files API

### List File Records

```
GET /api/v1/db/files
```

**Response:**

```json
[
  {
    "id": 1,
    "name": "document.pdf",
    "path": "/home/user/documents/document.pdf",
    "size": 1024000,
    "mime_type": "application/pdf",
    "checksum": "sha256:abc123...",
    "metadata": "{\"pages\": 10}",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
]
```

### Create File Record

```
POST /api/v1/db/files
```

**Request Body:**

```json
{
  "name": "report.xlsx",
  "path": "/home/user/reports/report.xlsx",
  "size": 50000,
  "mime_type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  "checksum": "sha256:def456...",
  "metadata": "{\"sheets\": 3}"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | File name |
| path | string | Yes | File path |
| size | int64 | No | File size in bytes |
| mime_type | string | No | MIME type |
| checksum | string | No | File checksum/hash |
| metadata | string | No | JSON metadata |

**Response:** `201 Created`

### Get File Record by ID

```
GET /api/v1/db/files/:id
```

### Update File Record

```
PUT /api/v1/db/files/:id
```

### Delete File Record

```
DELETE /api/v1/db/files/:id
```

---

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "error message description"
}
```

**HTTP Status Codes:**

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (invalid input) |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## Usage Examples

### cURL Examples

```bash
# Get database stats
curl http://localhost:34501/api/v1/db/stats

# Create a user
curl -X POST http://localhost:34501/api/v1/db/users \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "email": "test@example.com"}'

# Get all settings
curl http://localhost:34501/api/v1/db/settings

# Set a setting
curl -X POST http://localhost:34501/api/v1/db/settings \
  -H "Content-Type: application/json" \
  -d '{"key": "debug_mode", "value": "true", "type": "bool"}'

# Create a task
curl -X POST http://localhost:34501/api/v1/db/tasks \
  -H "Content-Type: application/json" \
  -d '{"name": "Process Data", "priority": 5}'

# Update task progress
curl -X PATCH http://localhost:34501/api/v1/db/tasks/1/progress \
  -H "Content-Type: application/json" \
  -d '{"progress": 50}'

# Get error logs
curl "http://localhost:34501/api/v1/db/logs?level=error&limit=10"
```

### JavaScript/TypeScript Examples

```typescript
// Using fetch
const API_BASE = 'http://localhost:34501/api/v1/db';

// Get stats
const stats = await fetch(`${API_BASE}/stats`).then(r => r.json());

// Create user
const newUser = await fetch(`${API_BASE}/users`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'john',
    email: 'john@example.com',
    role: 'user'
  })
}).then(r => r.json());

// Get settings
const settings = await fetch(`${API_BASE}/settings`).then(r => r.json());

// Create task and update progress
const task = await fetch(`${API_BASE}/tasks`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    name: 'Import Data',
    description: 'Import from CSV',
    priority: 3
  })
}).then(r => r.json());

// Update progress
await fetch(`${API_BASE}/tasks/${task.id}/progress`, {
  method: 'PATCH',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ progress: 50 })
});
```

---

## Database Schema

```sql
-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    password_hash TEXT,
    role TEXT DEFAULT 'user',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Settings table
CREATE TABLE settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT UNIQUE NOT NULL,
    value TEXT,
    type TEXT DEFAULT 'string',
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Logs table
CREATE TABLE logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level TEXT DEFAULT 'info',
    message TEXT NOT NULL,
    source TEXT,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tasks table
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'pending',
    priority INTEGER DEFAULT 0,
    progress INTEGER DEFAULT 0,
    result TEXT,
    error TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME
);

-- Files table
CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    size INTEGER DEFAULT 0,
    mime_type TEXT,
    checksum TEXT,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```
