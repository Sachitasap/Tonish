#!/bin/bash

# Tonish AI Setup Script
# This script configures Tonish to use existing Ollama at 192.168.5.10:11434

set -e

echo "╔════════════════════════════════════════════════════════╗"
echo "║        Tonish AI Integration Setup                     ║"
echo "║        Using Existing Ollama Instance                  ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

# Check if Ollama is accessible
echo "🔍 Checking Ollama connection at 192.168.5.10:11434..."
if curl -s http://192.168.5.10:11434/api/tags > /dev/null 2>&1; then
    echo "✓ Ollama is accessible!"
else
    echo "❌ Error: Cannot connect to Ollama at 192.168.5.10:11434"
    echo "Please ensure Ollama is running"
    exit 1
fi

echo ""
echo "📋 Listing available models..."
MODELS=$(curl -s http://192.168.5.10:11434/api/tags | grep -o '"name":"[^"]*"' | cut -d'"' -f4)

if echo "$MODELS" | grep -q "qwen2.5-coder"; then
    INSTALLED_MODEL=$(echo "$MODELS" | grep "qwen2.5-coder" | head -1)
    echo "✓ qwen2.5-coder model is already installed: $INSTALLED_MODEL"
    echo "  Using this model for Tonish AI"
else
    echo "⚠️  qwen2.5-coder model not found"
    echo "📥 Would you like to pull qwen2.5-coder:3b now? (y/n)"
    read -r RESPONSE
    if [[ "$RESPONSE" =~ ^[Yy]$ ]]; then
        echo "Pulling qwen2.5-coder:3b (this may take a few minutes)..."
        curl -X POST http://192.168.5.10:11434/api/pull \
            -d '{"name": "qwen2.5-coder:3b"}' \
            -H "Content-Type: application/json"
        echo ""
        echo "✓ Model pulled successfully!"
    fi
fi

echo ""
echo "🧪 Testing AI service..."

# Test the model
TEST_RESPONSE=$(curl -s http://192.168.5.10:11434/api/generate \
    -d '{
        "model": "qwen2.5-coder:3b",
        "prompt": "Reply with just the word ready",
        "stream": false
    }' | grep -o '"response":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -n "$TEST_RESPONSE" ]; then
    echo "✓ AI model is responding!"
    echo "  Response preview: ${TEST_RESPONSE:0:50}..."
else
    echo "⚠️  Warning: AI model may not be fully ready"
fi

echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║              Setup Complete! 🎉                        ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo "🤖 Ollama: http://192.168.5.10:11434"
echo "📖 Model: qwen2.5-coder:3b"
echo ""
echo "Configuration:"
echo "  OLLAMA_URL=http://192.168.5.10:11434"
echo "  OLLAMA_MODEL=qwen2.5-coder:3b"
echo ""
echo "Next steps:"
echo "  1. Start the full stack:"
echo "     $ docker-compose up -d"
echo ""
echo "  2. OR start development mode:"
echo "     $ ./start-dev.sh"
echo ""
echo "  3. Test AI health:"
echo "     $ curl http://192.168.5.10:8080/api/ai/health"
echo ""
echo "📚 Read AI_INTEGRATION_GUIDE.md for detailed documentation"
echo ""
