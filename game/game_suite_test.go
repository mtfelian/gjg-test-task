package game_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/mtfelian/gjg-test-task/game"
	"github.com/mtfelian/gjg-test-task/service"
	"github.com/mtfelian/gjg-test-task/storage/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAll(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = Describe("testing position", func() {
	Describe("position", func() {
		BeforeEach(func() {})
		AfterEach(func() {})

		It("checks that 2D array is being converted to flat array in the proper order in ToStorage method", func() {
			type tcs struct {
				in  game.Position
				out model.Level
			}
			for i, tc := range []tcs{
				{
					in: game.Position{
						X: 3, Y: 4,
						Maze: [][]byte{
							{0, 1, 0},
							{0, 0, 4},
							{0, 1, 0},
							{0, 1, 0},
						},
					},
					out: model.Level{X: 3, Y: 4, Maze: []byte{0, 1, 0, 0, 0, 4, 0, 1, 0, 0, 1, 0}},
				},
				{
					in: game.Position{
						X: 2, Y: 2,
						Maze: [][]byte{
							{0, 1},
							{0, 0},
						},
					},
					out: model.Level{X: 2, Y: 2, Maze: []byte{0, 1, 0, 0}},
				},
			} {
				Expect(tc.in.ToStorage()).To(Equal(tc.out), "case %d", i)
			}
		})

		Context("validation", func() {
			It("checks validation cases", func() {
				type tcs struct {
					in           game.Position
					expectedCode int
				}
				m := game.MaxDim
				r := m/3 + 1 // length of sequence for repeating
				for i, tc := range []tcs{
					{ // normal case
						in: game.Position{
							Maze: [][]byte{
								{0, 1},
								{0, 0},
								{0, 1},
								{0, 4},
							},
						},
						expectedCode: service.ErrOK,
					},
					{ // normal case near to edge (ErrValidationFieldIsTooLarge)
						in: game.Position{
							Maze: [][]byte{
								bytes.Repeat([]byte{0, 1, 0}, r)[:m],
								bytes.Repeat([]byte{0, 0, 4}, r)[:m], // it should be invalid to have many "4"
								bytes.Repeat([]byte{0, 1, 0}, r)[:m],
								bytes.Repeat([]byte{0, 1, 0}, r)[:m],
							},
						},
						expectedCode: service.ErrOK,
					},
					{
						in: game.Position{
							Maze: [][]byte{
								bytes.Repeat([]byte{0, 1, 0}, r)[:m+1],
								bytes.Repeat([]byte{0, 0, 4}, r)[:m+1], // it should be invalid to have many "4"
								bytes.Repeat([]byte{0, 1, 0}, r)[:m+1],
								bytes.Repeat([]byte{0, 1, 0}, r)[:m+1],
							},
						},
						expectedCode: service.ErrValidationFieldIsTooLarge,
					},
					{
						in: game.Position{
							Maze: [][]byte{
								{0},
								{4},
								{0},
								{0},
							},
						},
						expectedCode: service.ErrValidationFieldIsTooSmall,
					},
					{
						in: game.Position{
							Maze: [][]byte{
								{0, 1},
								{0, 0, 3},
								{0, 1},
								{0, 4},
							},
						},
						expectedCode: service.ErrValidationFieldIsNotRectangular,
					},
					{
						in: game.Position{
							Maze: [][]byte{
								{0, 1},
								{0, 5},
								{0, 1},
								{0, 4},
							},
						},
						expectedCode: service.ErrValidationFieldHasInvalidData,
					},
					// additional validation ideas:
					// - validate that there is only one "4"
					// - validate that maze is separated by walls and there is only one "exit"
					// - validate that there is a possible way to exit maze from starting position with survival conds
				} {
					Err := tc.in.Validate()
					if tc.expectedCode == service.ErrOK {
						Expect(Err).To(BeNil(), "case %d", i)
						continue
					}
					Expect(Err).NotTo(BeNil(), "case %d", i)
					Expect(Err.Code).To(Equal(tc.expectedCode), "case %d", i)
				}
			})
		})
		Context("traversal", func() {
			findStart := func(maze [][]byte) (i, j int) {
				for i, row := range maze {
					for j, cell := range row {
						if cell == game.CellPlayer {
							return i, j
						}
					}
				}
				return -1, -1
			}

			findFirstExit := func(maze [][]byte) (i, j int) {
				for i, row := range maze {
					for j, cell := range row {
						if j > 0 && j < len(row)-1 && i > 0 && i < len(maze)-1 {
							continue
						}
						if cell == game.CellOpen {
							return i, j
						}
					}
				}
				return -1, -1
			}

			It("checks DFS when right path has one 1hp trap and one 2hp trap -- right path should be chosen", func() {
				p := game.Position{
					Maze: [][]byte{
						{1, 1, 1, 1, 0, 1, 1, 1},
						{1, 0, 0, 0, 0, 0, 0, 1},
						{1, 0, 1, 1, 1, 3, 1, 1},
						{1, 0, 0, 0, 1, 0, 2, 1},
						{1, 1, 1, 0, 1, 1, 0, 1},
						{1, 0, 0, 0, 1, 0, 0, 1},
						{1, 0, 1, 1, 1, 0, 1, 1},
						{1, 0, 0, 4, 0, 0, 0, 1},
						{1, 1, 1, 1, 1, 1, 1, 1},
					},
					X: 8,
					Y: 9,
				}

				var sMaze string
				for _, row := range p.Maze {
					for _, cell := range row {
						sMaze += fmt.Sprintf("%d", cell)
					}
					sMaze += "\n"
				}
				sMaze = strings.TrimSpace(sMaze)

				startI, startJ := findStart(p.Maze)
				Expect(startI).To(BeNumerically(">=", 0))
				Expect(startJ).To(BeNumerically(">=", 0))

				exitI, exitJ := findFirstExit(p.Maze)
				Expect(exitI).To(BeNumerically(">=", 0))
				Expect(exitJ).To(BeNumerically(">=", 0))

				graph, err := p.ToGraph()
				Expect(err).NotTo(HaveOccurred())
				visitedOrder := []game.JI{}
				visitCB := func(cell game.JI) {
					visitedOrder = append(visitedOrder, cell)
				}
				bt := game.BreadthFirstSearch(graph,
					graph.Vertices[game.JI{startJ, startI}],
					graph.Vertices[game.JI{exitJ, exitI}],
					visitCB)
				fmt.Println(visitedOrder)
				Expect(bt).NotTo(BeNil())

				sMazeRows := strings.Split(sMaze, "\n")
				fmt.Println(sMazeRows)
				for _, v := range bt {
					sMazeRows[v.Idx.I] = sMazeRows[v.Idx.I][:v.Idx.J] + "*" + sMazeRows[v.Idx.I][v.Idx.J+1:]
				}
				fmt.Println(strings.Join(sMazeRows, "\n"))
				Expect(sMazeRows).To(Equal([]string{
					"1111*111",
					"1000**01",
					"10111*11",
					"10001**1",
					"111011*1",
					"10001**1",
					"10111*11",
					"100***01",
					"11111111",
				}))

			})

			It("checks DFS when right path has two 2hp traps -- left path should be chosen", func() {
				p := game.Position{
					Maze: [][]byte{
						{1, 1, 1, 1, 0, 1, 1, 1},
						{1, 0, 0, 0, 0, 0, 0, 1},
						{1, 0, 1, 1, 1, 3, 1, 1},
						{1, 0, 0, 0, 1, 0, 3, 1},
						{1, 1, 1, 0, 1, 1, 0, 1},
						{1, 0, 0, 0, 1, 0, 0, 1},
						{1, 0, 1, 1, 1, 0, 1, 1},
						{1, 0, 0, 4, 0, 0, 0, 1},
						{1, 1, 1, 1, 1, 1, 1, 1},
					},
					X: 8,
					Y: 9,
				}

				var sMaze string
				for _, row := range p.Maze {
					for _, cell := range row {
						sMaze += fmt.Sprintf("%d", cell)
					}
					sMaze += "\n"
				}
				sMaze = strings.TrimSpace(sMaze)

				startI, startJ := findStart(p.Maze)
				Expect(startI).To(BeNumerically(">=", 0))
				Expect(startJ).To(BeNumerically(">=", 0))

				exitI, exitJ := findFirstExit(p.Maze)
				Expect(exitI).To(BeNumerically(">=", 0))
				Expect(exitJ).To(BeNumerically(">=", 0))

				graph, err := p.ToGraph()
				Expect(err).NotTo(HaveOccurred())
				visitedOrder := []game.JI{}
				visitCB := func(cell game.JI) {
					visitedOrder = append(visitedOrder, cell)
				}
				bt := game.BreadthFirstSearch(graph,
					graph.Vertices[game.JI{startJ, startI}],
					graph.Vertices[game.JI{exitJ, exitI}],
					visitCB)
				fmt.Println(visitedOrder)
				Expect(bt).NotTo(BeNil())

				sMazeRows := strings.Split(sMaze, "\n")
				fmt.Println(sMazeRows)
				for _, v := range bt {
					sMazeRows[v.Idx.I] = sMazeRows[v.Idx.I][:v.Idx.J] + "*" + sMazeRows[v.Idx.I][v.Idx.J+1:]
				}
				fmt.Println(strings.Join(sMazeRows, "\n"))
				Expect(sMazeRows).To(Equal([]string{
					"1111*111",
					"1****001",
					"1*111311",
					"1***1031",
					"111*1101",
					"1***1001",
					"1*111011",
					"1***0001",
					"11111111",
				}))

			})
		})
	})
})
