# Wild UI

React-based user interface for Wild Cloud management.

## Development

```bash
# Install dependencies
make ui-install

# Start development server (http://localhost:3000)
make ui-dev

# Build for production
make ui-build
```

## Features

- Real-time status monitoring from Wild daemon API
- Auto-refresh every 30 seconds
- Responsive design with glassmorphism UI
- Error handling and offline states

## API Integration

The UI fetches status from `http://localhost:5065/api/status` and expects:

```json
{
  "status": "running",
  "timestamp": "2023-12-06T10:30:00Z",
  "version": "0.1.0",
  "uptime": "running"
}
```

## Testing with Wild Daemon

1. Start the Wild daemon: `wild --ui`
2. In another terminal, start UI dev server: `make ui-dev`
3. Open http://localhost:3000 to see the UI
4. The UI will fetch status from the daemon running on port 5065