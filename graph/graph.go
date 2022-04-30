package graph

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
