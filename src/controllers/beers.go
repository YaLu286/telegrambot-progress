package controllers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"telegrambot/progress/models"
)

func FindAllBeer() []models.Beer {
	var beer_list []models.Beer
	models.DB.Find(&beer_list)
	return beer_list
}

func FindBeer(favorite_breweries []string, favorite_styles []string) []models.Beer {
	var beer_list []models.Beer
	if len(favorite_breweries) == 0 {
		models.DB.Where("style IN ?", favorite_styles).Find(&beer_list)
	} else if len(favorite_styles) == 0 {
		models.DB.Where("brewery IN ?", favorite_breweries).Find(&beer_list)
	} else {
		models.DB.Where("brewery IN ? AND style IN ?", favorite_breweries, favorite_styles).Find(&beer_list)
	}
	return beer_list
}

func DisplayBeerist(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	ctx := context.Background()
	UserID := update.Message.From.ID
	msg := tgbotapi.NewMessage(UserID, "")
	var favorite_breweries []string
	var favorite_styles []string
	favorite_breweries = strings.Split(RedisClient.HGetAll(ctx, fmt.Sprint(UserID)).Val()["brewery"], ",")
	favorite_styles = strings.Split(RedisClient.HGetAll(ctx, fmt.Sprint(UserID)).Val()["style"], ",")
	favorite_breweries = RemoveStrFromArray(favorite_breweries, "")
	favorite_styles = RemoveStrFromArray(favorite_styles, "")

	var bottles []models.Beer

	if len(favorite_breweries) == 0 && len(favorite_styles) == 0 {
		bottles = FindAllBeer()
		msg.Text = "Фильтры\nПивоварни: -\nСтили: -"
	} else {
		bottles = FindBeer(favorite_breweries, favorite_styles)
		msg.Text = "Фильтры\nПивоварни: " + strings.Join(favorite_breweries, ", ") + "\nСтили: " + strings.Join(favorite_styles, ", ")
	}
	bot.Send(msg)
	for _, bottle := range bottles {
		bottle_description := fmt.Sprintf("%s от %s \nСтиль: %s\nABV: %.2f Rate: %.2f\n%s\n%d₽", bottle.Name, bottle.Brewery,
			bottle.Style, bottle.ABV,
			bottle.Rate, bottle.Brief, bottle.Price)
		photo := tgbotapi.NewPhoto(update.Message.From.ID, tgbotapi.FilePath(bottle.ImagePath))
		photo.Caption = bottle_description
		if _, err := bot.Send(photo); err != nil {
			panic(err)
		}
	}
}

func CreateBeer(newBeer models.Beer) {
	models.DB.Create(&newBeer)
}

func DeleteBeer(deleteID int64) error {
	return models.DB.Delete(&models.Beer{}, deleteID).Error
}

func ChangeBeer() {

}
