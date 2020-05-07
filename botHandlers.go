package main

import (
	"findTuEnvioBot/products"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (m MyBot) handleCallBackQuery(query *tgbotapi.CallbackQuery) {
	_, err := m.bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, query.Data))
	if err != nil {
		logrus.Warn(err)
	}
	m.sendUserPanel(query.Message.Chat.ID, fmt.Sprintf("Ha seleccionado %s", query.Data))

	err = InsertUser(strconv.FormatInt(query.Message.Chat.ID, 10), query.Data)
	if err != nil {
		logrus.Warn(err)
	}

	_, err = m.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
		ChatID:    query.Message.Chat.ID,
		MessageID: query.Message.MessageID,
	})
	if err != nil {
		logrus.Warn(err)
	}
}

func (m MyBot) handlePublicMessage(message *tgbotapi.Message) {

}

func (m MyBot) handlePrivateMessage(privateMsg *tgbotapi.Message) {
	if privateMsg.IsCommand() {
		switch {
		case privateMsg.Text == "/start", privateMsg.Text == "/help":
			// Send instructions
		case strings.Split(privateMsg.Text, " ")[0] == "/buscar":
			// Buscar
		case strings.Split(privateMsg.Text, " ")[0] == "/subscribirme":
			// Buscar
		default:
			// inserte un comando valido
		}
		return
	}

	switch privateMsg.Text {
	case "🆘 Help":
	// Send instuctions
	case "🗺 Seleccionar Provincia":
		// Send province list
		m.sendInlineKeyboardSelectProvince(privateMsg.Chat.ID)
		break
	case "📋 Adicionar subscripcion":
	// Add subscription
	case "👤 Mi Perfil":
	default:
		// Insert a valid message
	}
}

func (m MyBot) handleInlineQuery(query *tgbotapi.InlineQuery) {
	if len(query.Query) >= 2 {
		user, err := GetUser(strconv.Itoa(query.From.ID))
		switch {
		case err == errValEmpty || err == errBucketEmpty:
			_, err := m.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
				InlineQueryID: query.ID,
				Results: []interface{}{
					tgbotapi.NewInlineQueryResultArticleHTML(uuid.New().String(), "Necesitas empezar una"+
						" conversacion conmigo primero.", "Necesitas empezar una conversacion "+
						"conmigo para poder usarme <a href=\"https://t.me/buscarTuEnvioBot\">Empezar conversacion.</a>"),
				},
			})
			if err != nil {
				logrus.Warn(err)
			}
			break
		case err != nil:
			logrus.Warn(err)
			return
		default:

			var productList = make([]products.Product, 0)
			for _, store := range provinces[user.Province] {
				prods, err := m.GetProductsByPattern(store.rawName, query.Query)
				if err != nil {
					logrus.Print(err)
				}

				productList = append(productList, prods...)
			}
			m.sendQueryResultList(productList, query.ID)
		}
	}
}
