package model

import (
	uuid "github.com/satori/go.uuid"
)

// GetLevelsParams represents parameters for requesting levels
type GetLevelsParams struct {
}

// Level represents price level
type Level struct {
	tableName struct{} `pg:"levels"`

	ID   uuid.UUID `pg:"id,notnull,type:uuid,pk"`
	X    int       `pg:"x,notnull"`
	Y    int       `pg:"y,notnull"`
	Maze []byte    `pg:"maze,notnull"`
}
