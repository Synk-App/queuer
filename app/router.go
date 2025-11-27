package app

import (
	"net/http"
	"os"
	"synk/gateway/app/controller"
	"synk/gateway/app/util"
)

func Router(service *Service) {
	aboutController := controller.NewAbout(service.DB)
	postController := controller.NewPost(service.DB)

	http.HandleFunc("GET /about", aboutController.HandleAbout)
	http.HandleFunc("POST /send", postController.HandleSend)

	port := os.Getenv("PORT")
	util.Log("app running on port " + port)

	var err error
	env := os.Getenv("ENV")

	if env == "production" {
		util.Log("Running in PRODUCTION mode (HTTP)")
		err = http.ListenAndServe(":"+port, controller.Cors(http.DefaultServeMux))
	} else {
		util.Log("Running in DEV mode (HTTPS)")
		err = http.ListenAndServeTLS(
			":"+port,
			"/cert/cert.pem",
			"/cert/key.pem",
			controller.Cors(http.DefaultServeMux),
		)
	}

	if err != nil {
		util.Log("app failed on running on port " + port + ": " + err.Error())
	}
}
