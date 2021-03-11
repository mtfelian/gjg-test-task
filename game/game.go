package game

import (
	"fmt"

	"github.com/mtfelian/gjg-test-task/service"
	"github.com/mtfelian/gjg-test-task/storage/model"
)

// cell values
const (
	CellOpen = iota
	CellWall
	CellPit
	CellArrow
	CellPlayer
)

// dimension limits
const (
	MaxDim     = 100
	MinDim     = 2
	startingHP = 4
)

// Position represents game field
type Position struct {
	Maze [][]byte
	X, Y int // to be set after validation
}

// ToStorage converts position p to storage layer format
func (p Position) ToStorage() model.Level {
	level := model.Level{
		X:    p.X,
		Y:    p.Y,
		Maze: make([]byte, p.X*p.Y),
	}
	for i, row := range p.Maze {
		for j, cell := range row {
			level.Maze[i*p.X+j] = cell
		}
	}
	return level
}

// Validate the field
func (p Position) Validate() *Error {
	lenMaze := len(p.Maze)
	switch {
	case lenMaze > 100:
		return &Error{
			Code:    service.ErrValidationFieldIsTooLarge,
			Message: fmt.Sprintf("Position contains %d rows, max is %d", lenMaze, MaxDim),
			Params:  []interface{}{lenMaze, MaxDim},
		}
	case lenMaze < 2:
		return &Error{
			Code:    service.ErrValidationFieldIsTooSmall,
			Message: fmt.Sprintf("Position contains %d rows, min is %d", lenMaze, MinDim),
			Params:  []interface{}{lenMaze, MinDim},
		}
	}

	row0Length := len(p.Maze[0])
	for i, row := range p.Maze {
		lenRow := len(row)
		switch {
		case lenRow != row0Length:
			return &Error{
				Code:    service.ErrValidationFieldIsNotRectangular,
				Message: fmt.Sprintf("Row %d contains %d columns, while row 0 contains %d", i, lenRow, row0Length),
				Params:  []interface{}{i, lenRow, row0Length},
			}
		case lenRow > 100:
			return &Error{
				Code:    service.ErrValidationFieldIsTooLarge,
				Message: fmt.Sprintf("Row %d contains %d columns, max is %d", i, lenRow, MaxDim),
				Params:  []interface{}{i, lenRow, MaxDim},
			}
		case lenRow < 2:
			return &Error{
				Code:    service.ErrValidationFieldIsTooSmall,
				Message: fmt.Sprintf("Row %d contains %d columns, min is %d", i, lenRow, MinDim),
				Params:  []interface{}{i, lenRow, MinDim},
			}
		}

		for j, cell := range row {
			if cell < CellOpen || cell > CellPlayer {
				return &Error{
					Code:    service.ErrValidationFieldHasInvalidData,
					Message: fmt.Sprintf("Cell (%d,%d) contains invalid value %d", i, j, cell),
					Params:  []interface{}{i, lenRow, MinDim},
				}
			}
		}
	}
	return nil
}

// ToGraph converts a position to a graph
func (p Position) ToGraph() (*Graph, error) {
	res := NewGraph(startingHP)
	for i, row := range p.Maze {
		for j, cell := range row {
			fmt.Println("adding vertex", j, i)
			res.AddVertex(JI{j, i}, cell)
		}
	}
	for i, row := range p.Maze {
		for j := range row {
			if p.Maze[i][j] == CellWall {
				continue
			}
			if j > 0 && p.Maze[i][j-1] != CellWall {
				fmt.Println("adding edge", j, i, j-1, i)
				if err := res.AddEdge(JI{j, i}, JI{j - 1, i}); err != nil {
					return nil, err
				}
			}
			if j < len(row)-1 && p.Maze[i][j+1] != CellWall {
				fmt.Println("adding edge", j, i, j+1, i)
				if err := res.AddEdge(JI{j, i}, JI{j + 1, i}); err != nil {
					return nil, err
				}
			}
			if i > 0 && p.Maze[i-1][j] != CellWall {
				fmt.Println("adding edge", j, i, j, i-1)
				if err := res.AddEdge(JI{j, i}, JI{j, i - 1}); err != nil {
					return nil, err
				}
			}
			if i < len(p.Maze)-1 && p.Maze[i+1][j] != CellWall {
				fmt.Println("adding edge", j, i, j, i+1)
				if err := res.AddEdge(JI{j, i}, JI{j, i + 1}); err != nil {
					return nil, err
				}
			}
		}
	}
	return res, nil
}
