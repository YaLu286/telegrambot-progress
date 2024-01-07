package admin

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	// "net/url"
	"fmt"
	"io"
	"os"
	"strconv"
	"telegrambot/progress/controllers"
	"telegrambot/progress/models"
)

var adminCommandKeyboard = tgbotapi.NewReplyKeyboard(
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

var adminCreateKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Назад"),
	),
)

var adminChangeKeyboard = tgbotapi.NewReplyKeyboard(
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
		tgbotapi.NewKeyboardButton("Стоимость"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Наличие на Пресне"),
		tgbotapi.NewKeyboardButton("Наличие на Рижской"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Наличие на Соколе"),
		tgbotapi.NewKeyboardButton("Наличие на Фрунзе"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сохранить изменения"),
	),
)

var actionChoiseKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Изменить", "change"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "delete"),
	),
)

func Auth(UserID int64) bool {
	var user models.User
	models.DB.First(&user, "id = ?", UserID)
	return user.IsAdmin
}

func SavePhoto(fullURLFile string, fileName string) {

	file, _ := os.Create(fileName)

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, _ := client.Get(fullURLFile)

	defer resp.Body.Close()
	io.Copy(file, resp.Body)
	defer file.Close()
}

func CreateBeerPanel(bot *tgbotapi.BotAPI, admChan chan tgbotapi.Update) {
	var newBeer models.Beer

	var maxIDbeer models.Beer
	models.DB.Last(&maxIDbeer)
	newBeer.ID = maxIDbeer.ID + 1

	update := <-admChan
	msg := tgbotapi.NewMessage(update.Message.From.ID, "Введите название:")
	bot.Send(msg)
	update = <-admChan
	newBeer.Name = update.Message.Text

	msg.Text = "Введите название пивоварни:"
	bot.Send(msg)
	update = <-admChan
	newBeer.Brewery = update.Message.Text

	msg.Text = "Введите стиль пива:"
	bot.Send(msg)
	update = <-admChan
	newBeer.Style = update.Message.Text

	msg.Text = "Введите краткое описание пива:"
	bot.Send(msg)
	update = <-admChan
	newBeer.Brief = update.Message.Text

	msg.Text = "Введите значение ABV:"
	bot.Send(msg)
	update = <-admChan
	newBeer.ABV, _ = strconv.ParseFloat(update.Message.Text, 32)

	msg.Text = "Введите текущий рейтинг пива в Untappd:"
	bot.Send(msg)
	update = <-admChan
	newBeer.Rate, _ = strconv.ParseFloat(update.Message.Text, 32)

	msg.Text = "Введите стоимость пива:"
	bot.Send(msg)
	update = <-admChan
	newBeer.Price, _ = strconv.Atoi(update.Message.Text)

	msg.Text = "Отправьте фотографию пива"
	bot.Send(msg)
	update = <-admChan
	photoURL, _ := bot.GetFileDirectURL(update.Message.Photo[len(update.Message.Photo)-1].FileID)
	var fileName string = "/Users/yalu/images/" + newBeer.Name + ".jpg"
	SavePhoto(photoURL, fileName)
	newBeer.ImagePath = fileName

	controllers.CreateBeer(newBeer)
}

// var adminChangeKeyboard = tgbotapi.NewReplyKeyboard(
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("Название"),
// 		tgbotapi.NewKeyboardButton("Пивоварня"),
// 		tgbotapi.NewKeyboardButton("Стиль"),
// 	),
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("Краткое описание"),
// 	),
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("ABV"),
// 		tgbotapi.NewKeyboardButton("Рейтинг"),
// 		tgbotapi.NewKeyboardButton("Стоимость"),
// 	),
// 	tgbotapi.NewKeyboardButtonRow(

// 	),
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("Сохранить изменения"),

// 	),
// )

