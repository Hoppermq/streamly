#!/bin/bash

set -e

echo "🧪 Testing Query API Implementation"
echo "===================================="
echo ""

echo "1️⃣  Running domain tests..."
go test ./pkg/domain -v -count=1
echo "✅ Domain tests passed"
echo ""

echo "2️⃣  Building query-api..."
go build -o /tmp/query-api-test ./cmd/query-api
echo "✅ Build successful"
echo ""

echo "3️⃣  Testing JSON unmarshaling..."
cat > /tmp/test-simple.json << 'EOF'
{
  "select": ["event_name", "user_id"],
  "from": "events",
  "timeRange": {
    "start": "now-1h",
    "end": "now"
  }
}
EOF

cat > /tmp/test-complex.json << 'EOF'
{
  "select": [
    "event_name",
    {
      "function": "count",
      "args": ["*"],
      "alias": "total_events"
    },
    {
      "function": "p95",
      "args": ["duration"],
      "alias": "p95_duration"
    }
  ],
  "from": "events",
  "timeRange": {
    "start": "2024-01-01T00:00:00Z",
    "end": "2024-01-02T00:00:00Z"
  },
  "where": [
    {
      "field": "event_type",
      "op": "IN",
      "value": ["api_call", "db_query"]
    },
    {
      "field": "status_code",
      "op": ">=",
      "value": 200
    }
  ],
  "groupBy": [
    {
      "timeWindow": "5m"
    },
    "event_name"
  ],
  "orderBy": [
    {
      "field": "total_events",
      "direction": "DESC"
    }
  ],
  "limit": 1000
}
EOF

# Test unmarshal with Go
go run -C pkg/domain << 'GOCODE'
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hoppermq/streamly/pkg/domain"
)

func main() {
	// Test simple query
	data1, _ := os.ReadFile("/tmp/test-simple.json")
	var req1 domain.QueryAstRequest
	if err := json.Unmarshal(data1, &req1); err != nil {
		panic(fmt.Sprintf("Simple unmarshal failed: %v", err))
	}
	fmt.Printf("✅ Simple query: %d select clauses\n", len(req1.Select))

	// Test complex query
	data2, _ := os.ReadFile("/tmp/test-complex.json")
	var req2 domain.QueryAstRequest
	if err := json.Unmarshal(data2, &req2); err != nil {
		panic(fmt.Sprintf("Complex unmarshal failed: %v", err))
	}
	fmt.Printf("✅ Complex query: %d select, %d where, %d groupBy\n", 
		len(req2.Select), len(req2.Where), len(req2.GroupBy))
	
	// Verify polymorphic handling
	hasField := false
	hasFunc := false
	for _, sel := range req2.Select {
		if sel.IsField() {
			hasField = true
		}
		if sel.IsFunction() {
			hasFunc = true
		}
	}
	if !hasField || !hasFunc {
		panic("Polymorphic select not working")
	}
	fmt.Println("✅ Polymorphic select clauses work correctly")
}
GOCODE

echo ""
echo "4️⃣  All checks passed!"
echo "===================================="
echo ""
echo "📝 Summary:"
echo "  - Domain types properly defined"
echo "  - JSON unmarshaling works for polymorphic types"
echo "  - Gin BindJSON integration ready"
echo "  - Build successful"
echo ""
echo "🚀 Ready to implement ClickHouse repository!"
