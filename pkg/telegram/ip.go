package telegram

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	mydb "github.com/StepanShevelev/smh-tg-bot/pkg/db"
	ipstack "github.com/StepanShevelev/smh-tg-bot/pkg/ipstackapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ip(text string, user *tgbotapi.User, chatId int64) string {
	isAdmin := mydb.Client.CheckUser(user, chatId)
	log.WithFields(log.Fields{
		"ChatID":   chatId,
		"UserID":   user.ID,
		"UserName": user.UserName,
		"IsAdmin":  isAdmin,
		"Text":     text,
	}).Info("ip command")
	ip, ok := parseIpCmd(text)
	if !ok {
		return "Command arguments error"
	}
	if ipInfo, isSaved := checkIp(ip); !isSaved {
		info := ipstack.NewInfoIP()
		if err := ipstack.FillInfoIP(ip, info); err != nil {
			log.Error(err)
			return "Sorry something goes wrong, try again"
		}
		if err := mydb.Client.CreateGlobalHistory(info); err != nil {
			log.Error(err)
		} else {
			log.WithFields(log.Fields{
				"ip": ip,
			}).Info("New IP Info save in Data Base")
		}
		checkUserHistory(user, ip)
		return ipInfoPp(info)
	} else {
		log.WithFields(log.Fields{
			"ip": ip,
		}).Info("IP Info exist in Data Base")
		checkUserHistory(user, ip)
		return ipInfoPp(ipInfo)
	}
}

func parseIpCmd(text string) (string, bool) {
	args := strings.Fields(text)
	if len(args) != 2 || !validIP4(args[1]) {
		return "", false
	}
	return args[1], true
}

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func checkIp(ip string) (*ipstack.InfoIP, bool) {
	existInfo, err := mydb.Client.GiveIpInfoByIP(ip)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false
	} else if err != nil {
		log.Error(err)
		return nil, false
	}
	return existInfo, true
}

func ipInfoPp(info *ipstack.InfoIP) string {
	if info.CountryName == "" {
		ret := "ip = " + info.Ip + "\n" +
			"Not found"
		return ret
	}
	ret := "IP = " + info.Ip + "\n" +
		"Ip type:  " + fmt.Sprintf("%v", info.Iptype) + "\n" +
		"Continent Code: " + fmt.Sprintf("%v", info.ContinentCode) + "\n" +
		"Continent Name: " + fmt.Sprintf("%v", info.ContinentName) + "\n" +
		"Country: " + info.CountryName + "\n" +
		"Country code: " + fmt.Sprintf("%v", info.CountryCode) + "\n" +
		"City: " + info.City + "\n" +
		"Region Code: " + fmt.Sprintf("%v", info.RegionCode) + "\n" +
		"Region Name:  " + fmt.Sprintf("%v", info.RegionName) + "\n" +
		"Latitude:  " + fmt.Sprintf("%v", info.Latitude) + "\n" +
		"Longitude: " + fmt.Sprintf("%v", info.Longitude) + "\n"
	return ret
}

func checkUserHistory(user *tgbotapi.User, ip string) {
	_, err := mydb.Client.GiveUserHistoryByIP(ip, user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		mydb.Client.Create(&mydb.UserHistory{UserId: user.ID, Ip: ip})
		log.WithFields(log.Fields{
			"UserID": user.ID,
			"IP":     ip,
		}).Info("Add ip in user history")
		return
	} else if err != nil {
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"UserID": user.ID,
		"IP":     ip,
	}).Info("Ip already exists in user history")
}
