package domain

type Item struct {
	ID                 string   `db:"id" bson:"_id"`
	UserID             string   `db:"user_id" bson:"user_id"`
	Name               string   `db:"name" bson:"name"`
	Category           Category `db:"category" bson:"category"`
	SubCategory        string   `db:"sub_category" bson:"sub_category"`
	PrimaryColor       string   `db:"primary_color" bson:"primary_color"`
	SecondaryColor     string   `db:"secondary_color" bson:"secondary_color"`
	Pattern            string   `db:"pattern" bson:"pattern"`
	Material           string   `db:"material" bson:"material"`
	Fit                string   `db:"fit" bson:"fit"`
	Silhouette         string   `db:"silhouette" bson:"silhouette"`
	StyleTags          []string `db:"style_tags" bson:"style_tags"`
	Ocassions          []string `db:"ocassions" bson:"ocassions"`
	SeasonTags         []string `db:"season_tags" bson:"season_tags"`
	ImageAssetID       string   `db:"image_asset" bson:"image_asset"`
	AvailabilityStatus bool     `db:"avilable" bson:"available"`
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
