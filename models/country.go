package models

type country struct {
	Name         string
	OfficialName string
	Slug         string
	Flag         string
	Capital      string
	Population   string
	Region       string
	Subregion    string
	Currency     string
	Language     []string
}

type FeaturedCountry struct {
	Name    string
	Slug    string
	Capital string
	Region  string
	Flag    string
}
