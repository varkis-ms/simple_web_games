package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"simple_web_games/internal/service"
	"simple_web_games/internal/utils"
	"simple_web_games/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	router := httprouter.New()
	cfg, err := utils.LoadConfig(".")
	if err != nil {
		logger.WithError(err).Fatal("no config")
	}
	logger.Info("Create router")
	service.SetupCookieStorage(&cfg, &logger)
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
