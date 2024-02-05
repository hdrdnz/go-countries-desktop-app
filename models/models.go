package model

import "fyne.io/fyne"

type Countries struct {
	Name        string     `json:"name"`
	Capital     string     `json:"capital"`
	Region      string     `json:"region"`
	Population  int64      `json:"population"`
	NativeName  string     `json:"nativeName"`
	NumericCode string     `json:"numericCode"`
	Flags       Flags      `json:"flags"`
	Currencies  []Currency `json:"currencies"`
	Languages   []Language `json:"languages"`
	Flag        string     `json:"flag"`
}

type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Flags struct {
	SVG string `json:"svg"`
	PNG string `json:"png"`
}

type Language struct {
	Iso6391    string `json:"iso639_1"`
	Iso6392    string `json:"iso639_2"`
	Name       string `json:"name"`
	NativeName string `json:"nativeName"`
}
type Page1 struct {
	content fyne.CanvasObject
}
