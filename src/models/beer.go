package models

import (
	"encoding/json"
)

type Beer struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Brewery   string  `json:"brewery"`
	Style     string  `json:"style"`
	Brief     string  `json:"brief"`
	ABV       float64 `json:"abv"`
	Rate      float64 `json:"rate"`
	ImagePath string  `json:"imagepath"`
	Price     int     `json:"price"`
	Presnya   bool    `json:"presnya"`
	Rizhskaya bool    `json:"pizhskaya"`
	Sokol     bool    `json:"sokol"`
	Frunza    bool    `json:"frunza"`
}

func (b *Beer) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (t *Beer) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	return nil
}

func (newBeer *Beer) Create() {
	DB.Create(&newBeer)
}

func (b *Beer) Save() {
	DB.Save(b)
}

func (delBeer *Beer) Delete() error {
	return DB.Delete(delBeer).Error
}

func (b *Beer) Find() {

}
