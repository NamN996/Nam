package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type service struct {
	ginEngine *gin.Engine

	log *logrus.Entry
}

// Service repository
type Service interface {
	Engine() *gin.Engine
	ListenAndServe(c chan<- os.Signal, certFile, keyFile string, serv *http.Server)
}

// New function
func New(log *logrus.Entry) Service {

	gin.SetMode(gin.ReleaseMode)

	var engine = gin.Default()

	// default allow all origins
	handlerFunc := func() gin.HandlerFunc {
		config := cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
			AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "*"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}
		return cors.New(config)
	}()

	engine.Use(handlerFunc)
	engine.NoRoute(func(c *gin.Context) {
		defer func(begin time.Time) {
			log.Warningf("| %v | %15v | %15v | %10v | %v", http.StatusNotFound, time.Since(begin), c.Request.Method, c.Request.RequestURI)
		}(time.Now())
		c.AbortWithStatus(http.StatusNotFound)
	})

	return &service{
		ginEngine: engine,
		log:       log,
	}
}

func (ins *service) Engine() *gin.Engine {
	return ins.ginEngine
}

func (ins *service) ListenAndServe(c chan<- os.Signal, certFile, keyFile string, srv *http.Server) {
	if certFile == "" && keyFile == "" {
		ins.log.Infof("[CLOCK] SYSTEM TIME= [    %v    ]", time.Now().Format(time.RFC3339))
		ins.log.Infof("[ROUTE] Listening and serving HTTP on %v", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			ins.log.Error(err)
		}
	} else {
		ins.log.Infof("[CLOCK] SYSTEM TIME= [    %v    ]", time.Now().Format(time.RFC3339))
		ins.log.Infof("[ROUTE] listening and serving HTTPS on %v", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			ins.log.Error(err)
		}
	}
	c <- os.Interrupt
}
