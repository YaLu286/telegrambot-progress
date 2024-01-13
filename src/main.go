package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"telegrambot/progress/admin"
	"telegrambot/progress/controllers"
	"telegrambot/progress/keyboards"
	"telegrambot/progress/models"
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

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	var ChanRoutes []ChanRoute

	var beer_map map[int64][]models.Beer = make(map[int64][]models.Beer)

	// var admChan chan tgbotapi.Update = make(chan tgbotapi.Update)

	for update := range updates {

		if update.Message != nil {

			UserID := update.Message.From.ID
			session := &models.UserSession{UserID: UserID}
			session.LoadInfo()

			if session.State == "admin" {

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
					if models.DB.First(session, "user_id = ?", UserID).RowsAffected == 0 {
						session.NewSession(UserID)
					}
					session.SetUserState("welcome")
					DelMsg := tgbotapi.NewDeleteMessage(UserID, update.Message.MessageID)
					bot.Request(DelMsg)

				case "/admin":
					if admin.Auth(session) {
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
			session := &models.UserSession{UserID: UserID}
			session.LoadInfo()
			CallerMsgID := update.CallbackQuery.Message.MessageID
			CallbackData := update.CallbackQuery.Data
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

			if session.State == "admin" {

				SendUpdateToAdmin(ChanRoutes, UserID, update)

			} else {

				if update.CallbackQuery.Data == "right" {
					next_page := session.CurrentPage
					next_page++
					if next_page < len(beer_map[UserID]) {
						controllers.DisplayBeer(bot, UserID, &beer_map[UserID][next_page], update.CallbackQuery.Message.MessageID)
						session.SetCurrentPage(next_page)
					}
				} else if update.CallbackQuery.Data == "left" {
					prev_page := session.CurrentPage
					prev_page--
					if prev_page >= 0 {
						controllers.DisplayBeer(bot, UserID, &beer_map[UserID][prev_page], update.CallbackQuery.Message.MessageID)
						session.SetCurrentPage(prev_page)
					}
				} else if CallbackData == "presnya" || CallbackData == "rizhskaya" || CallbackData == "sokol" || CallbackData == "frunza" {
					session.SetLocation(CallbackData)
					controllers.DisplayStartMessage(bot, UserID, CallbackData, CallerMsgID)
				} else if CallbackData == "list" {
					var beers []models.Beer
					beers = controllers.GetBeerList(UserID)
					if len(beers) > 0 {
						beer_map[UserID] = beers
						session.SetCurrentPage(0)
						controllers.DisplayBeer(bot, UserID, &beer_map[UserID][0], CallerMsgID)
					} else {
						controllers.DisplayNotFoundMessage(bot, UserID, CallerMsgID)
					}
				} else if CallbackData == "filters" {
					msg := tgbotapi.NewEditMessageCaption(UserID, CallerMsgID, "Выберите фильтры")
					msg.ReplyMarkup = &keyboards.FiltersSelectKeyboard
					bot.Send(msg)
				} else if CallbackData == "backToMenu" {
					controllers.DisplayStartMessage(bot, UserID, session.Location, CallerMsgID)
				} else if update.CallbackQuery.Message.Caption == "Выберите фильтры" {
					re_msg := tgbotapi.NewEditMessageCaption(UserID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Message.Text)
					switch update.CallbackQuery.Data {
					case "styles":
						re_msg.Caption = "Выберите предпочитаемые стили"
						re_msg.ReplyMarkup = &keyboards.StyleSelectKeyboard
						bot.Send(re_msg)
					case "breweries":
						re_msg.Caption = "Выберите предпочитаемые пивоварни"
						re_msg.ReplyMarkup = &keyboards.BrewerySelectKeyboard
						bot.Send(re_msg)
					case "clear":
						session.CleanUserFilters()
						callback.Text = "Фильтры сброшены"
					}
					if _, err := bot.Request(callback); err != nil {
						panic(err)
					}
				} else if update.CallbackQuery.Message.Caption == "Выберите предпочитаемые стили" || update.CallbackQuery.Message.Caption == "Выберите предпочитаемые пивоварни" {

					if update.CallbackQuery.Data == "back" {
						re_msg := tgbotapi.NewEditMessageCaption(UserID, update.CallbackQuery.Message.MessageID, "Выберите фильтры")
						re_msg.ReplyMarkup = &keyboards.FiltersSelectKeyboard
						bot.Send(re_msg)
					} else {
						switch update.CallbackQuery.Message.Caption {
						case "Выберите предпочитаемые стили":
							controllers.UpdateUserStyles(UserID, CallbackData, &callback)
						case "Выберите предпочитаемые пивоварни":
							controllers.UpdateUserBreweries(UserID, CallbackData, &callback)
						}
						if _, err := bot.Request(callback); err != nil {
							panic(err)
						}
					}

				}

			}
		}

	}

}
