package controllers

import (
	"telegrambot/progress/models"
)

func FindAllBeer() []models.Beer {
	var beer_list []models.Beer
	models.DB.Find(&beer_list)
	return beer_list
}

func FindBeer(favorite_breweries []string, favorite_styles []string) []models.Beer {
	var beer_list []models.Beer
	models.DB.Where("brewery IN ? AND style IN ?", favorite_breweries, favorite_styles).Find(&beer_list)
	return beer_list
}
