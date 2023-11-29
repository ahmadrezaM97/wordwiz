package models

type Word struct {
	ID       string `json:"-"`
	Lang     string `json:"lang"`
	Word     string `json:"word"`
	Example  string `json:"example"`
	ImageURL string `json:"image_url"`
	Link     string `json:"link"`
}

type Definition struct {
	Lang       string `json:"lang"`
	Definition string `json:"definition"`
}

type AddWordRequest struct {
	Word        Word         `json:"word"`
	Definitions []Definition `json:"definitions"`
}

type UserWord struct {
	Word        Word         `json:"word"`
	Definitions []Definition `json:"definitions"`
}
type UserWords struct {
	Words []UserWord `json:"words"`
}
