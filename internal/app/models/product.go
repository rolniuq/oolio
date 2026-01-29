package models

type Image struct {
	Thumbnail string `json:"thumbnail"`
	Mobile    string `json:"mobile"`
	Tablet    string `json:"tablet"`
	Desktop   string `json:"desktop"`
}

type Product struct {
	ID       string  `json:"id" example:"10"`
	Name     string  `json:"name" example:"Chicken Waffle"`
	Price    float64 `json:"price" description:"Selling price"`
	Category string  `json:"category" example:"Waffle"`
	Image    Image   `json:"image"`
}
