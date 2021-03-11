package storage

import (
	"github.com/go-openapi/strfmt"
	"github.com/mtfelian/gjg-test-task/storage/model"
)

// Keeper abstracts data storage
type Keeper interface {
	ApplyMigrations(path, migrateCommand string) error
	Close() error
	AddLevel(levelData model.Level) (strfmt.UUID, error)
	RemoveAll() error
	GetLevels(p model.GetLevelsParams) (levels []model.Level, err error)
}
