package graph

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/msp301/zb/parser"
)

type Vertex struct {
	Id         uint64
	Label      string
	Properties map[string]interface{}
}

type Edge struct {
	Id         uint64
	Label      string
	Properties map[string]interface{}
	From       uint64
	To         uint64
}

type Graph struct {
	Vertices  map[uint64]Vertex
	Edges     map[uint64]Edge
	Adjacency map[uint64]map[uint64]int
}

func New() *Graph {
	return &Graph{
		Vertices:  map[uint64]Vertex{},
		Edges:     map[uint64]Edge{},
		Adjacency: map[uint64]map[uint64]int{},
	}
}

func (g *Graph) AddVertex(vertex Vertex) {
	if vertex.Id == 0 {
		vertex.Id = uint64(len(g.Vertices) + 1)
	}
	g.Vertices[vertex.Id] = vertex
}

func (g *Graph) AddEdge(edge Edge) error {
	if !g.IsVertex(edge.From) {
		return fmt.Errorf("vertex does not exist: %v", edge.From)
	}
	if !g.IsVertex(edge.To) {
		return fmt.Errorf("vertex does not exist: %v", edge.To)
	}
	g.addEdge(edge)

	reverse := edge
	reverse.From = edge.To
	reverse.To = edge.From
	g.addEdge(reverse)

	return nil
}

func (g *Graph) addEdge(edge Edge) {
	if _, ok := g.Adjacency[edge.From]; !ok {
		g.Adjacency[edge.From] = map[uint64]int{}
	}

	g.Adjacency[edge.From][edge.To] = len(g.Adjacency[edge.From]) + 1

	edgeId := uint64(len(g.Edges) + 1)
	edge.Id = edgeId
	g.Edges[edgeId] = edge
}

func (g *Graph) IsEdge(v1, v2 uint64) bool {
	_, ok := g.Adjacency[v1][v2]
	return ok
}

func (g *Graph) IsVertex(id uint64) bool {
	_, exists := g.Vertices[id]
	return exists
}

type WalkFunc func(vertex Vertex, depth int) bool

func (g *Graph) Walk(callback WalkFunc, maxDepth int) {
	g.orderedWalk(callback, maxDepth, func(a, b uint64) bool {
		return a < b
	})
}

func (g *Graph) WalkBackwards(callback WalkFunc, maxDepth int) {
	g.orderedWalk(callback, maxDepth, func(a, b uint64) bool {
		return a > b
	})
}

func (g *Graph) WalkGit(callback WalkFunc, maxDepth int) {
	g.orderedWalk(callback, maxDepth, func(a, b uint64) bool {
		refA := a
		refB := b

		switch vertexA := g.Vertices[a].Properties["Value"].(type) {
		case parser.Note:
			refA = gitNoteId(vertexA)
		}

		switch vertexB := g.Vertices[b].Properties["Value"].(type) {
		case parser.Note:
			refB = gitNoteId(vertexB)
		}

		return refA > refB
	})
}

func gitNoteId(note parser.Note) uint64 {
	noteDir := filepath.Dir(note.File)
	ref, err := git.PlainOpen(noteDir)
	if err != nil {
		fmt.Printf("Error opening git repo: %s\n", err)
		return note.Id
	}

	cIter, err := ref.Log(&git.LogOptions{FileName: &note.File})
	if err != nil {
		fmt.Printf("Error getting log: %s\n", err)
		return note.Id
	}

	commit, err := cIter.Next()
	if err != nil {
		return note.Id
	}

	commitTime := commit.Committer.When
	commitTimeString := fmt.Sprintf("%d%d%d%d%d", commitTime.Year(), commitTime.Month(), commitTime.Day(), commitTime.Hour(), commitTime.Minute())
	commitTimeInt, err := strconv.ParseUint(commitTimeString, 0, 64)
	if err != nil {
		fmt.Printf("Error parsing commit time: %s\n", err)
		return note.Id
	}

	return commitTimeInt
}

type WalkSort func(a, b uint64) bool

func (g *Graph) orderedWalk(callback WalkFunc, maxDepth int, sortVertices WalkSort) {
	visited := map[uint64]bool{}

	var sortedVertices []uint64
	for vertex := range g.Vertices {
		sortedVertices = append(sortedVertices, vertex)
	}
	sort.SliceStable(sortedVertices, func(a, b int) bool { return sortVertices(sortedVertices[a], sortedVertices[b]) })

	for _, id := range sortedVertices {
		vertex := g.Vertices[id]
		g.walk(vertex, 0, maxDepth, visited, callback)
	}
}

func (g *Graph) walk(vertex Vertex, depth int, maxDepth int, visited map[uint64]bool, callback WalkFunc) {
	if visited[vertex.Id] {
		return
	}

	if depth > maxDepth && maxDepth != -1 {
		return
	}

	callback(vertex, depth)

	visited[vertex.Id] = true

	for childId := range g.Adjacency[vertex.Id] {
		child := g.Vertices[childId]
		if visited[childId] {
			continue
		}

		g.walk(child, depth+1, maxDepth, visited, callback)
	}
}
