package db

import (
	"encoding/json"
	"errors"
	"fmt"
	cfg "github.com/StepanShevelev/smh-tg-bot/pkg/config"
	ipstack "github.com/StepanShevelev/smh-tg-bot/pkg/ipstackapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Client *client

type client struct {
	*gorm.DB
}

func New(config *cfg.Config) (*client, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=Europe/Moscow",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cant connect to db: %w", err)
	}

	tmp := &client{db}
	Client = tmp
	return tmp, nil
}

func (db *client) GiveUserByID(userId int) (*User, error) {
	user := NewUser()
	result := db.Where("user_id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (db *client) GiveUsers() ([]User, error) {
	var users []User

	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (db *client) GiveIpInfoByIP(ip string) (*ipstack.InfoIP, error) {
	history := NewGlobalHistory()
	result := db.Where("ip = ?", ip).First(&history)
	if result.Error != nil {
		return nil, result.Error
	}
	infoIP := ipstack.NewInfoIP()
	err := json.Unmarshal([]byte(history.IpInfo), infoIP)
	return infoIP, err
}

func (db *client) GiveUserHistoryByIP(ip string, userId int) (*UserHistory, error) {

	ret := NewUserHistory()
	result := db.Where("user_id = ?", userId).
		Where("ip = ?", ip).First(&ret)
	if result.Error != nil {
		return nil, result.Error
	}
	return ret, nil
}

func (db *client) GiveUserHistoryByID(id int) (*UserHistory, error) {

	ret := NewUserHistory()
	result := db.First(&ret, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return ret, nil
}

func (db *client) GiveUserHistory(userId int) ([]string, error) {
	var ips []string

	result := db.Table("user_histories").Where("user_id = ?", userId).
		Select("ip").Find(&ips)
	if result.Error != nil {
		return nil, result.Error
	}
	return ips, nil
}

func (db *client) GiveUserHistoryRet(userId int) ([]ipstack.InfoIP, error) {
	var ret []ipstack.InfoIP
	ips, err := db.GiveUserHistory(userId)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		info, errInfo := db.GiveIpInfoByIP(ip)
		if errInfo != nil {
			continue
		}
		ret = append(ret, *info)
	}
	return ret, nil
}

func (db *client) GiveGlobalUserHistory() ([]UserHistory, error) {
	var ret []UserHistory

	result := db.Find(&ret)
	if result.Error != nil {
		return nil, result.Error
	}
	return ret, nil
}

func (db *client) CreateGlobalHistory(info *ipstack.InfoIP) error {
	b, err := json.Marshal(info)
	if err != nil {
		return err
	}
	db.Create(&GlobalHistory{Ip: info.Ip, IpInfo: string(b)})
	return nil
}

func (db *client) CheckUser(user *tgbotapi.User, chatId int64) bool {
	existUser, err := db.GiveUserByID(user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		db.Create(&User{Name: user.UserName, UserId: user.ID,
			ChatId: chatId, UserRole: "user"})
		log.WithFields(log.Fields{
			"ChatID":   chatId,
			"UserID":   user.ID,
			"UserName": user.UserName,
		}).Info("Add a new user")
		return false
	} else if err != nil {
		log.Error(err)
		return false
	}
	return existUser.UserRole == "admin"
}

func (db *client) SetDB() error {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserHistory{})
	db.AutoMigrate(&GlobalHistory{})
	return nil
}

type User struct {
	gorm.Model
	Name     string
	UserId   int
	ChatId   int64
	UserRole string
}

func NewUser() *User {
	return &User{}
}

type UserHistory struct {
	gorm.Model
	UserId int
	Ip     string
}

func NewUserHistory() *UserHistory {
	return &UserHistory{}
}

type GlobalHistory struct {
	gorm.Model
	Ip     string
	IpInfo string
}

func NewGlobalHistory() *GlobalHistory {
	return &GlobalHistory{}
}
