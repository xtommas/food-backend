package data

type Dish struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Price       int32    `json:"price"`
	Description string   `json:"description"`
	Category    []string `json:"category"`
	Photo       string   `json:"photo"`
}
