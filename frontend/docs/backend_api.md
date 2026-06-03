# HRPML API Documentation

## Base URL

```
http://localhost:34501
```

---

## Endpoints

### 1. Health Check

**GET** `/`

Returns a welcome message.

**Response:**
```json
{
  "message": "Hello, Gin!"
}
```

---

### 2. Ping

**GET** `/ping`

Returns a simple pong response for health checks.

**Response:**
```json
{
  "message": "pong"
}
```

---

### 3. Hello

**GET** `/hello/:name`

Greets the user with the provided name parameter.

**Parameters:**
- `name` (path) - The name to greet

**Response:**
```json
{
  "message": "Hello <name>"
}
```

---

### 4. Echo

**POST** `/echo`

Echoes back the JSON request body.

**Request Body:**
```json
{
  "key": "value",
  "data": 123
}
```

**Response:**
```json
{
  "key": "value",
  "data": 123
}
```

---

### 5. Shell Execution

**POST** `/shell`

Executes a shell command with two execution modes.

**Request Body:**
```json
{
  "shell": "echo hello",
  "type": "simple" | "spawn"
}
```

**Parameters:**
- `shell` (string, required) - The shell command to execute
- `type` (string, optional) - Execution mode: `simple` (default) or `spawn`

#### Type: `simple`

Executes the command synchronously and returns the output after completion.

**Example Request:**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"shell":"echo hello && sleep 1 && echo world", "type":"simple"}' \
  http://localhost:8080/shell
```

**Example Response:**
```json
{
  "success": true,
  "output": "hello\nworld\n"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "error message",
  "output": "partial output if any"
}
```

---

#### Type: `spawn`

Executes the command and streams logs in real-time using chunked transfer encoding. Each log entry is a separate JSON line.

**Example Request:**
```bash
curl -N -X POST -H "Content-Type: application/json" \
  -d '{"shell":"echo line1 && sleep 1 && echo line2", "type":"spawn"}' \
  http://localhost:8080/shell
```

**Example Response (streaming):**
```json
{"timestamp":"2026-05-31 19:59:37.129","type":"stdout","content":"line1"}
{"timestamp":"2026-05-31 19:59:38.129","type":"stdout","content":"line2"}
{"timestamp":"2026-05-31 19:59:39.131","type":"system","content":"Process completed"}
```

**Log Entry Format:**
| Field | Type | Description |
|-------|------|-------------|
| `timestamp` | string | ISO 8601 format time when the log was captured |
| `type` | string | Log source: `stdout`, `stderr`, or `system` |
| `content` | string | The actual log/output content |

---

### 6. File Operations

**POST** `/file`

Performs file operations including create, delete, append, overwrite, and JSON key/value manipulation.

**Request Body:**
```json
{
  "operation": "create" | "delete" | "append" | "overwrite" | "add_json_key" | "modify_json_value" | "delete_json_key",
  "path": "/path/to/your/file",
  "content": "file content (for create, append, overwrite)",
  "key": "json key (for json operations)",
  "value": "json value (for add/modify json key)"
}
```

**Parameters:**
- `operation` (string, required) - The operation to perform
- `path` (string, required) - Path to the file
- `content` (string, optional) - Content for create/append/overwrite operations
- `key` (string, optional) - Key for JSON operations
- `value` (any, optional) - Value for JSON add/modify operations

#### Operations:

**1. Create File**
```json
{
  "operation": "create",
  "path": "./test.txt",
  "content": "Hello, World!"
}
```

**2. Delete File**
```json
{
  "operation": "delete",
  "path": "./test.txt"
}
```

**3. Append to File**
```json
{
  "operation": "append",
  "path": "./test.txt",
  "content": "\nMore content"
}
```

**4. Overwrite File**
```json
{
  "operation": "overwrite",
  "path": "./test.txt",
  "content": "New content"
}
```

**5. Add/Modify JSON Key**
```json
{
  "operation": "add_json_key",
  "path": "./config.json",
  "key": "version",
  "value": "1.0.0"
}
```

**6. Delete JSON Key**
```json
{
  "operation": "delete_json_key",
  "path": "./config.json",
  "key": "version"
}
```

**Success Response:**
```json
{
  "success": true,
  "message": "operation description"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "error message"
}
```

---

### 7. System Information

**GET** `/sysinfo`

Returns system architecture and operating system information.

**Response:**
```json
{
  "arch": "amd64",
  "os": "linux"
}
```

**Fields:**
| Field | Type | Description |
|-------|------|-------------|
| `arch` | string | CPU architecture (e.g., `amd64`, `arm64`) |
| `os` | string | Operating system (`linux`, `darwin`, `windows`) |

**Example Request:**
```bash
curl http://localhost:8080/sysinfo
```

**Example Response:**
```json
{
  "arch": "amd64",
  "os": "linux"
}
```

---

## Error Handling

All endpoints return appropriate HTTP status codes:

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request (invalid JSON, missing required fields) |
| 500 | Internal Server Error |

**Error Response Format:**
```json
{
  "error": "error description"
}
```
