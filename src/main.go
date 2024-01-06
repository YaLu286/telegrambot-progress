package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strings"
	"telegrambot/progress/admin"
	"telegrambot/progress/controllers"
	"telegrambot/progress/models"
)

var commandKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список"),
		tgbotapi.NewKeyboardButton("Инфо"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Фильтры"),
		tgbotapi.NewKeyboardButton("Помощь"),
	),
)

var filtersSelectKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Стили", "styles"),
		tgbotapi.NewInlineKeyboardButtonData("Пивоварни", "breweries"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Сбросить", "clear"),
	),
)

var styleSelectKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("APA", "APA"),
		tgbotapi.NewInlineKeyboardButtonData("Lager", "Lager"),
		tgbotapi.NewInlineKeyboardButtonData("Sour - Fruited", "Sour - Fruited"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("NE Pale Ale", "NE Pale Ale"),
		tgbotapi.NewInlineKeyboardButtonData("Gose", "Gose"),
		tgbotapi.NewInlineKeyboardButtonData("Sour - Fruited", "Sour - Fruited"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
	),
)

var brewerySelectKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("AFBrew", "AFBrew"),
		tgbotapi.NewInlineKeyboardButtonData("Zavod", "Zavod"),
		tgbotapi.NewInlineKeyboardButtonData("4Brewers", "4Brewers"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Stamm Brewing", "Stamm Brewing"),
		tgbotapi.NewInlineKeyboardButtonData("Red Button", "Red Button"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	models.ConnectDatabase()

	controllers.ConnectRedis()

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	var admChan chan tgbotapi.Update = make(chan tgbotapi.Update)

	for update := range updates {

		ctx := context.Background()

		if update.Message != nil {

			UserID := update.Message.From.ID
			UserState := controllers.GetUserState(UserID)

			if UserState == "admin" {

				admChan <- update

			} else {

				msg := tgbotapi.NewMessage(UserID, "")

				msg.ReplyMarkup = commandKeyboard
				msg.ParseMode = "markdown"

				switch update.Message.Text {
				case "/start":

					photo := tgbotapi.NewPhoto(UserID, tgbotapi.FilePath("/Users/yalu/images/progress.jpg"))
					photo.ParseMode = "markdown"
					photo.Caption = "Добро пожаловать в *Прогресс на Соколе*!\nС помощью этого бота вы можете ознакомиться с актуальным ассортиментом бутылочного пива/сидра и подобрать его по своим собственным предпочтениям\n📞Телефон:+7(925)433-52-94\n📩Email: progress.sokol@gmail.com"
					photo.ReplyMarkup = commandKeyboard
					bot.Send(photo)
					controllers.SetUserState(UserID, "start")

				case "Инфо":

					photo := tgbotapi.NewPhoto(UserID, tgbotapi.FilePath("/Users/yalu/images/progress.jpg"))
					photo.ParseMode = "markdown"
					photo.Caption = "Добро пожаловать в *Прогресс на Соколе*!\nС помощью этого бота вы можете ознакомиться с актуальным ассортиментом бутылочного пива и подобрать его по своим собственным предпочтениям\nТелефон:+7(925)433-52-94\nEmail: progress.sokol@gmail.com"
					bot.Send(photo)

				case "Список":

					controllers.DisplayBeerist(bot, update)

				case "Фильтры":

					controllers.RedisClient.HSet(ctx, fmt.Sprint(UserID), "state", "filters")
					msg.Text = "Выберите фильтры"
					msg.ReplyMarkup = filtersSelectKeyboard
					bot.Send(msg)

				case "Помощь":

					msg.Text = "Нажмите *Список* для получения списка пива/сидра в бутылках.\nНажмите *Фильтры* для редактирования поисковых фильтров\nНажмите *Инфо* для отображения начального сообщения с информацией\nНажмите *Помощь*, чтобы снова увидеть данное сообщение."
					bot.Send(msg)

				case "/admin":
					if admin.Auth(UserID) {
						// admMode = true
						controllers.SetUserState(UserID, "admin")
						go admin.AdmPanel(bot, admChan)
					}
				}
			}

		} else if update.CallbackQuery != nil {

			UserID := update.CallbackQuery.From.ID
			UserState := controllers.GetUserState(UserID)

			if UserState == "admin" {

				admChan <- update

			} else {

				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

				if update.CallbackQuery.Message.Text == "Выберите фильтры" {
					re_msg := tgbotapi.NewEditMessageText(UserID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Message.Text)
					switch update.CallbackQuery.Data {
					case "styles":
						re_msg.Text = "Выберите предпочитаемые стили"
						re_msg.ReplyMarkup = &styleSelectKeyboard
						bot.Send(re_msg)
					case "breweries":
						re_msg.Text = "Выберите предпочитаемые пивоварни"
						re_msg.ReplyMarkup = &brewerySelectKeyboard
						bot.Send(re_msg)
					case "clear":
						controllers.RedisClient.HDel(ctx, fmt.Sprint(UserID), "style")
						controllers.RedisClient.HDel(ctx, fmt.Sprint(UserID), "brewery")
						callback.Text = "Фильтры сброшены"
					}
				}

				if update.CallbackQuery.Message.Text == "Выберите предпочитаемые стили" || update.CallbackQuery.Message.Text == "Выберите предпочитаемые пивоварни" {

					if update.CallbackQuery.Data == "back" {
						re_msg := tgbotapi.NewEditMessageText(UserID, update.CallbackQuery.Message.MessageID, "Выберите фильтры")
						re_msg.ReplyMarkup = &filtersSelectKeyboard
						bot.Send(re_msg)
					} else {
						var category string
						switch update.CallbackQuery.Message.Text {
						case "Выберите предпочитаемые стили":
							category = "style"
						case "Выберите предпочитаемые пивоварни":
							category = "brewery"
						}
						controllers.UpdateFilters(category, &update, &callback)
					}

				}

				if _, err := bot.Request(callback); err != nil {
					panic(err)
				}

			}
		}

	}

}
