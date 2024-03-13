package controllers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrambot/progress/keyboards"
	"telegrambot/progress/models"
)

func DisplayLocationSelector(bot *tgbotapi.BotAPI, UserID int64, CallerID int) {
	var LocationSelectMsg tgbotapi.InputMediaPhoto
	LocationSelectMsg.Media = tgbotapi.FilePath("/images/progress.jpg")
	LocationSelectMsg.Caption = "Добро пожаловать в *Прогресс*!\nС помощью этого бота вы можете ознакомиться с актуальным ассортиментом бутылочного пива/сидра и подобрать его по своим собственным предпочтениям.\nЧтобы продолжить, пожалуйста, выберите локацию."
	LocationSelectMsg.Type = "photo"
	editMsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      UserID,
			MessageID:   CallerID,
			ReplyMarkup: &keyboards.LocationSelectKeys,
		},
		Media: LocationSelectMsg,
	}
	bot.Send(editMsg)
}

func DisplayStartMessage(bot *tgbotapi.BotAPI, UserID int64, UserLocationID string, CallerID int) {
	var UserLocation models.Location
	UserLocation.ID = UserLocationID
	UserLocation.LoadInfo()
	startCaption := UserLocation.WelcomeText + "\nС помощью этого бота вы можете ознакомиться с актуальным ассортиментом бутылочного пива и подобрать его по своим собственным предпочтениям\n📞: " + UserLocation.PhoneNumbers + "\n📩: " + UserLocation.Email
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

func DisplayBeer(bot *tgbotapi.BotAPI, UserID int64, beer *models.Beer, callerID int, first bool, last bool) {
	beer_description := fmt.Sprintf("%s от %s \nСтиль: %s\nABV: %.2f Rate: %.2f\n%s\n%d₽",
		beer.Name, beer.Brewery,
		beer.Style, beer.ABV,
		beer.Rate, beer.Brief, beer.Price)

	var beerImage tgbotapi.InputMediaPhoto
	beerImage.Media = tgbotapi.FilePath(beer.ImagePath)
	beerImage.Caption = beer_description
	beerImage.Type = "photo"
	editMsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      UserID,
			MessageID:   callerID,
			ReplyMarkup: &keyboards.ArrowsKeys,
		},
		Media: beerImage,
	}
	if first {
		editMsg.BaseEdit.ReplyMarkup = &keyboards.ArrowsKeysFirst
	} else if last {
		editMsg.BaseEdit.ReplyMarkup = &keyboards.ArrowsKeysLast
	} else {
		editMsg.BaseEdit.ReplyMarkup = &keyboards.ArrowsKeys
	}
	if _, err := bot.Request(editMsg); err != nil {
		panic(err)
	}
}

func DisplayHelpMessage(bot *tgbotapi.BotAPI, UserID int64, CallerID int) {
	var helpMsg tgbotapi.InputMediaPhoto
	helpMsg.Media = tgbotapi.FilePath("/images/progress.jpg")
	helpMsg.Caption = "Нажмите 'Список' для отображения всего ассортимента. Нажмите 'Фильтры' для настройки поисковых фильтров по категориям. Нажмите 'К выбору локации', чтобы сменить вашу локацию."
	helpMsg.Type = "photo"
	editMsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      UserID,
			MessageID:   CallerID,
			ReplyMarkup: &keyboards.BackKey,
		},
		Media: helpMsg,
	}
	bot.Send(editMsg)
}
