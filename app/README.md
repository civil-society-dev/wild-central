# Wild-Central Web UI

A React-based web interface for managing wild-central network appliance configuration and operations.

## Overview

This application provides a modern web UI for configuring and monitoring a wild-central network appliance, which serves as DNS, DHCP, and PXE boot infrastructure for local cloud deployments. The interface handles:

- **System Status**: Real-time monitoring of network services and system health
- **Configuration Management**: YAML-based configuration with validation and form-based editing
- **DNS/DHCP Management**: dnsmasq configuration and service control
- **PXE Boot Assets**: Talos Linux image management and boot configuration
- **Network Appliance Setup**: Guided setup for local cloud infrastructure

## Tech Stack

- **React 19** with TypeScript
- **Vite** for build tooling and development server
- **TanStack Query** for server state management
- **React Hook Form + Zod** for form validation
- **Tailwind CSS** for styling with shadcn/ui component library
- **Vitest + React Testing Library** for testing

## Development

### Prerequisites

- Node.js 18+ 
- npm
- Go backend server running on port 5055

### Getting Started

```bash
# Install dependencies
npm install

# Start development server (runs on port 5050)
npm run dev

# Run tests
npm test

# Build for production
npm run build
```

### Available Scripts

- `npm run dev` - Start development server on port 5050
- `npm run build` - Build for production to `build/` directory  
- `npm test` - Run test suite with Vitest
- `npm run test:ui` - Run tests with UI
- `npm run test:coverage` - Run tests with coverage report
- `npm run type-check` - TypeScript type checking

### Development Workflow

The app expects the Go backend to be running on port 5055. Use VSCode launch configurations for full-stack development:

1. **"Go Daemon"** - Starts the Go backend server
2. **"React App"** - Starts the React development server  
3. **"Full Stack"** - Compound configuration that starts both

### API Integration

The React app communicates with the Go backend via REST API:

- `GET /api/status` - System status and health
- `GET /api/config` - Configuration state
- `POST /api/config` - Create/update configuration
- `POST /api/dnsmasq/restart` - Restart dnsmasq service
- And more...

## Architecture

### Key Components

- **App.tsx** - Main application with routing and state management
- **ConfigurationForm** - Form-based YAML configuration editing
- **StatusSection** - Real-time system monitoring
- **DnsmasqSection** - DNS/DHCP service management
- **PxeAssetsSection** - Boot asset management
- **ErrorBoundary** - Application error handling

### State Management

- **TanStack Query** for server state caching and synchronization
- **React Hook Form** for form state and validation
- **Theme Context** for dark/light mode

### Testing

Comprehensive test suite covering:
- React Query hooks
- Form validation schemas
- Component rendering and interactions
- Error boundary behavior

## Configuration

The app manages YAML configuration with schema validation for:

- **Server**: Host and port settings
- **Cloud**: Domain, DHCP range, DNS, and network settings  
- **Cluster**: Endpoint IP and Talos version configuration

## Deployment

Built as a static React app that can be served by any web server. In production, it's served by the Go backend's static file handler.
