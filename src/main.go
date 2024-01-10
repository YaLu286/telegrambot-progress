package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"telegrambot/progress/admin"
	"telegrambot/progress/controllers"
	"telegrambot/progress/keyboards"
	"telegrambot/progress/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChanRoute struct {
	AdminID int64
	AdmChan chan tgbotapi.Update
}

func SendUpdateToAdmin(ChanRoutes []ChanRoute, AdminID int64, Update tgbotapi.Update) {
	for _, route := range ChanRoutes {
		if route.AdminID == AdminID {
			ok := true
			select {
			case _, ok = <-route.AdmChan:
			default:
			}
			if ok {
				route.AdmChan <- Update
			}
		}
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	models.ConnectDatabase()

	models.ConnectRedis()

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	var ChanRoutes []ChanRoute

	var beer_map map[int64][]models.Beer = make(map[int64][]models.Beer)

	// var admChan chan tgbotapi.Update = make(chan tgbotapi.Update)

	for update := range updates {

		if update.Message != nil {

			UserID := update.Message.From.ID
			UserState := controllers.GetUserState(UserID)

			if UserState == "admin" {

				SendUpdateToAdmin(ChanRoutes, UserID, update)

			} else {

				msg := tgbotapi.NewMessage(UserID, "")

				msg.ReplyMarkup = keyboards.CommandKeyboard
				msg.ParseMode = "markdown"

				switch update.Message.Text {
				case "/start":

					photo := tgbotapi.NewPhoto(UserID, tgbotapi.FilePath("/images/progress.jpg"))
					photo.ParseMode = "markdown"
					photo.Caption = "Добро пожаловать в *Прогресс*!\nС помощью этого бота вы можете ознакомиться с актуальным ассортиментом бутылочного пива/сидра и подобрать его по своим собственным предпочтениям.\nЧтобы продолжить, пожалуйста, выберите локацию."
					photo.ReplyMarkup = keyboards.LocationSelectKeys
					bot.Send(photo)
					controllers.SetUserState(UserID, "welcome")
					DelMsg := tgbotapi.NewDeleteMessage(UserID, update.Message.MessageID)
					bot.Request(DelMsg)

				case "/admin":
					if admin.Auth(UserID) {
						var NewChanRoute ChanRoute
						NewChanRoute.AdminID = UserID
						NewChanRoute.AdmChan = make(chan tgbotapi.Update)
						ChanRoutes = append(ChanRoutes, NewChanRoute)
						go admin.AdmPanel(bot, ChanRoutes[len(ChanRoutes)-1].AdmChan)
					}
				}
			}

		} else if update.CallbackQuery != nil {

			UserID := update.CallbackQuery.From.ID
			UserState := controllers.GetUserState(UserID)
			CallbackData := update.CallbackQuery.Data

			if UserState == "admin" {

				SendUpdateToAdmin(ChanRoutes, UserID, update)

			} else {

				if update.CallbackQuery.Data == "right" {
					ctx := context.Background()
					next_page, _ := strconv.Atoi(models.RedisClient.HGet(ctx, fmt.Sprint(UserID), "page").Val())
					next_page++
					if next_page < len(beer_map[UserID]) {
						delMsg := tgbotapi.NewDeleteMessage(UserID, update.CallbackQuery.Message.MessageID)
						controllers.DisplayBeer(bot, UserID, &beer_map[UserID][next_page])
						bot.Send(delMsg)
						models.RedisClient.HSet(ctx, fmt.Sprint(UserID), "page", next_page)
					}
				} else if update.CallbackQuery.Data == "left" {
					ctx := context.Background()
					prev_page, _ := strconv.Atoi(models.RedisClient.HGet(ctx, fmt.Sprint(UserID), "page").Val())
					prev_page--
					if prev_page >= 0 {
						delMsg := tgbotapi.NewDeleteMessage(UserID, update.CallbackQuery.Message.MessageID)
						controllers.DisplayBeer(bot, UserID, &beer_map[UserID][prev_page])
						bot.Send(delMsg)
						models.RedisClient.HSet(ctx, fmt.Sprint(UserID), "page", prev_page)
					}
				} else if CallbackData == "presnya" || CallbackData == "rizhskaya" || CallbackData == "sokol" || CallbackData == "frunza" {
					controllers.SetUserLocation(UserID, CallbackData)
					DelMsg := tgbotapi.NewDeleteMessage(UserID, update.CallbackQuery.Message.MessageID)
					controllers.DisplayStartMessage(bot, UserID, CallbackData)
					bot.Request(DelMsg)
				} else if CallbackData == "list" {
					ctx := context.Background()
					var beers []models.Beer
					beers = controllers.GetBeerList(bot, UserID)
					beer_map[UserID] = beers
					models.RedisClient.HSet(ctx, fmt.Sprint(UserID), "page", 0)
					controllers.DisplayBeer(bot, UserID, &beer_map[UserID][0])
					DelMsg := tgbotapi.NewDeleteMessage(UserID, update.CallbackQuery.Message.MessageID)
					bot.Request(DelMsg)
				} else if CallbackData == "filters" {
					msg := tgbotapi.NewMessage(UserID, "")
					msg.Text = "Выберите фильтры"
					msg.ReplyMarkup = keyboards.FiltersSelectKeyboard
					bot.Send(msg)
					DelMsg := tgbotapi.NewDeleteMessage(UserID, update.CallbackQuery.Message.MessageID)
					bot.Request(DelMsg)
				} else if CallbackData == "back" {
					UserLocation := controllers.GetUserLocation(UserID)
					controllers.DisplayStartMessage(bot, UserID, UserLocation)
					DelMsg := tgbotapi.NewDeleteMessage(UserID, update.CallbackQuery.Message.MessageID)
					bot.Request(DelMsg)
				}

				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

				if update.CallbackQuery.Message.Text == "Выберите фильтры" {
					re_msg := tgbotapi.NewEditMessageText(UserID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Message.Text)
					switch update.CallbackQuery.Data {
					case "styles":
						re_msg.Text = "Выберите предпочитаемые стили"
						re_msg.ReplyMarkup = &keyboards.StyleSelectKeyboard
						bot.Send(re_msg)
					case "breweries":
						re_msg.Text = "Выберите предпочитаемые пивоварни"
						re_msg.ReplyMarkup = &keyboards.BrewerySelectKeyboard
						bot.Send(re_msg)
					case "clear":
						controllers.CleanUserFilters(UserID)
						callback.Text = "Фильтры сброшены"
					}
					if _, err := bot.Request(callback); err != nil {
						panic(err)
					}
				}

				if update.CallbackQuery.Message.Text == "Выберите предпочитаемые стили" || update.CallbackQuery.Message.Text == "Выберите предпочитаемые пивоварни" {

					if update.CallbackQuery.Data == "back" {
						re_msg := tgbotapi.NewEditMessageText(UserID, update.CallbackQuery.Message.MessageID, "Выберите фильтры")
						re_msg.ReplyMarkup = &keyboards.FiltersSelectKeyboard
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
						if _, err := bot.Request(callback); err != nil {
							panic(err)
						}
					}

				}

			}
		}

	}

}
