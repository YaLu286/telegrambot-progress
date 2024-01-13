package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"slices"
	"telegrambot/progress/models"
)

func UpdateUserBreweries(UserID int64, selectedBrewery string, callback *tgbotapi.CallbackConfig) {

	session := &models.UserSession{
		UserID: UserID,
	}
	models.DB.Find(session)

	if session.Breweries == nil || !slices.Contains(session.Breweries, selectedBrewery) {
		session.AppendBrewery(selectedBrewery)
		callback.Text = "Добавлен фильтр: " + selectedBrewery
	} else {
		session.RemoveBrewery(selectedBrewery)
		callback.Text = "Удалён фильтр: " + selectedBrewery
	}
}

func UpdateUserStyles(UserID int64, selectedStyle string, callback *tgbotapi.CallbackConfig) {

	session := &models.UserSession{
		UserID: UserID,
	}
	session.LoadInfo()

	if session.Styles == nil || !slices.Contains(session.Styles, selectedStyle) {
		session.AppendStyle(selectedStyle)
		callback.Text = "Добавлен фильтр: " + selectedStyle
	} else {
		session.RemoveStyle(selectedStyle)
		callback.Text = "Удалён фильтр: " + selectedStyle
	}
}
