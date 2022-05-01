package telegram

import (
	"errors"
	"strconv"
	"strings"

	mydb "github.com/StepanShevelev/smh-tg-bot/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func adminSendAll(bot *tgbotapi.BotAPI, text string,
	user *tgbotapi.User, chatId int64) string {

	isAdmin := mydb.Client.CheckUser(user, chatId)
	log.WithFields(log.Fields{
		"ChatID":   chatId,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  isAdmin,
		"Text":     text,
	}).Info("adminSendAll command")
	if !isAdmin {
		return "Permission denied"
	}
	msg, ok := parseAdminSendAllCmd(text)
	if !ok {
		return "Command arguments error"
	}
	users, err := mydb.Client.GiveUsers()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "No users found"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	sendAllUsers(bot, msg, users)
	return "Success! Message sent to all users"
}

func sendAllUsers(bot *tgbotapi.BotAPI, msg string, users []mydb.User) {
	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ChatId, msg)
		bot.Send(msg)
	}
}

func parseAdminSendAllCmd(text string) (string, bool) {
	msg := strings.TrimPrefix(text, "/admin_send_all")
	msg = strings.Trim(msg, " ")
	if args := strings.Fields(msg); len(args) < 1 {
		return "", false
	}
	return msg, true
}

func adminNew(text string, user *tgbotapi.User, chatId int64) string {

	isAdmin := mydb.Client.CheckUser(user, chatId)
	log.WithFields(log.Fields{
		"ChatID":   chatId,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  isAdmin,
		"Text":     text,
	}).Info("adminNew command")
	if !isAdmin {
		return "Permission denied"
	}
	userId, ok := parseAdminCmd("/admin_new", text)
	if !ok {
		return "Command arguments error"
	}
	newAdmin, err := mydb.Client.GiveUserByID(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "Unrecognized user"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	if newAdmin.UserRole == "admin" {
		return "User with id: " + strconv.Itoa(userId) + " - Is already admin"
	}
	newAdmin.UserRole = "admin"
	mydb.Client.Save(&newAdmin)
	return "New admin added with id: " + strconv.Itoa(userId)
}

func adminDelete(text string, user *tgbotapi.User, chatId int64) string {

	isAdmin := mydb.Client.CheckUser(user, chatId)
	log.WithFields(log.Fields{
		"ChatID":   chatId,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  isAdmin,
		"Text":     text,
	}).Info("adminNew command")
	if !isAdmin {
		return "Permission denied"
	}
	userId, ok := parseAdminCmd("/admin_delete", text)
	if !ok {
		return "Command arguments error"
	}
	oldAdimn, err := mydb.Client.GiveUserByID(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "Unrecognized user"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	if oldAdimn.UserRole == "user" {
		return "User with user_id: " + strconv.Itoa(userId) + " - Not admin"
	}
	oldAdimn.UserRole = "user"
	mydb.Client.Save(&oldAdimn)
	return "Deleted admin permissions from user with id : " + strconv.Itoa(userId)
}

func adminUserHistory(text string, user *tgbotapi.User,
	chatId int64) string {

	isAdmin := mydb.Client.CheckUser(user, chatId)
	log.WithFields(log.Fields{
		"ChatID":   chatId,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  isAdmin,
		"Text":     text,
	}).Info("adminUserHistory command")
	if !isAdmin {
		return "Permission denied"
	}
	userId, ok := parseAdminCmd("/admin_user_history", text)
	if !ok {
		return "Command arguments error"
	}
	target, errFu := mydb.Client.GiveUserByID(userId)
	if errors.Is(errFu, gorm.ErrRecordNotFound) {
		return "Unrecognized user"
	} else if errFu != nil {
		log.Error(errFu)
		return "Sorry something goes wrong, try again"
	}
	ipReqList, err := mydb.Client.GiveUserHistory(target.UserId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "History is empty"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	return userHistoryPp(ipReqList, target.Name)
}

func parseAdminCmd(cmd, text string) (int, bool) {
	strId := strings.TrimPrefix(text, cmd)
	strId = strings.Trim(strId, " ")
	if args := strings.Fields(strId); len(args) != 1 {
		return 0, false
	}
	userId, err := strconv.Atoi(strId)
	if err != nil {
		return 0, false
	}
	return userId, true
}
