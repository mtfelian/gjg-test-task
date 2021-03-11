package api

import (
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/mtfelian/gjg-test-task/game"
	"github.com/mtfelian/gjg-test-task/service"
)

// SubmitLevelParams represents parameters for SubmitLevel handler
type SubmitLevelParams struct {
	Maze [][]byte `json:"maze"`
}

// SubmitLevelResponse represents response for SubmitLevel handler
type SubmitLevelResponse struct {
	LevelID strfmt.UUID `json:"id"`
}

// Model converts API model to storage layer model
func (p SubmitLevelParams) ToPosition() (*game.Position, *game.Error) {
	position := &game.Position{Maze: p.Maze}
	if err := position.Validate(); err != nil {
		return nil, err
	}
	position.X = len(p.Maze[0])
	position.Y = len(p.Maze)
	return position, nil
}

// SubmitLevel is an API handler to submit level
func SubmitLevel(c echo.Context) error {
	var code int
	p := new(SubmitLevelParams)
	if err := c.Bind(p); err != nil {
		code = http.StatusUnprocessableEntity
		return c.JSON(code, game.Error{Code: service.ErrValidationRequest, Message: err.Error()})
	}

	position, Err := p.ToPosition()
	if Err != nil {
		code = http.StatusBadRequest
		return c.JSON(code, *Err)
	}

	s := service.Get()
	newLevelID, err := s.Storage.AddLevel(position.ToStorage())
	if err != nil {
		code = http.StatusInternalServerError
		return c.JSON(code, game.Error{Code: service.ErrStorageFailed, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, SubmitLevelResponse{LevelID: newLevelID})
}
