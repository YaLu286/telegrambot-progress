package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// "log"
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

			photo := tgbotapi.NewPhoto(UserID, tgbotapi.FilePath("/images/progress.jpg"))
			photo.ParseMode = "markdown"
			photo.Caption = "Добро пожаловать в *Прогресс*!\nС помощью этого бота вы можете ознакомиться с актуальным ассортиментом бутылочного пива/сидра и подобрать его по своим собственным предпочтениям.\nЧтобы продолжить, пожалуйста, выберите локацию."
			photo.ReplyMarkup = keyboards.LocationSelectKeys
			bot.Send(photo)
			if models.DB.First(session, "user_id = ?", UserID).RowsAffected == 0 {
				session.NewSession(UserID)
			}
			session.SetUserState("location")
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

	switch session.State {
	case "admin":
		SendUpdateToAdmin(*ChanRoutes, UserID, update)
	case "location":
		controllers.SelectLocation(bot, update)
	case "main":
		controllers.MainMenuHandler(bot, update, &beerMap)
	case "beer_list":
		controllers.ScrollBeerList(bot, update, beerMap)
	case "filters_group":
		controllers.SelectFiltersGroup(bot, update)
	case "filters":
		controllers.SelectFilters(bot, update)
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
