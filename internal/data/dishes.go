package data

type Dish struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Price       Price    `json:"price"`
	Description string   `json:"description"`
	Category    []string `json:"category"`
	Photo       string   `json:"photo,omitempty"`
}
