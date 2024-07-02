package data

type Recipe struct {
	Name        string `validate:"required"`
	URL         string
	Ingredients string
	Steps       string
	Notes       string
}
