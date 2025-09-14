#!/bin/bash

# PingPulse Dashboard Startup Script
# This script starts both the PingPulse backend and the React dashboard

set -e

echo "🚀 Starting PingPulse with Dashboard..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go to run PingPulse backend."
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js to run the dashboard."
    exit 1
fi

# Build PingPulse backend if needed
if [ ! -f "./pingpulse" ]; then
    echo "📦 Building PingPulse backend..."
    go build -o pingpulse ./cmd/pingpulse
fi

# Check if config exists
if [ ! -f "./config.yaml" ]; then
    echo "❌ config.yaml not found. Please create a configuration file."
    echo "   You can copy config.example.yaml to config.yaml and modify it."
    exit 1
fi

# Install dashboard dependencies if needed
if [ ! -d "./dashboard/node_modules" ]; then
    echo "📦 Installing dashboard dependencies..."
    cd dashboard
    npm install
    cd ..
fi

# Function to cleanup background processes
cleanup() {
    echo ""
    echo "🛑 Shutting down services..."
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    if [ ! -z "$DASHBOARD_PID" ]; then
        kill $DASHBOARD_PID 2>/dev/null || true
    fi
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Start PingPulse backend
echo "🔧 Starting PingPulse backend on http://localhost:8080..."
./pingpulse &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Check if backend is running
if ! curl -s http://localhost:8080/healthz > /dev/null; then
    echo "❌ Backend failed to start or is not responding on port 8080"
    cleanup
fi

echo "✅ Backend is running and healthy"

# Start dashboard
echo "🎨 Starting React dashboard on http://localhost:3000..."
cd dashboard
npm start &
DASHBOARD_PID=$!
cd ..

echo ""
echo "🎉 PingPulse Dashboard is starting up!"
echo ""
echo "📊 Backend API: http://localhost:8080"
echo "🌐 Dashboard:   http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop both services"
echo ""

# Wait for processes
wait $BACKEND_PID $DASHBOARD_PID
