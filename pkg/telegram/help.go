package telegram

import (
	mydb "github.com/StepanShevelev/smh-tg-bot/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func help(text string, user *tgbotapi.User, chatId int64) string {
	isAdmin := mydb.Client.CheckUser(user, chatId)
	log.WithFields(log.Fields{
		"ChatID":   chatId,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  isAdmin,
		"Text":     text,
	}).Info("help command")
	if isAdmin {
		return helpAdmin()
	}
	return helpUser()
}

func helpAdmin() string {
	msg := helpUser() + "\n"
	msg += ` COMMANDS FOR ADMIN ` + "\n"
	msg += `"/admin_new [user_id]"` + "\n"
	msg += `Gives admin permissions to user with <user_id>` + "\n" + "\n"

	msg += `"/admin_delete [user_id]"` + "\n"
	msg += `Deletes admin permissions from user with <user_id>` + "\n" + "\n"

	msg += `"/admin_user_history [user_id]"` + "\n"
	msg += `Shows all "/ip" request results of user with <user_id>` + "\n" + "\n"

	msg += `"/admin_send_all [msg]" ` + "\n"
	msg += `Sends <some message> to all bot users` + "\n"
	return msg
}

func helpUser() string {
	msg := `COMMANDS FOR USER ` + "\n"
	msg += `"/ip [some_ipV4]" ` + "\n"
	msg += `Shows info about <some_ipV4>` + "\n" + "\n"
	msg += `"/history"` + "\n"
	msg += `Shows all your requested IPs` + "\n"
	msg += `--------------------------------------------`
	return msg
}
