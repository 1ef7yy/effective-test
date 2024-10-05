package main

import (
	"emobile/pkg/logger"
	"net/http"
	"os"

	"emobile/internal/routes"
	"emobile/internal/view"
)

func main() {
	logger := logger.NewLogger(nil)

	logger.Info("starting server...")

	view := view.NewView(logger)

	logger.Info("initializing router...")

	mux := routes.InitRouter(view)

	logger.Info("server started on " + os.Getenv("SERVER_ADDRESS"))
	if err := http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), mux); err != nil {
		logger.Fatal(err.Error())
	}
}
