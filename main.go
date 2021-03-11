package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mtfelian/gjg-test-task/api"
	"github.com/mtfelian/gjg-test-task/config"
	"github.com/mtfelian/gjg-test-task/service"
	"github.com/sirupsen/logrus"
)

// RegisterHTTPAPIHandlers registers HTTP API handlers
func RegisterHTTPAPIHandlers(router *echo.Echo) {
	router.Use(middleware.Recover())
	router.Use(middleware.CORS())
	router.Use(middleware.Gzip())
	//router.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) { fmt.Printf("@@: %s\n", resBody) }))

	router.POST("/submit", api.SubmitLevel)
}

func main() {
	conf, err := config.Parse()
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.SetLevel(logrus.DebugLevel)

	if err = service.NewWithPostgresClient(conf); err != nil {
		logrus.Fatal(err)
	}

	s := service.Get()
	RegisterHTTPAPIHandlers(s.HTTPServer)
	if err = s.HTTPServer.Start(fmt.Sprintf(":%d", s.Conf.GetInt(config.Port))); err != nil {
		s.Logger.Fatalf("HTTP server error: %v", err)
	}

}
