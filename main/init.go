package main

import (
	nam "Nam/init"
	"Nam/pkg/server"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	interput chan os.Signal

	service server.Service
	engine  *gin.Engine

	setting *nam.Config
	logger  *logrus.Entry

	logpath string
	logmode uint32
)

func init() {
	interput = make(chan os.Signal)

	var config = nam.NewService()

	setting = config.Init()

	if setting.LogToFile {
		logmode = 5
	} else {
		logmode = 4
	}

	logger = nam.SetupLogger(logpath, logmode)
}
