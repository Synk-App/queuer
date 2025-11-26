package app

import (
	"net/http"
	"os"
	"synk/gateway/app/controller"
	"synk/gateway/app/util"
)

func Router(service *Service) {
	aboutController := controller.NewAbout(service.DB)
	discordController := controller.NewDiscord(service.DB)
	telegramController := controller.NewTelegram(service.DB)

	http.HandleFunc("GET /about", aboutController.HandleAbout)
	http.HandleFunc("POST /discord/send", discordController.HandleSend)
	http.HandleFunc("POST /telegram/send", telegramController.HandleSend)

	port := os.Getenv("PORT")
	util.Log("app running on port " + port)

	err := http.ListenAndServeTLS(
		":"+port,
		"/cert/cert.pem",
		"/cert/key.pem",
		controller.Cors(http.DefaultServeMux),
	)
	if err != nil {
		util.Log("app failed on running on port " + port + ": " + err.Error())
	}
}
