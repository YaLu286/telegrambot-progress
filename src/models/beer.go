package models

import (
	"encoding/json"
	pq "github.com/lib/pq"
	"slices"
)

type Beer struct {
	ID           int64          `gorm:"primaryKey" json:"id"`
	Name         string         `json:"name"`
	Brewery      string         `json:"brewery"`
	Style        string         `json:"style"`
	Brief        string         `json:"brief"`
	ABV          float64        `json:"abv"`
	Rate         float64        `json:"rate"`
	ImagePath    string         `json:"imagepath"`
	Price        int            `json:"price"`
	Availability pq.StringArray `gorm:"type:text"`
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

func (b *Beer) SwitchAvailability(Location string) {
	if slices.Contains(b.Availability, Location) {
		for i, v := range b.Availability {
			if v == Location {
				b.Availability = append(b.Availability[:i], b.Availability[i+1:]...)
			}
		}
	} else {
		b.Availability = append(b.Availability, Location)
	}
	DB.Save(b)
}
