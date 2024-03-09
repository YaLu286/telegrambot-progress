package admin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"telegrambot/progress/controllers"
	"telegrambot/progress/keyboards"
	"telegrambot/progress/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Auth(session *models.UserSession) bool {
	var adm models.Admin
	res := models.DB.First(&adm, "id = ?", session.UserID)
	if res != nil {
		session.SetUserState("admin")
		session.SetAdminMode(adm.AdminMode)
		return true
	}
	return false
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
	var fileName string = "/images/" + newBeer.Name + ".jpg"
	SavePhoto(photoURL, fileName)
	newBeer.ImagePath = fileName

	newBeer.Create()
}

func ChangeBeerPanel(bot *tgbotapi.BotAPI, admChan chan tgbotapi.Update, changeID int64, AdminID int64) {
	var Beer models.Beer
	models.DB.Find(&Beer, changeID)

	msg := tgbotapi.NewMessage(AdminID, "Выберите поле, которое хотите изменить")
	msg.ReplyMarkup = keyboards.AdminChangeKeyboard
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

			case "Сохранить изменения":
				models.DB.Save(&Beer)
				msg.Text = "Позиция сохранена"
				msg.ReplyMarkup = keyboards.AdminCommandKeyboard
				bot.Send(msg)
				return
			}
		}
	}
}

func DisplayBeerListForAdmin(bot *tgbotapi.BotAPI, update tgbotapi.Update, AdminLocation string) {
	var bottles []models.Beer

	bottles = controllers.FindAllBeerForAdmin()
	for _, bottle := range bottles {
		var availability string
		if slices.Contains(bottle.Availability, AdminLocation) {
			availability = "✅"
		} else {
			availability = "🚫"
		}
		bottle_description := fmt.Sprintf("ID:%d\nНазвание: %s\nПивоварня: %s\nСтиль: %s\nABV: %.2f\nРейтинг:%.2f\nОписание: %s\nСтоимость:%d₽\nВ наличии:%s",
			bottle.ID, bottle.Name,
			bottle.Brewery, bottle.Style,
			bottle.ABV, bottle.Rate,
			bottle.Brief, bottle.Price,
			availability)
		photo := tgbotapi.NewPhoto(update.Message.From.ID, tgbotapi.FilePath(bottle.ImagePath))
		photo.Caption = bottle_description
		photo.ReplyMarkup = keyboards.ActionChoiseKeyboard
		if _, err := bot.Send(photo); err != nil {
			panic(err)
		}
	}
}

func AdmPanel(bot *tgbotapi.BotAPI, admChan chan tgbotapi.Update) {
	defer close(admChan)
	for {
		update := <-admChan
		if update.Message != nil {
			UserID := update.Message.From.ID
			session := &models.UserSession{UserID: UserID}
			session.LoadInfo()
			switch update.Message.Text {
			case "Добавить позицию":
				msg := tgbotapi.NewMessage(UserID, "Добавление новой позиции")
				bot.Send(msg)
				CreateBeerPanel(bot, admChan)
			case "Список позиций":
				DisplayBeerListForAdmin(bot, update, session.AdmMode)
			case "Выйти":
				msg := tgbotapi.NewMessage(session.UserID, "Выход из режима администрирования")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				bot.Send(msg)
				session.SetUserState("start")
				return
			default:
				msg := tgbotapi.NewMessage(UserID, "Режим администрирования")
				msg.ReplyMarkup = keyboards.AdminCommandKeyboard
				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil {
			UserID := update.CallbackQuery.From.ID
			session := &models.UserSession{UserID: UserID}
			session.LoadInfo()
			switch update.CallbackQuery.Data {
			case "change":
				text := update.CallbackQuery.Message.Caption
				var changeID int64
				fmt.Sscanf(text, "ID:%d%s", &changeID, &text)
				msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Изменение позиции")
				bot.Send(msg)
				ChangeBeerPanel(bot, admChan, changeID, update.CallbackQuery.From.ID)
			case "available_switch":
				text := update.CallbackQuery.Message.Caption
				var beerID int64
				fmt.Sscanf(text, "ID:%d", &beerID)
				var beer models.Beer
				models.DB.Find(&beer, beerID)
				beer.SwitchAvailability(session.AdmMode)
				var availability string
				if slices.Contains(beer.Availability, session.AdmMode) {
					availability = "✅"
				} else {
					availability = "🚫"
				}
				re_msg := tgbotapi.NewEditMessageCaption(session.UserID, update.CallbackQuery.Message.MessageID, "")
				re_msg.Caption = fmt.Sprintf("ID:%d\nНазвание: %s\nПивоварня: %s\nСтиль: %s\nABV: %.2f\nРейтинг:%.2f\nОписание: %s\nСтоимость:%d₽\nВ наличии:%s",
					beer.ID, beer.Name,
					beer.Brewery, beer.Style,
					beer.ABV, beer.Rate,
					beer.Brief, beer.Price,
					availability)
				re_msg.ReplyMarkup = &keyboards.ActionChoiseKeyboard
				bot.Send(re_msg)

			case "delete":
				text := update.CallbackQuery.Message.Caption
				var deleteBeer models.Beer
				fmt.Sscanf(text, "ID:%d%s", &deleteBeer.ID, &text)
				if err := deleteBeer.Delete(); err == nil {
					CallbackMsg := tgbotapi.NewCallback(update.CallbackQuery.ID, "Позиция удалена")
					DelMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID)
					bot.Request(CallbackMsg)
					bot.Request(DelMsg)
				}
			}
		}
	}

}
