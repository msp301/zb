package graph

import "sort"

type TraversalSource struct {
	channel  chan Vertex
	graph    *Graph
	what     string
	label    map[string]bool
	position *Vertex
	sorted   []uint64
	visited  map[uint64]bool
}

func Traversal(g *Graph) *TraversalSource {
	return &TraversalSource{
		graph:   g,
		visited: map[uint64]bool{},
	}
}

func (t *TraversalSource) V() *TraversalSource {
	t.what = "vertex"
	return t
}

func (t *TraversalSource) HasLabel(label string) *TraversalSource {
	if t.label == nil {
		t.label = map[string]bool{}
	}
	t.label[label] = true
	return t
}

func (t *TraversalSource) Next() *Vertex {
	if t.channel == nil {
		t.Iterate()
	}
	vertex, ok := <-t.channel
	if ok {
		return &vertex
	}
	return nil
}

func (t *TraversalSource) Iterate() <-chan Vertex {
	if t.channel == nil {
		t.channel = make(chan Vertex)
	}

	go func() {
		if t.position == nil {
			vertex, ok := t.graph.Vertices[t.sortedVertices()[0]]
			if ok {
				t.position = &vertex
			}
		}

		for _, id := range t.sortedVertices() {
			vertex := t.graph.Vertices[id]
			t.walk(t.channel, vertex, 0)
		}

		close(t.channel)
	}()

	return t.channel
}

func (t *TraversalSource) walk(channel chan Vertex, vertex Vertex, depth int) {
	if t.visited[vertex.Id] {
		return
	}
	if t.label != nil && t.label[vertex.Label] == false {
		return
	}

	channel <- vertex
	t.visited[vertex.Id] = true

	for _, childId := range t.graph.Adjacency[vertex.Id] {
		child := t.graph.Vertices[childId]
		if t.visited[childId] {
			continue
		}

		t.walk(channel, child, depth+1)
	}
}

func (t *TraversalSource) sortedVertices() []uint64 {
	if t.sorted != nil {
		return t.sorted
	}
	for vertex := range t.graph.Vertices {
		t.sorted = append(t.sorted, vertex)
	}
	sort.Slice(t.sorted, func(a, b int) bool { return t.sorted[a] < t.sorted[b] })
	return t.sorted
}
