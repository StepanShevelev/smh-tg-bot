package main

import (
	tgapi "github.com/StepanShevelev/smh-tg-bot/pkg/api"
	cfg "github.com/StepanShevelev/smh-tg-bot/pkg/config"
	mydb "github.com/StepanShevelev/smh-tg-bot/pkg/db"
	tlg "github.com/StepanShevelev/smh-tg-bot/pkg/telegram"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	c := make(chan os.Signal, 1)

	config := cfg.New()
	if err := config.Load("./configs", "config", "yml"); err != nil {
		log.Fatal(err)
	}

	db, err := mydb.New(config)
	if err != nil {
		log.Fatal(err)
	}
	db.SetDB()
	bot, update := tlg.BotInit(config)
	exitCh := make(chan struct{}, 1)
	go tlg.BotServe(bot, update, exitCh)
	tgapi.InitBackendApi(config)
	http.ListenAndServe(":"+config.Port, nil)

	signal.Notify(c, os.Interrupt)
	<-c
	exitCh <- struct{}{}
}
