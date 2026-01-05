# File Operations API

## Overview

mcp-language-server bridges LSP file operations to MCP clients via server-initiated notifications. This enables MCP clients to stay synchronized with workspace changes.

## Supported Operations

- **Create**: New files added to workspace
- **Rename**: Files moved or renamed (detected via delete+create within 100ms)
- **Delete**: Files removed from workspace

## Notification Method

**Method:** `notifications/resources/updated`

**Direction:** Server → Client (all connected clients)

**Payload:**

| Field | Type     | Description                              |
|-------|----------|------------------------------------------|
| type  | string   | Operation type: "created", "renamed", "deleted" |
| uris  | string[] | File URIs affected (file:// scheme)     |

## Rename Detection

The server uses a 100ms debounce window to detect rename operations:

1. LSP sends `didDeleteFiles` for old path
2. Within 100ms, LSP sends `didCreateFiles` for new path
3. Server detects pattern and emits single "renamed" notification with both URIs

## Example Notification Payloads

### File Created

```json
{
  "type": "created",
  "uris": ["file:///workspace/main.go"]
}
```

### File Renamed

For rename operations, the URIs array contains pairs of [oldPath, newPath]:

```json
{
  "type": "renamed",
  "uris": [
    "file:///workspace/old.go",
    "file:///workspace/new.go"
  ]
}
```

### File Deleted

```json
{
  "type": "deleted",
  "uris": ["file:///workspace/main.go"]
}
```

### Multiple Files

Batch operations include all affected files:

```json
{
  "type": "created",
  "uris": [
    "file:///workspace/file1.go",
    "file:///workspace/file2.go",
    "file:///workspace/file3.go"
  ]
}
```

## Event Flows

### File Created

```
User creates "main.go" in editor
  → LSP sends workspace/didCreateFiles
  → Server receives event
  → After 100ms debounce: notifications/resources/updated
  → {type: "created", uris: ["file:///workspace/main.go"]}
```

### File Renamed

```
User renames "old.go" → "new.go" in editor
  → LSP sends workspace/didDeleteFiles (old.go)
  → Server stores pending delete
  → LSP sends workspace/didCreateFiles (new.go) [within 100ms]
  → Server detects rename pattern
  → notifications/resources/updated
  → {type: "renamed", uris: ["file:///workspace/old.go", "file:///workspace/new.go"]}
```

### File Deleted

```
User deletes "main.go" in editor
  → LSP sends workspace/didDeleteFiles
  → Server waits 100ms (no matching create)
  → notifications/resources/updated
  → {type: "deleted", uris: ["file:///workspace/main.go"]}
```

## Client Implementation

Clients should:

1. Subscribe to `notifications/resources/updated`
2. Invalidate cached resources for affected URIs
3. Refresh UI/state as needed

Notifications are best-effort; clients should handle missed notifications gracefully.

### Example Client (JavaScript)

```javascript
mcpClient.on('notifications/resources/updated', (params) => {
  console.log(`Files ${params.type}:`, params.uris);

  // Invalidate cache for affected resources
  params.uris.forEach(uri => {
    resourceCache.invalidate(uri);
  });

  // Handle renames specially (URIs come in pairs)
  if (params.type === 'renamed') {
    for (let i = 0; i < params.uris.length; i += 2) {
      const oldUri = params.uris[i];
      const newUri = params.uris[i + 1];
      ui.updateFileTree(oldUri, newUri);
    }
  } else {
    ui.refreshFileTree();
  }
});
```

## Notes

- Notifications are sent to all connected MCP clients
- Delivery is best-effort (no retries or guarantees)
- File URIs always use the `file://` scheme
- All paths are absolute workspace paths
- The 100ms debounce window handles rapid file system operations
