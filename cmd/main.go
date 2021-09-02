package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"todoapp/router"

	"go.uber.org/zap"
)

// main function
func main() {
	var sugarLogger *zap.SugaredLogger
	router.Router()
	zap.L().Info(fmt.Sprintf("Listening & Serving on : %s", os.Getenv("APP_PORT")))
	err := http.ListenAndServe(os.Getenv("SERVER_PORT"), nil)
	if err != nil {
		sugarLogger.Errorf("Failed to start server %s", err.Error())
	}
}
