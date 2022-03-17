package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"todoapp/logger"
	"todoapp/router"

	"go.uber.org/zap"
)

// main function
func main() {
	sugarLogger := logger.InitLogger()
	router.Router()
	fmt.Printf("Listening & Serving on %s\n", os.Getenv("SERVER_PORT"))
	zap.L().Info(fmt.Sprintf("Listening & Serving on %s", os.Getenv("SERVER_PORT")))
	err := http.ListenAndServe(os.Getenv("SERVER_PORT"), nil)
	if err != nil {
		sugarLogger.Errorf("Failed to start server %s", err.Error())
	}

}
