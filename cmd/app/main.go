package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/tommy647/gramarr/internal/bot/telegram"

	"github.com/tommy647/gramarr/internal/app"
	"github.com/tommy647/gramarr/internal/auth"
	"github.com/tommy647/gramarr/internal/bot"
	"github.com/tommy647/gramarr/internal/config"
	"github.com/tommy647/gramarr/internal/conversation"
	"github.com/tommy647/gramarr/internal/radarr"
	"github.com/tommy647/gramarr/internal/router"
	"github.com/tommy647/gramarr/internal/sonarr"
	"github.com/tommy647/gramarr/internal/users"
)

// Flags
var configDir = flag.String("configDir", ".", "config dir for settings and logs")

func main() {
	flag.Parse()

	conf, err := config.LoadConfig(*configDir)
	if err != nil {
		log.Fatalf("failed to load config file: %s", err.Error())
	}

	//err = config.ValidateConfig(conf) // @todo: doesn't do anything
	//if err != nil {
	//	log.Fatal("config error: %s", err.Error())
	//}

	userPath := filepath.Join(*configDir, "users.json")
	usrs, err := users.NewUserDB(userPath)
	if err != nil {
		log.Fatalf("failed to load the usrs db %v", err)
	}

	var rc *radarr.Client
	if conf.Radarr != nil {
		rc, err = radarr.New(*conf.Radarr)
		if err != nil {
			log.Fatalf("failed to create radarr client: %v", err)
		}
	}

	var sn *sonarr.Client
	if conf.Sonarr != nil {
		sn, err = sonarr.New(*conf.Sonarr)
		if err != nil {
			log.Fatalf("failed to create sonarr client: %v", err)
		}
	}

	cm := conversation.NewConversationManager()
	r := router.NewRouter(cm)

	// @todo : move this into our bot service
	tbot, err := telegram.New(conf.Telegram)
	if err != nil {
		log.Fatalf("failed to create telegram bot client: %v", err)
	}

	boter := bot.New(conf.Bot, bot.WithBot(tbot), bot.WithAdmins(usrs.Admins()))

	authoriser := auth.New(conf.Auth, auth.WithBot(boter), auth.WithUsers(usrs))

	a := &app.Service{ // @todo: create contructor
		Auth:   authoriser,
		Bot:    boter,
		Users:  usrs,
		CM:     cm,
		Radarr: rc,
		Sonarr: sn,
	}

	a.SetupHandlers(r)
	log.Print("Gramarr is up and running. Go call your bot!")
	boter.Start()
}
