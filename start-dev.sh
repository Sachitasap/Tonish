#!/bin/bash

echo "Starting Tonish Development Environment..."

# Start backend
echo "Starting backend on port 3000..."
cd backend && go run main.go &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Start frontend
echo "Starting frontend on port 5173..."
cd ../frontend && npm run dev &
FRONTEND_PID=$!

echo ""
echo "✅ Tonish is running!"
echo "📱 Frontend: http://localhost:5173"
echo "🔧 Backend API: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop both servers"

# Wait for Ctrl+C
trap "kill $BACKEND_PID $FRONTEND_PID" EXIT
wait
