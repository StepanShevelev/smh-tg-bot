package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StepanShevelev/smh-tg-bot/pkg/ipstackapi"
	"strings"

	mydb "github.com/StepanShevelev/smh-tg-bot/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func history(text string, user *tgbotapi.User, chatId int64) string {

	isAdmin := mydb.Client.CheckUser(user, chatId)
	log.WithFields(log.Fields{
		"ChatID":   chatId,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  isAdmin,
		"Text":     text,
	}).Info("history comand")
	if ok := parseHistoryCmd(text); !ok {
		return "Command arguments error"
	}
	ipReqList, err := mydb.Client.GiveUserHistory(user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "History is empty"
	} else if err != nil {
		log.Error(err)
		return "Sorry something goes wrong, try again"
	}
	return userHistoryPp(ipReqList, user.UserName)
}

func userHistoryPp(ipList []string, userName string) string {
	var history []mydb.GlobalHistory
	mydb.Client.Where("ip IN ?", ipList).Find(&history)
	str := "History of user: " + userName + "\n"
	if len(ipList) < 1 {
		str += "Empty"
		return str
	}
	ret := ""
	for _, record := range history {
		var infoIP ipstackapi.InfoIP
		err := json.Unmarshal([]byte(record.IpInfo), &infoIP)
		if err != nil {
			log.Error(err)
		}
		ret += "IP = " + infoIP.Ip + "\n" +
			"Ip type:  " + fmt.Sprintf("%v", infoIP.Iptype) + "\n" +
			"Continent Code: " + fmt.Sprintf("%v", infoIP.ContinentCode) + "\n" +
			"Continent Name: " + fmt.Sprintf("%v", infoIP.ContinentName) + "\n" +
			"Country: " + infoIP.CountryName + "\n" +
			"Country code: " + fmt.Sprintf("%v", infoIP.CountryCode) + "\n" +
			"City: " + infoIP.City + "\n" +
			"Region Code: " + fmt.Sprintf("%v", infoIP.RegionCode) + "\n" +
			"Region Name:  " + fmt.Sprintf("%v", infoIP.RegionName) + "\n" +
			"Latitude:  " + fmt.Sprintf("%v", infoIP.Latitude) + "\n" +
			"Longitude: " + fmt.Sprintf("%v", infoIP.Longitude) + "\n" + "\n"

	}
	return ret
}

func parseHistoryCmd(text string) bool {
	args := strings.Fields(text)
	if len(args) != 1 {
		return false
	}
	return true
}
