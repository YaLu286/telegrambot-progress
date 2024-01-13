package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrambot/progress/keyboards"
	"telegrambot/progress/models"
)

func DisplayStartMessage(bot *tgbotapi.BotAPI, UserID int64, UserLocationID string, CallerID int) {
	var UserLocation models.Location
	UserLocation.ID = UserLocationID
	UserLocation.LoadInfo()
	startCaption := UserLocation.WelcomeText + "С помощью этого бота вы можете ознакомиться с актуальным ассортиментом бутылочного пива и подобрать его по своим собственным предпочтениям\n📞: " + UserLocation.PhoneNumbers + "\n📩: " + UserLocation.Email
	var startImage tgbotapi.InputMediaPhoto
	startImage.Media = tgbotapi.FilePath(UserLocation.ImagePath)
	startImage.Caption = startCaption
	startImage.Type = "photo"
	editMsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      UserID,
			MessageID:   CallerID,
			ReplyMarkup: &keyboards.CommandInlineKeyboard,
		},
		Media: startImage,
	}
	bot.Send(editMsg)
}

func DisplayNotFoundMessage(bot *tgbotapi.BotAPI, UserID int64, CallerID int) {
	var notFoundMsg tgbotapi.InputMediaPhoto
	notFoundMsg.Media = tgbotapi.FilePath("/images/notfound.jpg")
	notFoundMsg.Caption = "Ничего не найдено.Попробуйте изменить или убрать фильтры."
	notFoundMsg.Type = "photo"
	editMsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      UserID,
			MessageID:   CallerID,
			ReplyMarkup: &keyboards.BackKey,
		},
		Media: notFoundMsg,
	}
	bot.Send(editMsg)
}
