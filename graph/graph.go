package graph

import (
	"fmt"
	"sort"
)

type Vertex struct {
	Id         uint64
	Label      string
	Properties interface{}
}

type Edge struct {
	Id         uint64
	Label      string
	Properties interface{}
	From       uint64
	To         uint64
}

type Graph struct {
	Vertices  map[uint64]Vertex
	Edges     map[uint64]Edge
	Adjacency map[uint64][]uint64
}

func New() *Graph {
	return &Graph{
		Vertices:  map[uint64]Vertex{},
		Edges:     map[uint64]Edge{},
		Adjacency: map[uint64][]uint64{},
	}
}

func (g *Graph) AddVertex(vertex Vertex) {
	g.Vertices[vertex.Id] = vertex
}

func (g *Graph) AddEdge(v1, v2 uint64) {
	if !g.IsVertex(v1) {
		err := fmt.Sprintf("Vertex does not exist: %v", v1)
		panic(err)
	}
	if !g.IsVertex(v2) {
		err := fmt.Sprintf("Vertex does not exist: %v", v2)
		panic(err)
	}
	g.addEdge(v1, v2)
	g.addEdge(v2, v1)
}

func (g *Graph) addEdge(v1, v2 uint64) {
	g.Adjacency[v1] = append(g.Adjacency[v1], v2)

	edgeId := uint64(len(g.Edges) + 1)
	g.Edges[edgeId] = Edge{Id: edgeId, Label: "link", From: v1, To: v2}
}

func (g *Graph) IsEdge(v1, v2 uint64) bool {
	for _, v := range g.Adjacency[v1] {
		if v == v2 {
			return true
		}
	}
	return false
}

func (g *Graph) IsVertex(id uint64) bool {
	_, exists := g.Vertices[id]
	return exists
}

type WalkFunc func(vertex Vertex, depth int) bool

func (g *Graph) Walk(callback WalkFunc) {
	visited := map[uint64]bool{}

	sortedVertices := []uint64{}
	for vertex := range g.Vertices {
		sortedVertices = append(sortedVertices, vertex)
	}
	sort.Slice(sortedVertices, func(a, b int) bool { return sortedVertices[a] < sortedVertices[b] })

	for _, id := range sortedVertices {
		vertex := g.Vertices[id]
		g.walk(vertex, 0, visited, callback)
	}
}

func (g *Graph) walk(vertex Vertex, depth int, visited map[uint64]bool, callback WalkFunc) {
	if visited[vertex.Id] {
		return
	}

	callback(vertex, depth)

	visited[vertex.Id] = true

	for _, childId := range g.Adjacency[vertex.Id] {
		child := g.Vertices[childId]
		if visited[childId] {
			continue
		}

		g.walk(child, depth+1, visited, callback)
	}
}
