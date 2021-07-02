package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var Logger  *log.Logger

//init log
func InitLog(err error){
	var(
		file *os.File
		logger *log.Logger
	)
	logger = log.New()
	if file ,err = os.OpenFile("apollo-cron.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil{
		logger.Out = file
	}else{
		logger.Out = os.Stdout
	}

	// set log level (Only log the info severity or above)
	logger.SetLevel(log.InfoLevel)

	Logger = logger
	return

}