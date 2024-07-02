package data

type Recipe struct {
	ROW_ID      int
	Name        string `validate:"required"`
	URL         string
	Ingredients string
	Steps       string
	Notes       string
}
