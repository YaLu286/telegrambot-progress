package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegrambot/progress/keyboards"
	"telegrambot/progress/models"
)

func MainMenuHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, beerMap *map[int64][]models.Beer) {

	UserID, session, CallerMsgID := LoadCallbackInfo(update)

	switch update.CallbackQuery.Data {

	case "list":
		var beers []models.Beer
		beers = GetBeerList(session)
		if len(beers) > 0 {
			(*beerMap)[UserID] = beers
			session.SetCurrentPage(0)
			DisplayBeer(bot, UserID, &(*beerMap)[UserID][0], CallerMsgID)
		} else {
			DisplayNotFoundMessage(bot, UserID, CallerMsgID)
		}
		session.SetUserState("beer_list")
	case "filters":
		msg := tgbotapi.NewEditMessageCaption(UserID, CallerMsgID, "Выберите фильтры")
		msg.ReplyMarkup = &keyboards.FiltersSelectKeyboard
		session.SetUserState("filters_group")
		bot.Send(msg)
	case "select_location":
		DisplayLocationSelector(bot, UserID, CallerMsgID)
		session.SetUserState("location")
	case "help":
		// DisplayHelpMessage(bot, UserID, session.Location, CallerMsgID)
		session.SetUserState("help")
	}
}

func ScrollBeerList(bot *tgbotapi.BotAPI, update tgbotapi.Update, beerMap map[int64][]models.Beer) {

	UserID, session, CallerMsgID := LoadCallbackInfo(update)

	switch update.CallbackQuery.Data {
	case "right":
		if res, next_page := NextPage(session, len(beerMap[UserID])); res {
			DisplayBeer(bot, UserID, &beerMap[UserID][next_page], CallerMsgID)
		}
	case "left":
		if res, prev_page := PreviousPage(session); res {
			DisplayBeer(bot, UserID, &beerMap[UserID][prev_page], CallerMsgID)
		}
	case "backToMenu":
		DisplayStartMessage(bot, UserID, session.Location, CallerMsgID)
		session.SetUserState("main")
	}
}

func NextPage(session *models.UserSession, BeerListLenght int) (bool, int) {
	next_page := session.CurrentPage
	next_page++
	if next_page < BeerListLenght {
		session.SetCurrentPage(next_page)
		return true, next_page
	}
	return false, next_page
}

func PreviousPage(session *models.UserSession) (bool, int) {
	prev_page := session.CurrentPage
	prev_page--
	if prev_page >= 0 {
		session.SetCurrentPage(prev_page)
		return true, prev_page
	}
	return false, prev_page
}

func SelectFiltersGroup(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	UserID, session, CallerMsgID := LoadCallbackInfo(update)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	re_msg := tgbotapi.NewEditMessageCaption(UserID, CallerMsgID,
		update.CallbackQuery.Message.Text)
	switch update.CallbackQuery.Data {
	case "styles":
		re_msg.Caption = "Выберите предпочитаемые стили"
		re_msg.ReplyMarkup = &keyboards.StyleSelectKeyboard
		session.SetUserState("filters")
		bot.Send(re_msg)
	case "breweries":
		re_msg.Caption = "Выберите предпочитаемые пивоварни"
		re_msg.ReplyMarkup = &keyboards.BrewerySelectKeyboard
		session.SetUserState("filters")
		bot.Send(re_msg)
	case "clear":
		session.CleanUserFilters()
		callback.Text = "Фильтры сброшены"
	case "backToMenu":
		DisplayStartMessage(bot, UserID, session.Location, CallerMsgID)
		session.SetUserState("main")
	}
	if _, err := bot.Request(callback); err != nil {
		log.Println(err)
	}
}

func SelectFilters(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	UserID, session, CallerMsgID := LoadCallbackInfo(update)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	if update.CallbackQuery.Data == "back" {
		re_msg := tgbotapi.NewEditMessageCaption(UserID, CallerMsgID, "Выберите фильтры")
		re_msg.ReplyMarkup = &keyboards.FiltersSelectKeyboard
		session.SetUserState("filters_group")
		bot.Send(re_msg)
	} else {
		switch update.CallbackQuery.Message.Caption {
		case "Выберите предпочитаемые стили":
			UpdateUserStyles(UserID, update.CallbackQuery.Data, &callback)
		case "Выберите предпочитаемые пивоварни":
			UpdateUserBreweries(UserID, update.CallbackQuery.Data, &callback)
		}
		if _, err := bot.Request(callback); err != nil {
			log.Println(err)
		}
	}
}

func SelectLocation(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	UserID, session, CallerMsgID := LoadCallbackInfo(update)

	session.SetLocation(update.CallbackQuery.Data)
	session.SetUserState("main")

	DisplayStartMessage(bot, UserID, update.CallbackQuery.Data, CallerMsgID)
}

func LoadCallbackInfo(update tgbotapi.Update) (int64, *models.UserSession, int) {
	UserID := update.CallbackQuery.From.ID
	session := &models.UserSession{UserID: UserID}
	session.LoadInfo()

	CallerMsgID := update.CallbackQuery.Message.MessageID

	return UserID, session, CallerMsgID
}
