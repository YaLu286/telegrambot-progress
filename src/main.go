package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
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

var beerMap map[int64][]models.Beer

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

	beerMap = make(map[int64][]models.Beer)

	for update := range updates {

		if update.Message != nil {
			go UpdateMessageHandler(update, bot, &ChanRoutes)
		} else if update.CallbackQuery != nil {
			go UpdateCallbackHandler(update, bot, &ChanRoutes)
		}

	}
}

func UpdateMessageHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, ChanRoutes *[]ChanRoute) {

	UserID := update.Message.From.ID
	session := &models.UserSession{UserID: UserID}
	session.LoadInfo()

	if session.State == "admin" {

		SendUpdateToAdmin(*ChanRoutes, UserID, update)

	} else {

		switch update.Message.Text {
		case "/start":

			// controllers.DisplayLocationSelector(bot, UserID, update.Message.MessageID)

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
				NewChanRoute := ChanRoute{AdminID: UserID}
				NewChanRoute.AdmChan = make(chan tgbotapi.Update)
				*ChanRoutes = append(*ChanRoutes, NewChanRoute)
				go admin.AdmPanel(bot, (*ChanRoutes)[len(*ChanRoutes)-1].AdmChan)
			}

		}
	}
}

func UpdateCallbackHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, ChanRoutes *[]ChanRoute) {
	UserID := update.CallbackQuery.From.ID
	session := &models.UserSession{UserID: UserID}
	session.LoadInfo()
	CallerMsgID := update.CallbackQuery.Message.MessageID
	CallbackData := update.CallbackQuery.Data
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

	if session.State == "admin" {

		SendUpdateToAdmin(*ChanRoutes, UserID, update)

	} else {

		switch CallbackData {
		case "right":
			if res, next_page := controllers.NextPage(session, len(beerMap[UserID])); res {
				controllers.DisplayBeer(bot, UserID, &beerMap[UserID][next_page], CallerMsgID)
			}
		case "left":
			if res, prev_page := controllers.PreviousPage(session); res {
				controllers.DisplayBeer(bot, UserID, &beerMap[UserID][prev_page], CallerMsgID)
			}
		case "presnya", "rizhskaya", "sokol", "frunza":
			session.SetLocation(CallbackData)
			controllers.DisplayStartMessage(bot, UserID, CallbackData, CallerMsgID)
		case "list":
			var beers []models.Beer
			beers = controllers.GetBeerList(session)
			if len(beers) > 0 {
				beerMap[UserID] = beers
				session.SetCurrentPage(0)
				controllers.DisplayBeer(bot, UserID, &beerMap[UserID][0], CallerMsgID)
			} else {
				controllers.DisplayNotFoundMessage(bot, UserID, CallerMsgID)
			}
		case "select_location":
			controllers.DisplayLocationSelector(bot, UserID, CallerMsgID)
		case "filters":
			msg := tgbotapi.NewEditMessageCaption(UserID, CallerMsgID, "Выберите фильтры")
			msg.ReplyMarkup = &keyboards.FiltersSelectKeyboard
			bot.Send(msg)
		case "backToMenu":

			controllers.DisplayStartMessage(bot, UserID, session.Location, CallerMsgID)

		}
		switch update.CallbackQuery.Message.Caption {
		case "Выберите фильтры":
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
				log.Println(err)
			}
		case "Выберите предпочитаемые стили", "Выберите предпочитаемые пивоварни":
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
					log.Println(err)
				}
			}
		}
	}
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
