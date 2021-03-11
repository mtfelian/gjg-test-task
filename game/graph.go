package game

import (
	"errors"
	"fmt"
)

// JI is j and i coordinate pair
type JI struct{ J, I int }

// Vertex with coordinates
type Vertex struct {
	Idx      JI
	Vertices map[JI]*Vertex

	// task-specific
	Value       byte
	RemainingHP int
	BackTrace   *Vertex
}

// NewVertex returns a pointer to a new vertex
func NewVertex(key JI, value byte) *Vertex {
	return &Vertex{
		Idx:      key,
		Value:    value,
		Vertices: map[JI]*Vertex{},
	}
}

// Graph with vertices
type Graph struct {
	Vertices map[JI]*Vertex

	// task-specific
	StartingHP int
}

// NewGraph returns a pointer to a new graph
func NewGraph(HP int) *Graph {
	return &Graph{
		Vertices:   map[JI]*Vertex{},
		StartingHP: HP,
	}
}

// AddVertex creates a new vertex with the given key and value and adds it to the graph
func (g *Graph) AddVertex(key JI, value byte) {
	v := NewVertex(key, value)
	g.Vertices[key] = v
}

// The AddEdge method adds an edge between two vertices in the graph
func (g *Graph) AddEdge(k1, k2 JI) error {
	v1 := g.Vertices[k1]
	v2 := g.Vertices[k2]

	if v1 == nil || v2 == nil {
		return errors.New("not all vertices exist")
	}

	if _, ok := v1.Vertices[v2.Idx]; ok {
		return nil
	}

	v1.Vertices[v2.Idx] = v2
	if v1.Idx != v2.Idx {
		v2.Vertices[v1.Idx] = v1
	}

	g.Vertices[v1.Idx] = v1
	g.Vertices[v2.Idx] = v2
	return nil
}

type node struct {
	v    *Vertex
	next *node
}

type queue struct {
	head *node
	tail *node
}

// enqueue adds a new node to the tail of the queue
func (q *queue) enqueue(v *Vertex) {
	n := &node{v: v}

	if q.tail == nil {
		q.head = n
		q.tail = n
		return
	}

	q.tail.next = n
	q.tail = n
}

// dequeue removes the head from the queue and returns it
func (q *queue) dequeue() *Vertex {
	n := q.head
	if n == nil {
		return nil
	}

	q.head = q.head.next

	if q.head == nil {
		q.tail = nil
	}

	return n.v
}

// BreadthFirstSearch performs breadth-first search for shortest path from startVertex to endVertex in the graph g.
// visitFunc callback is being invoked on each vertex visit
func BreadthFirstSearch(g *Graph, startVertex, endVertex *Vertex, visitFunc func(JI)) (bt []*Vertex) {
	vertexQueue := &queue{}
	visitedVertices := map[JI]bool{}

	currentVertex := startVertex
	for currentVertex.Idx != endVertex.Idx {
		fmt.Println("next iteration, cv=", currentVertex.Idx)
		visitFunc(currentVertex.Idx)
		if currentVertex.BackTrace == nil { // starting point
			fmt.Println("(no bt) set rem hp of", currentVertex.Idx, "to", g.StartingHP)
			currentVertex.RemainingHP = g.StartingHP
		} else { // currentVertex.BackTrace != nil
			switch currentVertex.Value {
			case CellPit:
				fmt.Println("(bt) set rem hp of", currentVertex.Idx, "to", currentVertex.BackTrace.RemainingHP-1)
				currentVertex.RemainingHP = currentVertex.BackTrace.RemainingHP - 1
			case CellArrow:
				fmt.Println("(bt) set rem hp of", currentVertex.Idx, "to", currentVertex.BackTrace.RemainingHP-2)
				currentVertex.RemainingHP = currentVertex.BackTrace.RemainingHP - 2
			default:
				fmt.Println("(bt) set rem hp of", currentVertex.Idx, "to", currentVertex.BackTrace.RemainingHP)
				currentVertex.RemainingHP = currentVertex.BackTrace.RemainingHP
			}
		}

		fmt.Println("rem hp of vertex", currentVertex.Idx, "is", currentVertex.RemainingHP)
		if currentVertex.RemainingHP > 0 { // else can't go further
			fmt.Println("marking visited", currentVertex.Idx)
			visitedVertices[currentVertex.Idx] = true

			for _, v := range currentVertex.Vertices {
				if !visitedVertices[v.Idx] {
					fmt.Println("enqueuing v", v.Idx)
					v.BackTrace = currentVertex
					vertexQueue.enqueue(v)
				}
			}
		}

		prevVertex := currentVertex
		currentVertex = vertexQueue.dequeue()
		if currentVertex == nil {
			currentVertex = prevVertex // edge case when last vertex of traversal is also endVertex
			break
		}
	}

	fmt.Println("break", currentVertex, currentVertex.Idx, endVertex.Idx)
	if currentVertex != nil && currentVertex.Idx == endVertex.Idx {
		fmt.Println("append", currentVertex.Idx)
		bt = append(bt, currentVertex)
		for currentVertex.BackTrace != nil {
			currentVertex = currentVertex.BackTrace
			fmt.Println("append", currentVertex.Idx)
			bt = append(bt, currentVertex)
		}
	}
	return
}
