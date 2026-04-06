package model

import (
	"time"
)

type Outfit struct {
	Name      string
	Occasion  []string
	Season    []string
	Mood      string
	Fragrance string
	Notes     string
	Favorite  bool
	WearCount int
	LastWorn  time.Time
	Items     []Item
}

type OutfitBuilder struct {
	outfit Outfit
}

func NewOutfitBuilder() *OutfitBuilder {
	return &OutfitBuilder{
		outfit: Outfit{
			Occasion: []string{},
			Season:   []string{},
			Items:    []Item{},
		},
	}
}

func (o *OutfitBuilder) WithName(name string) *OutfitBuilder {
	o.outfit.Name = name
	return o
}

func (o *OutfitBuilder) WithOccasion(occasion []string) *OutfitBuilder {
	cp := make([]string, len(occasion))
	copy(cp, occasion)
	o.outfit.Occasion = cp
	return o
}

func (o *OutfitBuilder) AddOccasion(occasion string) *OutfitBuilder {
	o.outfit.Occasion = append(o.outfit.Occasion, occasion)
	return o
}

func (o *OutfitBuilder) WithSeason(season []string) *OutfitBuilder {
	cp := make([]string, len(season))
	copy(cp, season)
	o.outfit.Season = cp
	return o
}

func (o *OutfitBuilder) AddSeason(season string) *OutfitBuilder {
	o.outfit.Season = append(o.outfit.Season, season)
	return o
}

func (o *OutfitBuilder) WithMood(mood string) *OutfitBuilder {
	o.outfit.Mood = mood
	return o
}

func (o *OutfitBuilder) WithFragrance(fragrance string) *OutfitBuilder {
	o.outfit.Fragrance = fragrance
	return o
}

func (o *OutfitBuilder) WithNotes(notes string) *OutfitBuilder {
	o.outfit.Notes = notes
	return o
}

func (o *OutfitBuilder) WithFavorite(favorite bool) *OutfitBuilder {
	o.outfit.Favorite = favorite
	return o
}

func (o *OutfitBuilder) WithWearCount(wearCount int) *OutfitBuilder {
	o.outfit.WearCount = wearCount
	return o
}

func (o *OutfitBuilder) WithLastWorn(lastWorn time.Time) *OutfitBuilder {
	o.outfit.LastWorn = lastWorn
	return o
}

func (o *OutfitBuilder) WithItems(items []Item) *OutfitBuilder {
	cp := make([]Item, len(items))
	copy(cp, items)
	o.outfit.Items = cp
	return o
}

func (o *OutfitBuilder) AddItem(item Item) *OutfitBuilder {
	o.outfit.Items = append(o.outfit.Items, item)
	return o
}

func (o *OutfitBuilder) Build() Outfit {
	return o.outfit
}
