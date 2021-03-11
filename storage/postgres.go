package storage

import (
	"context"
	"os"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/mtfelian/gjg-test-task/config"
	"github.com/mtfelian/gjg-test-task/storage/model"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// PostgresKeeper provides access to PostgreSQL storage
type PostgresKeeper struct {
	pdb    *pg.DB
	Logger *logrus.Logger
}

// PostgresKeeperSettings contains PostgresKeeper settings
type PostgresKeeperSettings struct {
	Conn   PostgresConnection
	Logger *logrus.Logger
}

// PostgresConnection contains PostgreSQL connection params
type PostgresConnection struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Schema   string `json:"schema"`
}

// SetIntoViper sets config values into Viper configuration
func (pc PostgresConnection) SetIntoViper(conf *viper.Viper) {
	conf.Set(config.DBLogin, pc.User)
	conf.Set(config.DBPassword, pc.Password)
	conf.Set(config.DBHost, pc.Host)
	conf.Set(config.DBPort, pc.Port)
	conf.Set(config.DBName, pc.Name)
	conf.Set(config.DBSchema, pc.Schema)
}

type dbLogger struct{ logger *logrus.Logger }

// BeforeQuery hook
func (d dbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

// AfterQuery hook
func (d dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	d.logger.Debugln(q.FormattedQuery())
	return nil
}

func (keeper *PostgresKeeper) initMigrations() error {
	if _, _, err := migrations.Run(keeper.pdb, "init"); err != nil {
		return err
	}
	migrateCommand, passedMigrationCommand := viper.GetString(config.DBMigrate), true
	if migrateCommand == "" {
		migrateCommand, passedMigrationCommand = "up", false
	}
	if err := keeper.ApplyMigrations("migrations", migrateCommand); err != nil {
		return err
	}
	if passedMigrationCommand {
		keeper.Logger.Infof("Finished. Migration tool successfully ran command: %s", migrateCommand)
		os.Exit(0)
	}
	return nil
}

// ApplyMigrations with the given migrateCommand. Migrations will be searched according to path given.
func (keeper *PostgresKeeper) ApplyMigrations(path, migrateCommand string) error {
	if err := migrations.DefaultCollection.DiscoverSQLMigrations(path); err != nil {
		return err
	}
	oldVersion, newVersion, err := migrations.Run(keeper.pdb, strings.Split(migrateCommand, ",")...)
	if err != nil {
		return err
	}

	if newVersion != oldVersion {
		keeper.Logger.Infof("DB was migrated from version %d to %d", oldVersion, newVersion)
	} else {
		keeper.Logger.Infof("DB migrations version is %d", oldVersion)
	}
	return nil
}

// NewPostgresKeeper returns a pointer to a new keeper based on PostgreSQL storage
// valuesKeepMode parameter determines where to keep data values (DB or files)
func NewPostgresKeeper(pks PostgresKeeperSettings) *PostgresKeeper {
	db := pg.Connect(&pg.Options{
		Addr:     pks.Conn.Host + ":" + pks.Conn.Port,
		User:     pks.Conn.User,
		Password: pks.Conn.Password,
		Database: pks.Conn.Name,
	})
	db.AddQueryHook(dbLogger{logger: pks.Logger})

	// set schema and table name to be singular
	//orm.SetTableNameInflector(func(s string) string { return pks.Conn.Schema + "." + s })

	keeper := &PostgresKeeper{pdb: db, Logger: pks.Logger}

	if err := keeper.initMigrations(); err != nil {
		panic(err)
	}
	return keeper
}

// Close the DB connection
func (keeper *PostgresKeeper) Close() error { return keeper.pdb.Close() }

// AddLevel to the storage
func (keeper *PostgresKeeper) AddLevel(level model.Level) (strfmt.UUID, error) {
	level.ID = uuid.NewV4()
	_, err := keeper.pdb.Model(&level).Insert()
	return strfmt.UUID(level.ID.String()), err
}

// AddLevel to the storage
func (keeper *PostgresKeeper) GetLevels(p model.GetLevelsParams) (levels []model.Level, err error) {
	err = keeper.modifyLevelsQuery(keeper.pdb.Model(&levels), p).Select()
	return
}

// modifyLevelsQuery with given params p
func (keeper *PostgresKeeper) modifyLevelsQuery(query *orm.Query, p model.GetLevelsParams) *orm.Query {
	// construct query by p
	return query
}

// RemoveLevels removes levels according to the given params
func (keeper *PostgresKeeper) RemoveLevels(p model.GetLevelsParams) (err error) {
	_, err = keeper.modifyLevelsQuery(keeper.pdb.Model((*model.Level)(nil)).Where("TRUE"), p).Delete()
	return
}

// RemoveAll entities
func (keeper *PostgresKeeper) RemoveAll() (err error) {
	for _, f := range []func() error{
		func() error { return keeper.RemoveLevels(model.GetLevelsParams{}) },
		// add more removal funcs
	} {
		if err = f(); err != nil {
			return
		}
	}
	return
}
