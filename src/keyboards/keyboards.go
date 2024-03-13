package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var LocationSelectKeys = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Волоколамское шоссе, 1с1", "sokol"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Гиляровского, 68с1", "rizhskaya"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пресненский вал, 38с1", "presnya"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Фрунзенская набережная, 30c5", "frunza"),
	),
)

var CommandInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Список", "list"),
		tgbotapi.NewInlineKeyboardButtonData("Фильтры", "filters"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("К выбору локации", "select_location"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Помощь", "help"),
	),
)

var FiltersSelectKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Стили", "styles"),
		tgbotapi.NewInlineKeyboardButtonData("Пивоварни", "breweries"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Сбросить", "clear"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад в меню", "backToMenu"),
	),
)

var StyleSelectKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("APA", "APA"),
		tgbotapi.NewInlineKeyboardButtonData("Lager", "Lager"),
		tgbotapi.NewInlineKeyboardButtonData("Sour - Fruited", "Sour - Fruited"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("NE Pale Ale", "NE Pale Ale"),
		tgbotapi.NewInlineKeyboardButtonData("Gose", "Gose"),
		tgbotapi.NewInlineKeyboardButtonData("IPA", "IPA"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
	),
)

var BrewerySelectKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

var BackKey = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад в меню", "backToMenu"),
	),
)

var ArrowsKeysFirst = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔽", "right"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад в меню", "backToMenu"),
	),
)

var ArrowsKeys = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔼", "left"),
		tgbotapi.NewInlineKeyboardButtonData("🔽", "right"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад в меню", "backToMenu"),
	),
)

var ArrowsKeysLast = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔼", "left"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад в меню", "backToMenu"),
	),
)

var AdminCommandKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить позицию"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список позиций"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Выйти"),
	),
)

var AdminCreateKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Назад"),
	),
)

var AdminChangeKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Название"),
		tgbotapi.NewKeyboardButton("Пивоварня"),
		tgbotapi.NewKeyboardButton("Стиль"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Краткое описание"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("ABV"),
		tgbotapi.NewKeyboardButton("Рейтинг"),
		tgbotapi.NewKeyboardButton("Цена"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сохранить изменения"),
	),
)

var ActionChoiseKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("❌", "delete"),
		tgbotapi.NewInlineKeyboardButtonData("✅|🚫", "available_switch"),
		tgbotapi.NewInlineKeyboardButtonData("✏️", "change"),
	),
)
