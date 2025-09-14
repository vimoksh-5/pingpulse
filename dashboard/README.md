# PingPulse Dashboard

A modern React TypeScript dashboard for monitoring PingPulse services. This dashboard provides a user-friendly interface for product teams, managers, and other non-technical users to monitor service health without needing to interact with APIs directly.

## Features

### 🎯 **Real-time Monitoring**

- Live service status updates every 30 seconds
- Real-time connection status indicator
- Automatic refresh with manual refresh option

### 📊 **Comprehensive Overview**

- **Statistics Cards**: Total projects, monitors, uptime percentage, and average latency
- **Status Overview**: Quick counts of healthy, warning, critical, and unknown monitors
- **Project Filtering**: Filter monitors by specific projects
- **Search Functionality**: Search monitors by name, URL, or project

### 🚦 **Status Indicators**

- **Healthy**: Green - Service is responding normally (200-299 status codes, low latency)
- **Warning**: Yellow - Service has issues but is functional (high latency or 4xx errors)
- **Critical**: Red - Service is down or severely impacted (5xx errors, timeouts, or very high latency)
- **Unknown**: Gray - Status cannot be determined

### 💡 **Monitor Details**

Each monitor card displays:

- Service name and project
- Current status with color-coded badge
- HTTP status code
- Response latency with color coding
- Uptime percentage
- Last check timestamp
- Service URL

### 🎨 **Modern UI/UX**

- Clean, responsive design built with Tailwind CSS
- Mobile-friendly layout
- Loading states and error handling
- Intuitive navigation and filtering

## Prerequisites

Before running the dashboard, ensure you have:

1. **PingPulse Backend Running**: The dashboard connects to `http://localhost:8080` by default
2. **Node.js**: Version 14 or higher
3. **npm**: Comes with Node.js

## Installation

1. **Navigate to the dashboard directory**:

   ```bash
   cd dashboard
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

## Running the Dashboard

### Development Mode

Start the development server:

```bash
npm start
```

The dashboard will open at `http://localhost:3000` and automatically reload when you make changes.

### Production Build

Build the dashboard for production:

```bash
npm run build
```

This creates an optimized build in the `build` folder ready for deployment.

## Configuration

### Backend URL

By default, the dashboard connects to `http://localhost:8080`. To change this:

1. **Environment Variable**: Set `REACT_APP_API_URL`

   ```bash
   REACT_APP_API_URL=http://your-pingpulse-server:8080 npm start
   ```

2. **Code Change**: Modify the base URL in `src/services/api.ts`

### Auto-refresh Settings

The dashboard auto-refreshes every 30 seconds by default. To customize:

- Edit the `refreshInterval` parameter in `src/components/Dashboard.tsx`
- Set to `false` to disable auto-refresh

## Usage Guide

### For Product Teams

- **Monitor Service Health**: Get a quick overview of all services across projects
- **Track Uptime**: Monitor service reliability and uptime percentages
- **Identify Issues**: Quickly spot critical or warning services that need attention
- **Project Focus**: Filter by specific projects to focus on relevant services

### For Engineering Teams

- **System Overview**: Get a high-level view of system health
- **Performance Monitoring**: Track response times and identify slow services
- **Incident Response**: Quickly identify and prioritize issues during incidents
- **Historical Context**: View recent check results and trends

### Dashboard Sections

1. **Header**

   - Connection status indicator
   - Last updated timestamp
   - Manual refresh button

2. **Statistics Cards**

   - Total projects and monitors
   - Overall system uptime
   - Average response latency

3. **Filters**

   - Project dropdown filter
   - Search box for monitors
   - Status count overview

4. **Monitor Grid**
   - Individual monitor cards
   - Status, latency, and uptime information
   - Responsive grid layout

## API Integration

The dashboard integrates with PingPulse's REST API:

- **GET /api/v1/statistics**: Overall system statistics
- **GET /api/v1/projects**: List of all projects
- **GET /api/v1/monitors**: List of all monitors
- **GET /api/v1/results**: Monitor results with filtering
- **GET /healthz**: Health check endpoint

## Troubleshooting

### Dashboard Won't Load

1. Ensure PingPulse backend is running on `http://localhost:8080`
2. Check browser console for CORS or network errors
3. Verify the backend has CORS enabled (should be automatic with recent updates)

### No Data Showing

1. Confirm monitors are configured in PingPulse
2. Check that monitors have run at least once
3. Verify API endpoints are responding (visit `http://localhost:8080/healthz`)

### Connection Issues

1. Check network connectivity to PingPulse backend
2. Verify firewall settings allow connections to port 8080
3. Ensure backend is accessible from your machine

## Development

### Project Structure

```
src/
├── components/          # React components
│   ├── Dashboard.tsx    # Main dashboard component
│   ├── MonitorCard.tsx  # Individual monitor display
│   ├── StatusBadge.tsx  # Status indicator component
│   ├── StatsCard.tsx    # Statistics card component
│   └── LoadingSpinner.tsx
├── hooks/              # Custom React hooks
│   └── usePingPulse.ts # Main data fetching hook
├── services/           # API integration
│   └── api.ts          # PingPulse API client
├── types/              # TypeScript type definitions
│   └── api.ts          # API response types
├── utils/              # Utility functions
│   └── helpers.ts      # Helper functions
└── App.tsx             # Main app component
```

### Adding Features

1. **New Components**: Add to `src/components/`
2. **API Integration**: Extend `src/services/api.ts`
3. **Type Definitions**: Update `src/types/api.ts`
4. **Styling**: Use Tailwind CSS classes

### Contributing

1. Follow the existing code style and patterns
2. Add TypeScript types for new features
3. Test with the PingPulse backend
4. Ensure responsive design works on mobile



### Environment Variables

For production deployments, set:

- `REACT_APP_API_URL`: PingPulse backend URL
- `REACT_APP_REFRESH_INTERVAL`: Auto-refresh interval in milliseconds

## License

This dashboard is part of the PingPulse project and follows the same license terms.