func ChangeBeerPanel(bot *tgbotapi.BotAPI, admChan chan tgbotapi.Update, changeID int64, AdminID int64) {
	var Beer models.Beer
	models.DB.Find(&Beer, changeID)

	msg := tgbotapi.NewMessage(AdminID, "Выберите поле, которое хотите изменить")
	msg.ReplyMarkup = adminChangeKeyboard
	bot.Send(msg)

	for {
		update := <-admChan
		if update.Message != nil {
			switch update.Message.Text {
			case "Название":
				msg.Text = "Текущее название:"
				bot.Send(msg)
				msg.Text = Beer.Name
				bot.Send(msg)
				msg.Text = "Введите новое название или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				if update.Message.Text != "-" {
					Beer.Name = update.Message.Text
				}

			case "Пивоварня":
				msg.Text = "Текущее название пивоварни:"
				bot.Send(msg)
				msg.Text = Beer.Brewery
				bot.Send(msg)
				msg.Text = "Введите новое название пивоварни или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				if update.Message.Text != "-" {
					Beer.Brewery = update.Message.Text
				}
			case "Стиль":
				msg.Text = "Текущий стиль:"

				bot.Send(msg)
				msg.Text = Beer.Style
				bot.Send(msg)
				msg.Text = "Введите новое название стиля или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				if update.Message.Text != "-" {
					Beer.Style = update.Message.Text
				}

			case "Краткое описание":
				msg.Text = "Текущее краткое описание:"
				bot.Send(msg)
				msg.Text = Beer.Brief
				bot.Send(msg)
				msg.Text = "Введите новое краткое описание или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				if update.Message.Text != "-" {
					Beer.Brief = update.Message.Text
				}
			case "ABV":
				msg.Text = "Текущее значение ABV:"
				bot.Send(msg)
				msg.Text = strconv.FormatFloat(Beer.ABV, 'f', 2, 32)
				bot.Send(msg)
				msg.Text = "Введите новое значение ABV или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				if update.Message.Text != "-" {
					Beer.ABV, _ = strconv.ParseFloat(update.Message.Text, 32)
				}
			case "Рейтинг":
				msg.Text = "Текущее значение рейтинга на Untappd:"
				bot.Send(msg)
				msg.Text = strconv.FormatFloat(Beer.Rate, 'f', 2, 32)
				bot.Send(msg)
				msg.Text = "Введите новое значение рейтинга на Untappd или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				if update.Message.Text != "-" {
					Beer.Rate, _ = strconv.ParseFloat(update.Message.Text, 32)
				}
			case "Стоимость":
				msg.Text = "Текущая стоимость:"
				bot.Send(msg)
				msg.Text = fmt.Sprint(Beer.Price)
				bot.Send(msg)
				msg.Text = "Введите новую стоимость или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				if update.Message.Text != "-" {
					Beer.Price, _ = strconv.Atoi(update.Message.Text)
				}
			case "Наличие на Пресне":
				msg.Text = "В наличии на Пресне:"
				bot.Send(msg)
				if Beer.Presnya {
					msg.Text = "Да"
				} else {
					msg.Text = "Нет"
				}
				bot.Send(msg)
				msg.Text = "Введите Да, Нет или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				switch update.Message.Text {
				case "Да":
					Beer.Presnya = true
				case "Нет":
					Beer.Presnya = false
				}

			case "Наличие на Рижской":
				if Beer.Rizhskaya {
					msg.Text = "Да"
				} else {
					msg.Text = "Нет"
				}
				bot.Send(msg)
				msg.Text = "Введите Да, Нет или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				switch update.Message.Text {
				case "Да":
					Beer.Rizhskaya = true
				case "Нет":
					Beer.Rizhskaya = false
				}
			case "Наличие на Соколе":
				msg.Text = "В наличии на Соколе:"
				bot.Send(msg)
				if Beer.Sokol {
					msg.Text = "Да"
				} else {
					msg.Text = "Нет"
				}
				bot.Send(msg)
				msg.Text = "Введите Да, Нет или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				switch update.Message.Text {
				case "Да":
					Beer.Sokol = true
				case "Нет":
					Beer.Sokol = false
				}
			case "Наличие на Фрунзе":
				msg.Text = "В наличии на Фрунзенской:"
				bot.Send(msg)
				if Beer.Frunza {
					msg.Text = "Да"
				} else {
					msg.Text = "Нет"
				}
				bot.Send(msg)
				msg.Text = "Введите Да, Нет или - , чтобы пропустить"
				bot.Send(msg)
				update = <-admChan
				switch update.Message.Text {
				case "Да":
					Beer.Frunza = true
				case "Нет":
					Beer.Frunza = false
				}

			case "Сохранить изменения":
				models.DB.Save(&Beer)
				msg.Text = "Позиция сохранена"
				msg.ReplyMarkup = adminCommandKeyboard
				bot.Send(msg)
				return
			}
		}
	}
}

func DisplayBeerListForAdmin(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var bottles []models.Beer

	bottles = controllers.FindAllBeer()
	for _, bottle := range bottles {
		bottle_description := fmt.Sprintf("ID:%d\nНазвание: %s\nПивоварня: %s\nСтиль: %s\nABV: %.2f\nРейтинг:%.2f\nОписание: %s\nСтоимость:%d₽",
			bottle.ID, bottle.Name,
			bottle.Brewery, bottle.Style,
			bottle.ABV, bottle.Rate,
			bottle.Brief, bottle.Price)
		photo := tgbotapi.NewPhoto(update.Message.From.ID, tgbotapi.FilePath(bottle.ImagePath))
		photo.Caption = bottle_description
		photo.ReplyMarkup = actionChoiseKeyboard
		if _, err := bot.Send(photo); err != nil {
			panic(err)
		}
	}
}

func AdmPanel(bot *tgbotapi.BotAPI, admChan chan tgbotapi.Update) {
	for {
		update := <-admChan
		if update.Message != nil {
			UserID := update.Message.From.ID
			switch update.Message.Text {
			case "Добавить позицию":
				msg := tgbotapi.NewMessage(UserID, "Добавление новой позиции")
				bot.Send(msg)
				CreateBeerPanel(bot, admChan)
			case "Список позиций":
				DisplayBeerListForAdmin(bot, update)
			case "Выйти":
				controllers.SetUserState(UserID, "start")
				break
			default:
				msg := tgbotapi.NewMessage(UserID, "Режим администрирования")
				msg.ReplyMarkup = adminCommandKeyboard
				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "change":
				text := update.CallbackQuery.Message.Caption
				var changeID int64
				fmt.Sscanf(text, "ID:%d%s", &changeID, &text)
				msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Изменение позиции")
				bot.Send(msg)
				ChangeBeerPanel(bot, admChan, changeID, update.CallbackQuery.From.ID)
			case "delete":
				text := update.CallbackQuery.Message.Caption
				var deleteID int64
				fmt.Sscanf(text, "ID:%d%s", &deleteID, &text)
				if err := controllers.DeleteBeer(deleteID); err == nil {
					CallbackMsg := tgbotapi.NewCallback(update.CallbackQuery.ID, "Позиция удалена")
					DelMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID)
					bot.Request(CallbackMsg)
					bot.Request(DelMsg)
				}
			}
		}
	}

}
