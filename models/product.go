package models

type Product struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name"`
	Category string  `json:"category"` // Office furniture, Textile, etc.
	Segment  string  `json:"segment"`  // Table, Curtains, etc.
	Price    float64 `json:"price"`
}
