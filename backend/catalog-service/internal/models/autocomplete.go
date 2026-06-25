package models

// AutocompleteSuggestion — подсказка для поиска
type AutocompleteSuggestion struct {
	Text  string `json:"text"`
	Type  string `json:"type"`  // "product", "category", "tag"
	Slug  string `json:"slug,omitempty"`
	Score int    `json:"score"`
}