# Batch Configuration Update Endpoint

## Overview

The batch configuration update endpoint allows updating multiple configuration values in a single atomic request.

## Endpoint

```
PATCH /api/v1/instances/{name}/config
```

## Request Format

```json
{
  "updates": [
    {"path": "string", "value": "any"},
    {"path": "string", "value": "any"}
  ]
}
```

## Response Format

Success (200 OK):
```json
{
  "message": "Configuration updated successfully",
  "updated": 3
}
```

Error (400 Bad Request / 404 Not Found / 500 Internal Server Error):
```json
{
  "error": "error message"
}
```

## Implementation Details

### File Location
- Handler: `/daemon/internal/api/v1/handlers_config.go`
- Route registration: `/daemon/internal/api/v1/handlers.go` (line 62)

### Key Features

1. **Validation**:
   - Validates instance exists before processing
   - Validates all paths are non-empty before applying updates
   - Returns clear error messages for validation failures

2. **Atomicity**:
   - Uses existing config manager's file locking mechanism
   - Each update is applied sequentially within a lock
   - If any update fails, the operation stops and returns an error

3. **Type Handling**:
   - Accepts any JSON value type (string, number, boolean, object, array)
   - Converts values to strings for storage using `fmt.Sprintf("%v", value)`

## Usage Examples

### Example 1: Update Basic Configuration Values

```bash
curl -X PATCH http://localhost:8080/api/v1/instances/my-cloud/config \
  -H "Content-Type: application/json" \
  -d '{
    "updates": [
      {"path": "baseDomain", "value": "example.com"},
      {"path": "domain", "value": "wild.example.com"},
      {"path": "internalDomain", "value": "int.wild.example.com"}
    ]
  }'
```

Response:
```json
{
  "message": "Configuration updated successfully",
  "updated": 3
}
```

### Example 2: Update Nested Configuration Values

```bash
curl -X PATCH http://localhost:8080/api/v1/instances/my-cloud/config \
  -H "Content-Type: application/json" \
  -d '{
    "updates": [
      {"path": "cluster.name", "value": "prod-cluster"},
      {"path": "cluster.loadBalancerIp", "value": "192.168.1.100"},
      {"path": "cluster.ipAddressPool", "value": "192.168.1.100-192.168.1.200"}
    ]
  }'
```

### Example 3: Update Array Values

```bash
curl -X PATCH http://localhost:8080/api/v1/instances/my-cloud/config \
  -H "Content-Type: application/json" \
  -d '{
    "updates": [
      {"path": "cluster.nodes.activeNodes[0]", "value": "node-1"},
      {"path": "cluster.nodes.activeNodes[1]", "value": "node-2"}
    ]
  }'
```

### Example 4: Error Handling - Invalid Instance

```bash
curl -X PATCH http://localhost:8080/api/v1/instances/nonexistent/config \
  -H "Content-Type: application/json" \
  -d '{
    "updates": [
      {"path": "baseDomain", "value": "example.com"}
    ]
  }'
```

Response (404):
```json
{
  "error": "Instance not found: instance nonexistent does not exist"
}
```

### Example 5: Error Handling - Empty Updates

```bash
curl -X PATCH http://localhost:8080/api/v1/instances/my-cloud/config \
  -H "Content-Type: application/json" \
  -d '{
    "updates": []
  }'
```

Response (400):
```json
{
  "error": "updates array is required and cannot be empty"
}
```

### Example 6: Error Handling - Missing Path

```bash
curl -X PATCH http://localhost:8080/api/v1/instances/my-cloud/config \
  -H "Content-Type: application/json" \
  -d '{
    "updates": [
      {"path": "", "value": "example.com"}
    ]
  }'
```

Response (400):
```json
{
  "error": "update[0]: path is required"
}
```

## Configuration Path Syntax

The `path` field uses YAML path syntax as implemented by the `yq` tool:

- Simple fields: `baseDomain`
- Nested fields: `cluster.name`
- Array elements: `cluster.nodes.activeNodes[0]`
- Array append: `cluster.nodes.activeNodes[+]`

Refer to the yq documentation for advanced path syntax.

## Bonus: Bug Fix

While implementing this endpoint, a bug was discovered and fixed in the existing `UpdateConfig` (PUT) handler.

**Issue**: The handler was passing the instance name instead of the config file path to `SetConfigValue`.

**Fixed in**: `/daemon/internal/api/v1/handlers.go` (line 269)

**Before**:
```go
if err := api.config.SetConfigValue(name, key, valueStr); err != nil {
```

**After**:
```go
configPath := api.instance.GetInstanceConfigPath(name)
// ...
if err := api.config.SetConfigValue(configPath, key, valueStr); err != nil {
```

## Testing

To test the endpoint:

1. Start the daemon:
   ```bash
   cd daemon
   go run main.go
   ```

2. Create a test instance:
   ```bash
   curl -X POST http://localhost:8080/api/v1/instances \
     -H "Content-Type: application/json" \
     -d '{"name": "test-instance"}'
   ```

3. Update configuration using batch endpoint:
   ```bash
   curl -X PATCH http://localhost:8080/api/v1/instances/test-instance/config \
     -H "Content-Type: application/json" \
     -d '{
       "updates": [
         {"path": "baseDomain", "value": "example.com"},
         {"path": "domain", "value": "wild.example.com"}
       ]
     }'
   ```

4. Verify the changes:
   ```bash
   curl http://localhost:8080/api/v1/instances/test-instance/config
   ```
