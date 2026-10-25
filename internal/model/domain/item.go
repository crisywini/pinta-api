package domain

type Item struct {
	ID                 string
	UserID             string
	Name               string
	Category           Category
	SubCategory        string
	PrimaryColor       string
	SecondaryColor     string
	Pattern            string
	Material           string
	Fit                string
	Silhouette         string
	StyleTags          []string
	Ocassions          []string
	SeasonTags         []string
	ImageAssetID       string
	AvailabilityStatus bool
}

type ItemBuilder struct {
	item Item
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

func NewItemBuilder(name string, category Category) *ItemBuilder {
	return &ItemBuilder{
		item: Item{
			Name:     name,
			Category: category,
		},
	}
}

func (b *ItemBuilder) WithUserID(userID string) *ItemBuilder {
	b.item.UserID = userID
	return b
}

func (b *ItemBuilder) WithName(name string) *ItemBuilder {
	b.item.Name = name
	return b
}

func (b *ItemBuilder) WithCategory(category Category) *ItemBuilder {
	b.item.Category = category
	return b
}

func (b *ItemBuilder) WithSubCategory(subCategory string) *ItemBuilder {
	b.item.SubCategory = subCategory
	return b
}

func (b *ItemBuilder) WithPrimaryColor(primaryColor string) *ItemBuilder {
	b.item.PrimaryColor = primaryColor
	return b
}

func (b *ItemBuilder) WithSecondaryColor(secondaryColor string) *ItemBuilder {
	b.item.SecondaryColor = secondaryColor
	return b
}

func (b *ItemBuilder) WithPattern(pattern string) *ItemBuilder {
	b.item.Pattern = pattern
	return b
}

func (b *ItemBuilder) WithMaterial(material string) *ItemBuilder {
	b.item.Material = material
	return b
}

func (b *ItemBuilder) WithFit(fit string) *ItemBuilder {
	b.item.Fit = fit
	return b
}

func (b *ItemBuilder) WithSilhouette(silhouette string) *ItemBuilder {
	b.item.Silhouette = silhouette
	return b
}

func (b *ItemBuilder) WithStyleTags(styleTags []string) *ItemBuilder {
	b.item.StyleTags = styleTags
	return b
}

func (b *ItemBuilder) WithOcassions(occasions []string) *ItemBuilder {
	b.item.Ocassions = occasions
	return b
}

func (b *ItemBuilder) WithSeasonTags(seasonTags []string) *ItemBuilder {
	b.item.SeasonTags = seasonTags
	return b
}

func (b *ItemBuilder) WithImageAssetID(imageAssetID string) *ItemBuilder {
	b.item.ImageAssetID = imageAssetID
	return b
}

func (b *ItemBuilder) WithAvailabilityStatus(availabilityStatus bool) *ItemBuilder {
	b.item.AvailabilityStatus = availabilityStatus
	return b
}

func (b *ItemBuilder) Build() Item {
	return b.item
}
