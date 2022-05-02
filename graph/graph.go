package graph

import (
	"sort"
)

type Graph struct {
	Vertices map[uint64]bool
	Edges    map[uint64][]uint64
}

func New() *Graph {
	return &Graph{
		Vertices: map[uint64]bool{},
		Edges:    map[uint64][]uint64{},
	}
}

func (g *Graph) AddVertex(id uint64) {
	_, exists := g.Vertices[id]
	if exists {
		return
	}
	g.Vertices[id] = true
}

func (g *Graph) AddEdge(v1, v2 uint64) {
	if g.IsEdge(v1, v2) {
		return
	}
	g.Edges[v1] = append(g.Edges[v1], v2)
	g.Edges[v2] = append(g.Edges[v2], v1)
}

func (g *Graph) IsEdge(v1, v2 uint64) bool {
	for _, v := range g.Edges[v1] {
		if v == v2 {
			return true
		}
	}
	return false
}

type WalkFunc func(id uint64, depth int) bool

func (g *Graph) Walk(callback WalkFunc) {
	visited := map[uint64]bool{}

	sortedVertices := []uint64{}
	for vertex := range g.Vertices {
		sortedVertices = append(sortedVertices, vertex)
	}
	sort.Slice(sortedVertices, func(a, b int) bool { return sortedVertices[a] < sortedVertices[b] })

	for _, vertex := range sortedVertices {
		g.walk(vertex, 0, visited, callback)
	}
}

func (g *Graph) walk(vertex uint64, depth int, visited map[uint64]bool, callback WalkFunc) {
	if visited[vertex] {
		return
	}

	callback(vertex, depth)

	visited[vertex] = true

	for _, child := range g.Edges[vertex] {
		if visited[child] {
			continue
		}

		g.walk(child, depth+1, visited, callback)
	}
}
