package telegram

import (
	mydb "github.com/StepanShevelev/smh-tg-bot/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func BotServe(bot *tgbotapi.BotAPI, updatesCh tgbotapi.UpdatesChannel, exitCh chan struct{}) {

	defer logrus.Print("shooting BotServe")

	for {
		select {
		case update := <-updatesCh:
			if update.Message == nil {
				continue
			}
			processMessage(bot, update)
		case <-exitCh:
			return

		}
	}

}

func processMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	mux := &sync.Mutex{}

	logrus.WithFields(logrus.Fields{
		"UserName": update.Message.From.UserName,
		"Text":     update.Message.Text,
	}).Info("Message from User")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "help":
			msg.Text = help(update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		case "start":
			msg.Text = help(update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		case "ip":
			msg.Text = ip(update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		case "history":
			msg.Text = history(update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		case "admin_send_all":
			msg.Text = adminSendAll(bot, update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		case "admin_new":
			msg.Text = adminNew(update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		case "admin_delete":
			msg.Text = adminDelete(update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		case "admin_user_history":
			msg.Text = adminUserHistory(update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		default:
			msg.Text = unknown(update.Message.Text,
				update.Message.From, update.Message.Chat.ID)
		}
	} else {
		msg.Text = unknown(update.Message.Text,
			update.Message.From, update.Message.Chat.ID)
	}
	//bot.Send(msg)

	go func() {
		for counter := 0; counter < 5; counter++ {
			mux.Lock()
			_, err := bot.Send(msg)
			mux.Unlock()
			if err != nil {
				logrus.Error(err)
				time.Sleep(time.Second * 1)
				continue
			}
			break
		}

	}()
}

func unknown(text string, user *tgbotapi.User,
	chatId int64) string {

	isAdmin := mydb.Client.CheckUser(user, chatId)
	logrus.WithFields(logrus.Fields{
		"ChatID":   chatId,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  isAdmin,
		"Text":     text,
	}).Warn("Unknown request")
	return "I don't know that command, try \"/help\""
}
