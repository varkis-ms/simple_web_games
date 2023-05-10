package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"simple_web_games/internal/service"
	"simple_web_games/internal/utils"
	"simple_web_games/pkg/logging"
)

// @title Simple Web Games
// @version 0.1
// @description Web service for playing turn-based games

// @contact.name   Markin Sergey
// @contact.email  markin-2002@yandex.ru

// @host 127.0.0.1:8080
// @BasePath /

func main() {
	logger := logging.GetLogger()
	router := httprouter.New()
	cfg, err := utils.LoadConfig(".")
	if err != nil {
		logger.WithError(err).Fatal("no config")
	}
	logger.Info("Create router")
	service.SetupCookieStorage(&cfg)
	logger.Info("Cookie session storage successfully created")
	handler := service.New(logger)
	handler.Register(router)
	runHTTPServer(&cfg, &logger, router)
}

func runHTTPServer(cfg *utils.StorageConfig, logger *logging.Logger, router *httprouter.Router) {
	mux := http.NewServeMux()
	mux.Handle("/", router)
	logger.Infof("Starting application on port :%s", cfg.PortHttp)
	logger.Fatalln(http.ListenAndServe(":"+cfg.PortHttp, mux))
}
