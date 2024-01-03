package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"os"
	"slices"
	"strings"
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

func remove_str_from_arr(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

var RedisClient *redis.Client

func update_filters(category string, update *tgbotapi.Update, callback *tgbotapi.CallbackConfig) {

	ctx := context.Background()

	var new_filters_str string
	current_filters_array := strings.Split(RedisClient.HGetAll(ctx, fmt.Sprint(update.CallbackQuery.From.ID)).Val()[category], ",")

	if !slices.Contains(current_filters_array, update.CallbackQuery.Data) {
		current_filters_str := strings.Join(current_filters_array, ",")
		current_filters_array = remove_str_from_arr(current_filters_array, "")
		new_filters_str = strings.Join([]string{current_filters_str, update.CallbackQuery.Data}, ",")
		callback.Text = "Добавлен фильтр: " + update.CallbackQuery.Data
	} else {
		current_filters_array = remove_str_from_arr(current_filters_array, update.CallbackQuery.Data)
		new_filters_str = strings.Join(current_filters_array, ",")
		callback.Text = "Удалён фильтр: " + update.CallbackQuery.Data
	}

	RedisClient.HSet(ctx, fmt.Sprint(update.CallbackQuery.From.ID), category, new_filters_str)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	models.ConnectDatabase()

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	for update := range updates {

		ctx := context.Background()

		if update.Message != nil {

			UserID := update.Message.From.ID
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
				err := RedisClient.HSet(ctx, fmt.Sprint(UserID), "state", "start").Err()
				if err != nil {
					panic(err)
				}
			case "Инфо":
				photo := tgbotapi.NewPhoto(UserID, tgbotapi.FilePath("/Users/yalu/images/progress.jpg"))
				photo.ParseMode = "markdown"
				photo.Caption = "Добро пожаловать в *Прогресс на Соколе*!\nС помощью этого бота вы можете ознакомиться с актуальным ассортиментом бутылочного пива и подобрать его по своим собственным предпочтениям\nТелефон:+7(925)433-52-94\nEmail: progress.sokol@gmail.com"
				bot.Send(photo)

			case "Список":
				var favorite_breweries []string
				var favorite_styles []string
				favorite_breweries = strings.Split(RedisClient.HGetAll(ctx, fmt.Sprint(UserID)).Val()["brewery"], ",")
				favorite_styles = strings.Split(RedisClient.HGetAll(ctx, fmt.Sprint(UserID)).Val()["style"], ",")
				favorite_breweries = remove_str_from_arr(favorite_breweries, "")
				favorite_styles = remove_str_from_arr(favorite_styles, "")

				var bottles []models.Beer

				if len(favorite_breweries) == 0 && len(favorite_styles) == 0 {
					bottles = controllers.FindAllBeer()
					msg.Text = "Фильтры\nПивоварни: -\nСтили: -"
				} else {
					bottles = controllers.FindBeer(favorite_breweries, favorite_styles)
					msg.Text = "Фильтры\nПивоварни: " + strings.Join(favorite_breweries, ", ") + "\nСтили: " + strings.Join(favorite_styles, ", ")
				}
				bot.Send(msg)
				for _, bottle := range bottles {
					bottle_description := fmt.Sprintf("%s от %s \nСтиль: %s\nABV: %.2f Rate: %.2f\n%s\n%d₽", bottle.Name, bottle.Brewery,
						bottle.Style, bottle.ABV,
						bottle.Rate, bottle.Brief, bottle.Price)
					photo := tgbotapi.NewPhoto(update.Message.From.ID, tgbotapi.FilePath(bottle.ImagePath))
					photo.Caption = bottle_description
					if _, err = bot.Send(photo); err != nil {
						panic(err)
					}
				}

			case "Фильтры":
				RedisClient.HSet(ctx, fmt.Sprint(UserID), "state", "filters")
				msg.Text = "Выберите фильтры"
				msg.ReplyMarkup = filtersSelectKeyboard
				bot.Send(msg)

			case "Помощь":

				msg.Text = "Нажмите *Список* для получения списка пива/сидра в бутылках.\nНажмите *Фильтры* для редактирования поисковых фильтров\nНажмите *Инфо* для отображения начального сообщения с информацией\nНажмите *Помощь*, чтобы снова увидеть данное сообщение."
				bot.Send(msg)
			}

		} else if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			UserID := update.CallbackQuery.From.ID

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
					RedisClient.HDel(ctx, fmt.Sprint(UserID), "style")
					RedisClient.HDel(ctx, fmt.Sprint(UserID), "brewery")
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
					update_filters(category, &update, &callback)
				}

			}

			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

		}

	}

}
