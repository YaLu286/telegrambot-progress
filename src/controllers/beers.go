package controllers

import (
	// "fmt"
	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// "telegrambot/progress/keyboards"
	"telegrambot/progress/models"
)

func FindAllBeer(Location string) []models.Beer {
	var beer_list []models.Beer
	models.DB.Where(`? IN (select unnest(string_to_array(availability, '"')))`, Location).Find(&beer_list)
	return beer_list
}

func FindAllBeerForAdmin() []models.Beer {
	var beer_list []models.Beer
	models.DB.Find(&beer_list)
	return beer_list
}

func FindBeer(location string, favorite_breweries []string, favorite_styles []string) []models.Beer {
	var beer_list []models.Beer
	args := [][]string{favorite_styles}
	args = append(args, favorite_breweries, favorite_styles)
	if len(favorite_breweries) == 0 {
		models.DB.Where(`style IN ? AND ? IN (select unnest(string_to_array(availability, '"')))`, favorite_styles, location).Find(&beer_list)
	} else if len(favorite_styles) == 0 {
		models.DB.Where(`brewery IN ? AND ? IN (select unnest(string_to_array(availability, '"')))`, favorite_breweries, location).Find(&beer_list)
	} else {
		models.DB.Where(`brewery IN ? AND style IN ? AND ? IN (select unnest(string_to_array(availability, '"')))`, favorite_breweries, favorite_styles, location).Find(&beer_list)
	}
	return beer_list
}

func GetBeerList(session *models.UserSession) []models.Beer {

	var beers []models.Beer

	if len(session.Breweries) == 0 && len(session.Styles) == 0 {
		beers = FindAllBeer(session.Location)
	} else {
		beers = FindBeer(session.Location, session.Breweries, session.Styles)
	}
	return beers
}
