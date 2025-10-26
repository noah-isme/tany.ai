#!/bin/bash

# Test login endpoint
API_URL="${1:-http://localhost:8080}"

echo "Testing login endpoint at: $API_URL/api/auth/login"
echo ""

# Test with correct credentials
echo "Testing with admin@example.com / Admin#12345..."
response=$(curl -s -X POST "$API_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"Admin#12345"}')

echo "Response:"
echo "$response" | jq '.' 2>/dev/null || echo "$response"
echo ""

# Check if login successful
if echo "$response" | grep -q "accessToken"; then
  echo "✅ Login successful!"
else
  echo "❌ Login failed!"
fi
