#!/bin/bash

BASE_URL="http://localhost:8080"

# ─────────────────────────────────────────
# ITEMS
# ─────────────────────────────────────────

# POST /items — Create a new item
curl -s -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "White T-Shirt",
    "category": "top",
    "color": "white",
    "brand": "Uniqlo",
    "material": "cotton",
    "season": ["spring", "summer"],
    "occasion": ["casual"],
    "photo": "https://example.com/white-tshirt.jpg",
    "condition": "new",
    "wear_count": 0
  }' | jq .

echo ""

# GET /items — Get all items
curl -s -X GET "$BASE_URL/items" | jq .

echo ""

# GET /items/:id — Get item by ID  (replace ITEM_ID with a real ObjectID)
ITEM_ID="000000000000000000000001"
curl -s -X GET "$BASE_URL/items/$ITEM_ID" | jq .

echo ""

# PUT /items/:id — Update an item  (replace ITEM_ID with a real ObjectID)
curl -s -X PUT "$BASE_URL/items/$ITEM_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "White T-Shirt (Updated)",
    "category": "top",
    "color": "off-white",
    "brand": "Uniqlo",
    "material": "cotton",
    "season": ["spring", "summer", "autumn"],
    "occasion": ["casual", "work"],
    "photo": "https://example.com/white-tshirt-v2.jpg",
    "condition": "good",
    "wear_count": 5
  }' | jq .

echo ""

# DELETE /items/:id — Delete an item  (replace ITEM_ID with a real ObjectID)
curl -s -X DELETE "$BASE_URL/items/$ITEM_ID" | jq .

echo ""

# ─────────────────────────────────────────
# OUTFITS
# ─────────────────────────────────────────

# POST /outfits — Create a new outfit
# Items referenced inside must exist in the DB (replace item IDs accordingly)
curl -s -X POST "$BASE_URL/outfits" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Summer Casual",
    "occasion": ["casual", "weekend"],
    "season": ["summer"],
    "mood": "relaxed",
    "fragrance": "fresh",
    "notes": "Great for beach days",
    "favorite": true,
    "wear_count": 0,
    "photo": "https://example.com/summer-casual.jpg",
    "last_worn": "2026-04-01T00:00:00Z",
    "items": [
      {
        "id": "000000000000000000000001",
        "name": "White T-Shirt",
        "category": "top",
        "color": "white",
        "brand": "Uniqlo",
        "material": "cotton",
        "season": ["spring", "summer"],
        "occasion": ["casual"],
        "photo": "",
        "condition": "new",
        "wear_count": 0
      }
    ]
  }' | jq .

echo ""

# GET /outfits — Get all outfits
curl -s -X GET "$BASE_URL/outfits" | jq .

echo ""

# GET /outfits/:id — Get outfit by ID  (replace OUTFIT_ID with a real ObjectID)
OUTFIT_ID="000000000000000000000002"
curl -s -X GET "$BASE_URL/outfits/$OUTFIT_ID" | jq .

echo ""

# PUT /outfits/:id — Update an outfit  (replace OUTFIT_ID with a real ObjectID)
curl -s -X PUT "$BASE_URL/outfits/$OUTFIT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Summer Casual (Updated)",
    "occasion": ["casual", "weekend", "beach"],
    "season": ["summer"],
    "mood": "happy",
    "fragrance": "citrus",
    "notes": "Updated notes",
    "favorite": true,
    "wear_count": 3,
    "photo": "https://example.com/summer-casual-v2.jpg",
    "last_worn": "2026-04-07T00:00:00Z",
    "items": []
  }' | jq .

echo ""

# DELETE /outfits/:id — Delete an outfit  (replace OUTFIT_ID with a real ObjectID)
curl -s -X DELETE "$BASE_URL/outfits/$OUTFIT_ID" | jq .
