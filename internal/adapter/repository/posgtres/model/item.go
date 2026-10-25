package model

type ItemRow struct {
	ID                 string   `db:"id"`
	UserID             string   `db:"user_id"`
	Name               string   `db:"name"`
	Category           string   `db:"category"`
	SubCategory        string   `db:"sub_category"`
	PrimaryColor       string   `db:"primary_color"`
	SecondaryColor     string   `db:"secondary_color"`
	Pattern            string   `db:"pattern"`
	Material           string   `db:"material"`
	Fit                string   `db:"fit"`
	Silhouette         string   `db:"silhouette"`
	StyleTags          []string `db:"style_tags"`
	Occasions          []string `db:"occasions"`
	SeasonTags         []string `db:"season_tags"`
	ImageAssetID       string   `db:"image_asset_id"`
	AvailabilityStatus bool     `db:"available"`
}
