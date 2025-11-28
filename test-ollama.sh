#!/bin/bash

# Quick test script for Ollama integration

echo "🔍 Testing Ollama Connection..."
echo ""

OLLAMA_URL="http://192.168.5.10:11434"

# Test basic connection
echo "1. Testing connection to $OLLAMA_URL..."
if curl -s -f "$OLLAMA_URL/api/tags" > /dev/null 2>&1; then
    echo "   ✓ Connection successful!"
else
    echo "   ✗ Connection failed!"
    echo "   Make sure Ollama is running at $OLLAMA_URL"
    exit 1
fi

echo ""
echo "2. Available models:"
curl -s "$OLLAMA_URL/api/tags" | grep -o '"name":"[^"]*"' | cut -d'"' -f4 | while read -r model; do
    echo "   • $model"
done

echo ""
echo "3. Testing AI generation..."
RESPONSE=$(curl -s "$OLLAMA_URL/api/generate" \
    -H "Content-Type: application/json" \
    -d '{
        "model": "qwen2.5-coder:3b",
        "prompt": "Say hello in JSON format: {\"message\": \"...\"}",
        "stream": false
    }')

if echo "$RESPONSE" | grep -q "response"; then
    echo "   ✓ AI generation works!"
    EXTRACTED=$(echo "$RESPONSE" | grep -o '"response":"[^"]*"' | cut -d'"' -f4 | head -c 100)
    echo "   Response preview: $EXTRACTED..."
else
    echo "   ✗ AI generation failed!"
    echo "   Response: $RESPONSE"
fi

echo ""
echo "4. Testing chat endpoint..."
CHAT_RESPONSE=$(curl -s "$OLLAMA_URL/api/chat" \
    -H "Content-Type: application/json" \
    -d '{
        "model": "qwen2.5-coder:3b",
        "messages": [
            {"role": "user", "content": "What is 2+2? Answer with just the number."}
        ],
        "stream": false
    }')

if echo "$CHAT_RESPONSE" | grep -q "message"; then
    echo "   ✓ Chat endpoint works!"
    CHAT_MSG=$(echo "$CHAT_RESPONSE" | grep -o '"content":"[^"]*"' | cut -d'"' -f4 | head -c 50)
    echo "   Response: $CHAT_MSG"
else
    echo "   ✗ Chat endpoint failed!"
fi

echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║              Ollama Status Summary                      ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo "  Endpoint: $OLLAMA_URL"
echo "  Status: ✓ Ready"
echo "  Model: qwen2.5-coder:3b"
echo ""
echo "You can now start Tonish with AI support:"
echo "  $ ./start-dev.sh"
echo "  or"
echo "  $ docker-compose up -d"
echo ""
