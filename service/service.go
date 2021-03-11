package service

import (
	"github.com/labstack/echo/v4"
	"github.com/mtfelian/gjg-test-task/config"
	"github.com/mtfelian/gjg-test-task/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Service represents service singleton
type Service struct {
	Storage    storage.Keeper
	Logger     *logrus.Logger
	Conf       *viper.Viper
	HTTPServer *echo.Echo
}

var singleton *Service

// Get provides access to a service components
func Get() *Service { return singleton }

// newService creates new service object
func newService(conf *viper.Viper, keeper storage.Keeper) error {
	singleton = &Service{
		Storage:    keeper,
		Logger:     logrus.New(),
		Conf:       conf,
		HTTPServer: echo.New(),
	}
	return nil
}

// NewWithPostgresClient creates new service object with PostgreSQL-based keeper
func NewWithPostgresClient(conf *viper.Viper) (err error) {
	if err = newService(conf, nil); err != nil {
		return
	}
	keeper := storage.NewPostgresKeeper(storage.PostgresKeeperSettings{
		Conn: storage.PostgresConnection{
			User:     conf.GetString(config.DBLogin),
			Password: conf.GetString(config.DBPassword),
			Host:     conf.GetString(config.DBHost),
			Port:     conf.GetString(config.DBPort),
			Name:     conf.GetString(config.DBName),
			Schema:   conf.GetString(config.DBSchema),
		},
		Logger: singleton.Logger,
	})
	singleton.Storage = keeper
	return
}
