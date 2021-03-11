package config

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// misc constants
const (
	FileName = "config.json"
)

// executable flag names
const (
	Port = "port"

	DBHost     = "db_host"
	DBPort     = "db_port"
	DBSchema   = "db_schema"
	DBName     = "db_name"
	DBLogin    = "db_login"
	DBPassword = "db_password"
	DBMigrate  = "db_migrate"

	LogLevel = "loglevel"
)

// errors
var (
	ErrorSpecifyPort = errors.New("should specify --port command line flag")
)

func checkRequired() error {
	if viper.GetInt(Port) == 0 {
		return ErrorSpecifyPort
	}
	return nil
}

// Parse the configuration data from different sources
func Parse() (*viper.Viper, error) {
	if err := parseFlags(); err != nil {
		log.Fatalln("Error parsing flags:", err)
	}
	if err := parseConfigFile(); err != nil {
		log.Fatalln("Error parsing config file:", err)
	}
	if err := checkRequired(); err != nil {
		log.Fatalln(err)
	}
	return viper.GetViper(), nil
}

// parseFlags from command line
func parseFlags() error {
	pflag.UintVar(&params.Port, Port, 0, "application HTTP port")

	pflag.StringVar(&params.LogLevel, LogLevel, logrus.InfoLevel.String(), "logging level")

	pflag.StringVar(&params.DBHost, DBHost, "", "DB host")
	pflag.StringVar(&params.DBPort, DBPort, "5432", "DB port")
	pflag.StringVar(&params.DBSchema, DBSchema, "public", "DB schema")
	pflag.StringVar(&params.DBName, DBName, "gjg_games", "DB name")
	pflag.StringVar(&params.DBLogin, DBLogin, "postgres", "DB login")
	pflag.StringVar(&params.DBPassword, DBPassword, "", "DB password")
	pflag.StringVar(&params.DBMigrate, DBMigrate, "", "DB migration commands")

	pflag.Parse()
	return viper.BindPFlags(pflag.CommandLine)
}

// parseConfigFile
func parseConfigFile(paths ...string) error {
	viper.AddConfigPath(".")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigName(viper.GetString(FileName))
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return fmt.Errorf("failed to load config file: %v", err)
	}
	return nil
}

// Options describes command line options
type Options struct {
	sync.Mutex

	// Port to listen
	Port uint

	// DBHost is DB host
	DBHost string
	// DBPort is DB port
	DBPort string
	// DBSchema is DB scheme
	DBSchema string
	// DBName is DB name
	DBName string
	// DBLogin is DB user login
	DBLogin string
	// DBPassword is DB user password
	DBPassword string
	// DBMigrate is DB migration tool commands
	DBMigrate string

	// LogLevel is a logging level
	LogLevel string
}

// params is an application command line parameters
var params Options
