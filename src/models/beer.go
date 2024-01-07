package models

type Beer struct {
	ID        int64
	Name      string
	Brewery   string
	Style     string
	Brief     string
	ABV       float64
	Rate      float64
	ImagePath string
	Price     int
	Presnya   bool
	Rizhskaya bool
	Sokol     bool
	Frunza    bool
}
