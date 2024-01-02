package models

type Beer struct {
	ID        int64
	Name      string
	Brewery   string
	Style     string
	Brief     string
	ABV       float32
	Rate      float32
	ImagePath string
	Price     int16
}
