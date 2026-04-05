package model

type Item struct {
	Name      string
	Category  Category
	Color     string
	Brand     string
	Material  string
	Season    []string
	Occasion  []string
	Photo     string
	Condition string
	WearCount int
}

type ItemBuilder struct {
	item Item
}

func NewItemBuilder(name string, category Category) *ItemBuilder {
	return &ItemBuilder{
		item: Item{
			Name:     name,
			Category: category,
		},
	}
}

func (b *ItemBuilder) WithColor(color string) *ItemBuilder {
	b.item.Color = color
	return b
}

func (b *ItemBuilder) WithBrand(brand string) *ItemBuilder {
	b.item.Brand = brand
	return b
}

func (b *ItemBuilder) WithMaterial(material string) *ItemBuilder {
	b.item.Material = material
	return b
}

func (b *ItemBuilder) WithSeason(season []string) *ItemBuilder {
	b.item.Season = season
	return b
}

func (b *ItemBuilder) WithOccasion(occasion []string) *ItemBuilder {
	b.item.Occasion = occasion
	return b
}

func (b *ItemBuilder) WithPhoto(photo string) *ItemBuilder {
	b.item.Photo = photo
	return b
}

func (b *ItemBuilder) WithCondition(condition string) *ItemBuilder {
	b.item.Condition = condition
	return b
}

func (b *ItemBuilder) WithWearCount(count int) *ItemBuilder {
	b.item.WearCount = count
	return b
}

func (b *ItemBuilder) Build() Item {
	return b.item
}

type Category string

const (
	Top        Category = "top"
	Bottom     Category = "bottom"
	Shoes      Category = "shoes"
	Hosiery    Category = "hosiery"
	Outerwear  Category = "outerwear"
	Accesories Category = "accesories"
	Jewelry    Category = "jewelry"
	Innerwear  Category = "Innerwear"
	Fragance   Category = "fragance"
)

func (c Category) IsValid() bool {
	switch c {
	case Top, Bottom, Shoes, Hosiery, Outerwear, Accesories, Jewelry, Innerwear, Fragance:
		return true
	}
	return false
}
