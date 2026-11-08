package model

import "time"

type OutfitRow struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Occasions []string  `db:"occasions"`
	Seasons   []string  `db:"seasons"`
	Mood      string    `db:"mood"`
	Fragrance string    `db:"fragrance"`
	Notes     string    `db:"notes"`
	Favorite  bool      `db:"favorite"`
	WearCount int       `db:"wear_count"`
	Photo     string    `db:"photo"`
	LastWorn  time.Time `db:"last_worn"`
}
