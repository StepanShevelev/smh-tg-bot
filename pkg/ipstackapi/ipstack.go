package ipstackapi

import (
	"encoding/json"
	cfg "github.com/StepanShevelev/smh-tg-bot/pkg/config"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	log.Info("Receive data from ipstack api")
	r, err := myClient.Get(url)
	if err != nil {
		log.Warn("No response from request")
		return err
	}
	defer r.Body.Close()
	log.Info("Decode recive data to InfoIP struct")
	return json.NewDecoder(r.Body).Decode(target)
}

func NewInfoIP() *InfoIP {
	log.Info("Creating new InfoIP")
	return &InfoIP{}
}

func FillInfoIP(ip string, info *InfoIP) error {
	config := cfg.New()
	if err := config.Load("./configs", "config", "yml"); err != nil {
		log.Fatal(err)
	}
	var url = "http://api.ipstack.com/" + ip +
		"?access_key=" + config.IPStackAccessKey

	return getJson(url, info)
}

type InfoIP struct {
	City          string  `json:"city"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	Ip            string  `json:"ip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	Iptype        string  `json:"type"`
	Zip           string  `json:"zip"`
}
