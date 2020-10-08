package main

import (
	"Nam/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	/* create server engine */
	service = server.New(logger)
	engine = service.Engine()

	/* Add Routers */
	/* --- */

	/* Notify Ginlog */
	for _, route := range engine.Routes() {
		logger.Infof("[ROUTE] %-10v %-30v --> %v", route.Method, route.Path, route.Handler)
	}

	/* Run Server */
	go service.ListenAndServe(interput, setting.HostService.TLS.CertFile, setting.HostService.TLS.KeyFile, &http.Server{
		Addr:         setting.HostService.ParseAddr(),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	})

	/* Listening interput command */
	signal.Notify(interput, os.Interrupt, syscall.SIGTERM)
	println("\n", <-interput)
	logger.Warningln("The service has been shutdown!")
}
