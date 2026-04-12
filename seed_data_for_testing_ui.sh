#!/bin/bash
# Seed script for pinta-api UI testing
# Creates 22 items and 5 outfits with shared items across outfits
set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"

extract_id() {
  python3 -c "import sys,json; print(json.load(sys.stdin)['id'])"
}

extract_outfit_id() {
  python3 -c "import sys,json; print(json.load(sys.stdin)['outfit'])"
}

echo "=============================="
echo " Seeding pinta-api test data"
echo " Base URL: $BASE_URL"
echo "=============================="

# ─────────────────────────────────────
# ITEMS — TOPS (5)
# ─────────────────────────────────────
echo ""
echo ">>> Creating tops..."

ITEM_WHITE_TEE=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "White Classic Tee",
    "category": "top",
    "color": "White",
    "brand": "Uniqlo",
    "material": "Cotton",
    "season": ["spring", "summer"],
    "occasion": ["casual", "sport"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] White Classic Tee -> $ITEM_WHITE_TEE"

ITEM_BLACK_BLAZER=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Black Blazer",
    "category": "top",
    "color": "Black",
    "brand": "Zara",
    "material": "Polyester",
    "season": ["fall", "winter"],
    "occasion": ["work", "formal"],
    "condition": "new"
  }' | extract_id)
echo "  [✓] Black Blazer -> $ITEM_BLACK_BLAZER"

