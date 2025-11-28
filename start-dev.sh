#!/bin/bash

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

echo "Starting Tonish Development Environment..."

# Start backend
echo "Starting backend on ${HOST_IP:-192.168.5.10}:${BACKEND_PORT:-8080}..."
cd backend && PORT=${BACKEND_PORT:-8080} OLLAMA_URL=${OLLAMA_URL:-http://192.168.5.10:11434} OLLAMA_MODEL=${OLLAMA_MODEL:-qwen2.5-coder:3b} go run main.go &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Start frontend
echo "Starting frontend on ${HOST_IP:-192.168.5.10}:${FRONTEND_PORT:-5173}..."
cd ../frontend && npm run dev -- --host ${HOST_IP:-192.168.5.10} --port ${FRONTEND_PORT:-5173} &
FRONTEND_PID=$!

echo ""
echo "✅ Tonish is running!"
echo "📱 Frontend: http://${HOST_IP:-192.168.5.10}:${FRONTEND_PORT:-5173}"
echo "🔧 Backend API: http://${HOST_IP:-192.168.5.10}:${BACKEND_PORT:-8080}"
echo "🤖 AI (Ollama): ${OLLAMA_URL:-http://192.168.5.10:11434}"
echo ""
echo "Note: Using external Ollama instance at 192.168.5.10:11434"
echo "  To verify Ollama: curl http://192.168.5.10:11434/api/tags"
echo "  To test AI: curl http://192.168.5.10:8080/api/ai/health"
echo ""
echo "Press Ctrl+C to stop both servers"

# Wait for Ctrl+C
trap "kill $BACKEND_PID $FRONTEND_PID" EXIT
wait
