package main

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"softarch/initializers"
	"softarch/middleware"
	"softarch/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "access-control-allow-origin", "access-control-allow-headers"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	router.Use(middleware.TokenMiddleware())

	router.ForwardedByClientIP = true

	if router.SetTrustedProxies([]string{"127.0.0.1"}) != nil {
		logrus.Fatal("SetTrustedProxies failed")
	}

	routes.SetupRouter(router)

	logrus.Info("Starting AnonChat API Service")

	httpServer := &http.Server{
		Addr:    ":8085",
		Handler: router,
	}

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatal("Failed to start HTTP Service: AnonChat API ", err)
	}

	logrus.Info("HTTP Service: AnonChat API started")
}