ITEM_FLORAL_BLOUSE=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Floral Blouse",
    "category": "top",
    "color": "Pink",
    "brand": "Anthropologie",
    "material": "Silk",
    "season": ["spring", "summer"],
    "occasion": ["brunch", "casual", "date"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Floral Blouse -> $ITEM_FLORAL_BLOUSE"

ITEM_NAVY_CREWNECK=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Navy Crewneck Sweater",
    "category": "top",
    "color": "Navy",
    "brand": "Gap",
    "material": "Cotton",
    "season": ["fall", "winter"],
    "occasion": ["casual", "sport"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Navy Crewneck Sweater -> $ITEM_NAVY_CREWNECK"

ITEM_STRIPED_SHIRT=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Striped Button Shirt",
    "category": "top",
    "color": "Blue",
    "brand": "J.Crew",
    "material": "Cotton",
    "season": ["spring", "fall"],
    "occasion": ["work", "casual", "date"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Striped Button Shirt -> $ITEM_STRIPED_SHIRT"

# ─────────────────────────────────────
# ITEMS — BOTTOMS (5)
# ─────────────────────────────────────
echo ""
echo ">>> Creating bottoms..."

ITEM_BLUE_JEANS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Blue Slim Jeans",
    "category": "bottom",
    "color": "Blue",
    "brand": "Levi'\''s",
    "material": "Denim",
    "season": ["spring", "summer", "fall"],
    "occasion": ["casual", "date"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Blue Slim Jeans -> $ITEM_BLUE_JEANS"

ITEM_BLACK_TROUSERS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Black Dress Trousers",
    "category": "bottom",
    "color": "Black",
    "brand": "Hugo Boss",
    "material": "Wool",
    "season": ["fall", "winter"],
    "occasion": ["work", "formal"],
    "condition": "new"
  }' | extract_id)
echo "  [✓] Black Dress Trousers -> $ITEM_BLACK_TROUSERS"

ITEM_BEIGE_CHINOS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Beige Chinos",
    "category": "bottom",
    "color": "Beige",
    "brand": "Banana Republic",
    "material": "Cotton",
    "season": ["spring", "fall"],
    "occasion": ["casual", "work", "brunch", "date"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Beige Chinos -> $ITEM_BEIGE_CHINOS"

ITEM_FLORAL_SKIRT=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Floral Midi Skirt",
    "category": "bottom",
    "color": "Pink",
    "brand": "Zara",
    "material": "Chiffon",
    "season": ["spring", "summer"],
    "occasion": ["brunch", "casual", "date"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Floral Midi Skirt -> $ITEM_FLORAL_SKIRT"

ITEM_GREY_SWEATPANTS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Grey Sweatpants",
    "category": "bottom",
    "color": "Grey",
    "brand": "Nike",
    "material": "Fleece",
    "season": ["spring", "summer", "fall", "winter"],
    "occasion": ["sport", "casual"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Grey Sweatpants -> $ITEM_GREY_SWEATPANTS"

# ─────────────────────────────────────
# ITEMS — SHOES (5)
# ─────────────────────────────────────
echo ""
echo ">>> Creating shoes..."

ITEM_WHITE_SNEAKERS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "White Sneakers",
    "category": "shoes",
    "color": "White",
    "brand": "Adidas",
    "material": "Leather",
    "season": ["spring", "summer"],
    "occasion": ["casual", "sport"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] White Sneakers -> $ITEM_WHITE_SNEAKERS"

ITEM_BLACK_HEELS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Black Heels",
    "category": "shoes",
    "color": "Black",
    "brand": "Steve Madden",
    "material": "Leather",
    "season": ["fall", "winter"],
    "occasion": ["formal", "work", "date"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Black Heels -> $ITEM_BLACK_HEELS"

ITEM_BROWN_LOAFERS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Brown Loafers",
    "category": "shoes",
    "color": "Brown",
    "brand": "Cole Haan",
    "material": "Leather",
    "season": ["spring", "fall"],
    "occasion": ["work", "casual", "brunch", "date"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Brown Loafers -> $ITEM_BROWN_LOAFERS"

ITEM_WHITE_RUNNERS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "White Running Shoes",
    "category": "shoes",
    "color": "White",
    "brand": "New Balance",
    "material": "Mesh",
    "season": ["spring", "summer"],
    "occasion": ["sport"],
    "condition": "new"
  }' | extract_id)
echo "  [✓] White Running Shoes -> $ITEM_WHITE_RUNNERS"

ITEM_NUDE_SANDALS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nude Strappy Sandals",
    "category": "shoes",
    "color": "Nude",
    "brand": "Sam Edelman",
    "material": "Suede",
    "season": ["spring", "summer"],
    "occasion": ["brunch", "casual", "date"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Nude Strappy Sandals -> $ITEM_NUDE_SANDALS"

# ─────────────────────────────────────
# ITEMS — OUTERWEAR (2)
# ─────────────────────────────────────
echo ""
echo ">>> Creating outerwear..."

ITEM_BLACK_COAT=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Black Wool Coat",
    "category": "outerwear",
    "color": "Black",
    "brand": "Massimo Dutti",
    "material": "Wool",
    "season": ["fall", "winter"],
    "occasion": ["work", "formal"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Black Wool Coat -> $ITEM_BLACK_COAT"

ITEM_DENIM_JACKET=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Denim Jacket",
    "category": "outerwear",
    "color": "Blue",
    "brand": "Levi'\''s",
    "material": "Denim",
    "season": ["spring", "fall"],
    "occasion": ["casual", "date"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Denim Jacket -> $ITEM_DENIM_JACKET"

# ─────────────────────────────────────
# ITEMS — JEWELRY (3)
# ─────────────────────────────────────
echo ""
echo ">>> Creating jewelry..."

ITEM_GOLD_NECKLACE=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gold Chain Necklace",
    "category": "jewelry",
    "color": "Gold",
    "brand": "Mejuri",
    "material": "Gold",
    "season": ["spring", "summer", "fall", "winter"],
    "occasion": ["casual", "date", "formal"],
    "condition": "new"
  }' | extract_id)
echo "  [✓] Gold Chain Necklace -> $ITEM_GOLD_NECKLACE"

ITEM_PEARL_EARRINGS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Pearl Drop Earrings",
    "category": "jewelry",
    "color": "White",
    "brand": "Madewell",
    "material": "Pearl",
    "season": ["spring", "summer", "fall", "winter"],
    "occasion": ["formal", "brunch", "work"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Pearl Drop Earrings -> $ITEM_PEARL_EARRINGS"

ITEM_SILVER_HOOPS=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Silver Hoop Earrings",
    "category": "jewelry",
    "color": "Silver",
    "brand": "Pandora",
    "material": "Silver",
    "season": ["spring", "summer", "fall", "winter"],
    "occasion": ["casual", "date", "party"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Silver Hoop Earrings -> $ITEM_SILVER_HOOPS"

# ─────────────────────────────────────
# ITEMS — ACCESSORIES (2)
# ─────────────────────────────────────
echo ""
echo ">>> Creating accessories..."

ITEM_LEATHER_BELT=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Brown Leather Belt",
    "category": "accesories",
    "color": "Brown",
    "brand": "Coach",
    "material": "Leather",
    "season": ["spring", "summer", "fall", "winter"],
    "occasion": ["casual", "work"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Brown Leather Belt -> $ITEM_LEATHER_BELT"

ITEM_CANVAS_TOTE=$(curl -sf -X POST "$BASE_URL/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Canvas Tote Bag",
    "category": "accesories",
    "color": "Beige",
    "brand": "Baggu",
    "material": "Canvas",
    "season": ["spring", "summer"],
    "occasion": ["casual", "brunch"],
    "condition": "good"
  }' | extract_id)
echo "  [✓] Canvas Tote Bag -> $ITEM_CANVAS_TOTE"

# ─────────────────────────────────────
# OUTFITS (5)
# ─────────────────────────────────────
echo ""
echo ">>> Creating outfits..."

# ── Outfit 1: Office Ready ──────────────────────────────
# top: Black Blazer, bottom: Black Dress Trousers, shoes: Black Heels
# outerwear: Black Wool Coat, jewelry: Pearl Drop Earrings, accesories: Brown Leather Belt (shared with outfit 5)
OUTFIT_1=$(curl -sf -X POST "$BASE_URL/outfits" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Office Ready\",
    \"occasion\": [\"work\", \"formal\"],
    \"season\": [\"fall\", \"winter\"],
    \"mood\": \"Professional\",
    \"fragrance\": \"Tom Ford Black Orchid\",
    \"notes\": \"Sharp look for important meetings and presentations\",
    \"favorite\": true,
    \"wear_count\": 3,
    \"items\": [
      {
        \"id\": \"$ITEM_BLACK_BLAZER\",
        \"name\": \"Black Blazer\",
        \"category\": \"top\",
        \"color\": \"Black\",
        \"brand\": \"Zara\",
        \"material\": \"Polyester\",
        \"season\": [\"fall\", \"winter\"],
        \"occasion\": [\"work\", \"formal\"],
        \"condition\": \"new\"
      },
      {
        \"id\": \"$ITEM_BLACK_TROUSERS\",
        \"name\": \"Black Dress Trousers\",
        \"category\": \"bottom\",
        \"color\": \"Black\",
        \"brand\": \"Hugo Boss\",
        \"material\": \"Wool\",
        \"season\": [\"fall\", \"winter\"],
        \"occasion\": [\"work\", \"formal\"],
        \"condition\": \"new\"
      },
      {
        \"id\": \"$ITEM_BLACK_HEELS\",
        \"name\": \"Black Heels\",
        \"category\": \"shoes\",
        \"color\": \"Black\",
        \"brand\": \"Steve Madden\",
        \"material\": \"Leather\",
        \"season\": [\"fall\", \"winter\"],
        \"occasion\": [\"formal\", \"work\", \"date\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_BLACK_COAT\",
        \"name\": \"Black Wool Coat\",
        \"category\": \"outerwear\",
        \"color\": \"Black\",
        \"brand\": \"Massimo Dutti\",
        \"material\": \"Wool\",
        \"season\": [\"fall\", \"winter\"],
        \"occasion\": [\"work\", \"formal\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_PEARL_EARRINGS\",
        \"name\": \"Pearl Drop Earrings\",
        \"category\": \"jewelry\",
        \"color\": \"White\",
        \"brand\": \"Madewell\",
        \"material\": \"Pearl\",
        \"season\": [\"spring\", \"summer\", \"fall\", \"winter\"],
        \"occasion\": [\"formal\", \"brunch\", \"work\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_LEATHER_BELT\",
        \"name\": \"Brown Leather Belt\",
        \"category\": \"accesories\",
        \"color\": \"Brown\",
        \"brand\": \"Coach\",
        \"material\": \"Leather\",
        \"season\": [\"spring\", \"summer\", \"fall\", \"winter\"],
        \"occasion\": [\"casual\", \"work\"],
        \"condition\": \"good\"
      }
    ]
  }" | extract_outfit_id)
echo "  [✓] Outfit 1 - Office Ready -> $OUTFIT_1"

# ── Outfit 2: Weekend Casual ───────────────────────────
# top: White Classic Tee, bottom: Blue Slim Jeans, shoes: White Sneakers
# jewelry: Gold Chain Necklace (shared with outfit 5), accesories: Canvas Tote Bag (shared with outfit 3)
OUTFIT_2=$(curl -sf -X POST "$BASE_URL/outfits" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Weekend Casual\",
    \"occasion\": [\"casual\"],
    \"season\": [\"spring\", \"summer\"],
    \"mood\": \"Relaxed\",
    \"notes\": \"Go-to weekend look for errands and meetups\",
    \"favorite\": false,
    \"wear_count\": 8,
    \"items\": [
      {
        \"id\": \"$ITEM_WHITE_TEE\",
        \"name\": \"White Classic Tee\",
        \"category\": \"top\",
        \"color\": \"White\",
        \"brand\": \"Uniqlo\",
        \"material\": \"Cotton\",
        \"season\": [\"spring\", \"summer\"],
        \"occasion\": [\"casual\", \"sport\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_BLUE_JEANS\",
        \"name\": \"Blue Slim Jeans\",
        \"category\": \"bottom\",
        \"color\": \"Blue\",
        \"brand\": \"Levi's\",
        \"material\": \"Denim\",
        \"season\": [\"spring\", \"summer\", \"fall\"],
        \"occasion\": [\"casual\", \"date\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_WHITE_SNEAKERS\",
        \"name\": \"White Sneakers\",
        \"category\": \"shoes\",
        \"color\": \"White\",
        \"brand\": \"Adidas\",
        \"material\": \"Leather\",
        \"season\": [\"spring\", \"summer\"],
        \"occasion\": [\"casual\", \"sport\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_GOLD_NECKLACE\",
        \"name\": \"Gold Chain Necklace\",
        \"category\": \"jewelry\",
        \"color\": \"Gold\",
        \"brand\": \"Mejuri\",
        \"material\": \"Gold\",
        \"season\": [\"spring\", \"summer\", \"fall\", \"winter\"],
        \"occasion\": [\"casual\", \"date\", \"formal\"],
        \"condition\": \"new\"
      },
      {
        \"id\": \"$ITEM_CANVAS_TOTE\",
        \"name\": \"Canvas Tote Bag\",
        \"category\": \"accesories\",
        \"color\": \"Beige\",
        \"brand\": \"Baggu\",
        \"material\": \"Canvas\",
        \"season\": [\"spring\", \"summer\"],
        \"occasion\": [\"casual\", \"brunch\"],
        \"condition\": \"good\"
      }
    ]
  }" | extract_outfit_id)
echo "  [✓] Outfit 2 - Weekend Casual -> $OUTFIT_2"

# ── Outfit 3: Sunday Brunch ────────────────────────────
# top: Floral Blouse, bottom: Floral Midi Skirt, shoes: Nude Strappy Sandals
# jewelry: Silver Hoop Earrings, accesories: Canvas Tote Bag (shared with outfit 2)
OUTFIT_3=$(curl -sf -X POST "$BASE_URL/outfits" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Sunday Brunch\",
    \"occasion\": [\"brunch\", \"casual\"],
    \"season\": [\"spring\", \"summer\"],
    \"mood\": \"Breezy\",
    \"fragrance\": \"Chloe Eau de Parfum\",
    \"notes\": \"Floral matching set perfect for outdoor brunch\",
    \"favorite\": true,
    \"wear_count\": 2,
    \"items\": [
      {
        \"id\": \"$ITEM_FLORAL_BLOUSE\",
        \"name\": \"Floral Blouse\",
        \"category\": \"top\",
        \"color\": \"Pink\",
        \"brand\": \"Anthropologie\",
        \"material\": \"Silk\",
        \"season\": [\"spring\", \"summer\"],
        \"occasion\": [\"brunch\", \"casual\", \"date\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_FLORAL_SKIRT\",
        \"name\": \"Floral Midi Skirt\",
        \"category\": \"bottom\",
        \"color\": \"Pink\",
        \"brand\": \"Zara\",
        \"material\": \"Chiffon\",
        \"season\": [\"spring\", \"summer\"],
        \"occasion\": [\"brunch\", \"casual\", \"date\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_NUDE_SANDALS\",
        \"name\": \"Nude Strappy Sandals\",
        \"category\": \"shoes\",
        \"color\": \"Nude\",
        \"brand\": \"Sam Edelman\",
        \"material\": \"Suede\",
        \"season\": [\"spring\", \"summer\"],
        \"occasion\": [\"brunch\", \"casual\", \"date\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_SILVER_HOOPS\",
        \"name\": \"Silver Hoop Earrings\",
        \"category\": \"jewelry\",
        \"color\": \"Silver\",
        \"brand\": \"Pandora\",
        \"material\": \"Silver\",
        \"season\": [\"spring\", \"summer\", \"fall\", \"winter\"],
        \"occasion\": [\"casual\", \"date\", \"party\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_CANVAS_TOTE\",
        \"name\": \"Canvas Tote Bag\",
        \"category\": \"accesories\",
        \"color\": \"Beige\",
        \"brand\": \"Baggu\",
        \"material\": \"Canvas\",
        \"season\": [\"spring\", \"summer\"],
        \"occasion\": [\"casual\", \"brunch\"],
        \"condition\": \"good\"
      }
    ]
  }" | extract_outfit_id)
echo "  [✓] Outfit 3 - Sunday Brunch -> $OUTFIT_3"

# ── Outfit 4: Gym Session ──────────────────────────────
# top: Navy Crewneck Sweater, bottom: Grey Sweatpants, shoes: White Running Shoes
OUTFIT_4=$(curl -sf -X POST "$BASE_URL/outfits" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Gym Session\",
    \"occasion\": [\"sport\"],
    \"season\": [\"spring\", \"summer\"],
    \"mood\": \"Energized\",
    \"notes\": \"Comfortable and functional gym outfit\",
    \"favorite\": false,
    \"wear_count\": 12,
    \"items\": [
      {
        \"id\": \"$ITEM_NAVY_CREWNECK\",
        \"name\": \"Navy Crewneck Sweater\",
        \"category\": \"top\",
        \"color\": \"Navy\",
        \"brand\": \"Gap\",
        \"material\": \"Cotton\",
        \"season\": [\"fall\", \"winter\"],
        \"occasion\": [\"casual\", \"sport\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_GREY_SWEATPANTS\",
        \"name\": \"Grey Sweatpants\",
        \"category\": \"bottom\",
        \"color\": \"Grey\",
        \"brand\": \"Nike\",
        \"material\": \"Fleece\",
        \"season\": [\"spring\", \"summer\", \"fall\", \"winter\"],
        \"occasion\": [\"sport\", \"casual\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_WHITE_RUNNERS\",
        \"name\": \"White Running Shoes\",
        \"category\": \"shoes\",
        \"color\": \"White\",
        \"brand\": \"New Balance\",
        \"material\": \"Mesh\",
        \"season\": [\"spring\", \"summer\"],
        \"occasion\": [\"sport\"],
        \"condition\": \"new\"
      }
    ]
  }" | extract_outfit_id)
echo "  [✓] Outfit 4 - Gym Session -> $OUTFIT_4"

# ── Outfit 5: Autumn Date Night ────────────────────────
# top: Striped Button Shirt, bottom: Beige Chinos, shoes: Brown Loafers
# outerwear: Denim Jacket, jewelry: Gold Chain Necklace (shared with outfit 2), accesories: Brown Leather Belt (shared with outfit 1)
OUTFIT_5=$(curl -sf -X POST "$BASE_URL/outfits" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Autumn Date Night\",
    \"occasion\": [\"date\", \"casual\"],
    \"season\": [\"fall\"],
    \"mood\": \"Charming\",
    \"fragrance\": \"Dior Sauvage\",
    \"notes\": \"Smart-casual look for a dinner or evening walk\",
    \"favorite\": true,
    \"wear_count\": 1,
    \"items\": [
      {
        \"id\": \"$ITEM_STRIPED_SHIRT\",
        \"name\": \"Striped Button Shirt\",
        \"category\": \"top\",
        \"color\": \"Blue\",
        \"brand\": \"J.Crew\",
        \"material\": \"Cotton\",
        \"season\": [\"spring\", \"fall\"],
        \"occasion\": [\"work\", \"casual\", \"date\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_BEIGE_CHINOS\",
        \"name\": \"Beige Chinos\",
        \"category\": \"bottom\",
        \"color\": \"Beige\",
        \"brand\": \"Banana Republic\",
        \"material\": \"Cotton\",
        \"season\": [\"spring\", \"fall\"],
        \"occasion\": [\"casual\", \"work\", \"brunch\", \"date\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_BROWN_LOAFERS\",
        \"name\": \"Brown Loafers\",
        \"category\": \"shoes\",
        \"color\": \"Brown\",
        \"brand\": \"Cole Haan\",
        \"material\": \"Leather\",
        \"season\": [\"spring\", \"fall\"],
        \"occasion\": [\"work\", \"casual\", \"brunch\", \"date\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_DENIM_JACKET\",
        \"name\": \"Denim Jacket\",
        \"category\": \"outerwear\",
        \"color\": \"Blue\",
        \"brand\": \"Levi's\",
        \"material\": \"Denim\",
        \"season\": [\"spring\", \"fall\"],
        \"occasion\": [\"casual\", \"date\"],
        \"condition\": \"good\"
      },
      {
        \"id\": \"$ITEM_GOLD_NECKLACE\",
        \"name\": \"Gold Chain Necklace\",
        \"category\": \"jewelry\",
        \"color\": \"Gold\",
        \"brand\": \"Mejuri\",
        \"material\": \"Gold\",
        \"season\": [\"spring\", \"summer\", \"fall\", \"winter\"],
        \"occasion\": [\"casual\", \"date\", \"formal\"],
        \"condition\": \"new\"
      },
      {
        \"id\": \"$ITEM_LEATHER_BELT\",
        \"name\": \"Brown Leather Belt\",
        \"category\": \"accesories\",
        \"color\": \"Brown\",
        \"brand\": \"Coach\",
        \"material\": \"Leather\",
        \"season\": [\"spring\", \"summer\", \"fall\", \"winter\"],
        \"occasion\": [\"casual\", \"work\"],
        \"condition\": \"good\"
      }
    ]
  }" | extract_outfit_id)
echo "  [✓] Outfit 5 - Autumn Date Night -> $OUTFIT_5"

# ─────────────────────────────────────
# SUMMARY
# ─────────────────────────────────────
echo ""
echo "=============================="
echo " Seeding complete!"
echo "=============================="
echo ""
echo "Items created (22):"
echo "  Tops:       $ITEM_WHITE_TEE, $ITEM_BLACK_BLAZER, $ITEM_FLORAL_BLOUSE, $ITEM_NAVY_CREWNECK, $ITEM_STRIPED_SHIRT"
echo "  Bottoms:    $ITEM_BLUE_JEANS, $ITEM_BLACK_TROUSERS, $ITEM_BEIGE_CHINOS, $ITEM_FLORAL_SKIRT, $ITEM_GREY_SWEATPANTS"
echo "  Shoes:      $ITEM_WHITE_SNEAKERS, $ITEM_BLACK_HEELS, $ITEM_BROWN_LOAFERS, $ITEM_WHITE_RUNNERS, $ITEM_NUDE_SANDALS"
echo "  Outerwear:  $ITEM_BLACK_COAT, $ITEM_DENIM_JACKET"
echo "  Jewelry:    $ITEM_GOLD_NECKLACE, $ITEM_PEARL_EARRINGS, $ITEM_SILVER_HOOPS"
echo "  Accesories: $ITEM_LEATHER_BELT, $ITEM_CANVAS_TOTE"
echo ""
echo "Outfits created (5):"
echo "  1. Office Ready       -> $OUTFIT_1  [work/formal, fall/winter]"
echo "  2. Weekend Casual     -> $OUTFIT_2  [casual, spring/summer]"
echo "  3. Sunday Brunch      -> $OUTFIT_3  [brunch, spring/summer]"
echo "  4. Gym Session        -> $OUTFIT_4  [sport, spring/summer]"
echo "  5. Autumn Date Night  -> $OUTFIT_5  [date, fall]"
echo ""
echo "Shared items across outfits:"
echo "  Canvas Tote Bag      -> outfits 2 & 3"
echo "  Gold Chain Necklace  -> outfits 2 & 5"
echo "  Brown Leather Belt   -> outfits 1 & 5"
