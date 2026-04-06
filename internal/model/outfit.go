package model

import (
	"time"
)

type Outfit struct {
	Name      string
	Ocassion  []string
	Season    []string
	Mood      string
	Fragance  string
	Notes     string
	Favorite  bool
	WearCount int
	LastWorn  time.Time
	Items     []Item
}
