package notebook

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/msp301/graph"
)

type tagQueryEdge struct{ note, tag uint64 }

type tagQuerySubgraph struct {
	vertices map[uint64]graph.Vertex
	edges    map[tagQueryEdge]struct{}
	matches  map[uint64][]matchedTag
	returned map[uint64]struct{}
}

// VisualizeTagQuery returns a DOT-format graph of the subgraph involved in the
// given tag query. Tags matching a search term and notes returned by
// SearchByTags are filled.
func (book *Notebook) VisualizeTagQuery(searchTags ...string) string {
	return book.tagQuerySubgraph(searchTags...).dot()
}

func (book *Notebook) tagQuerySubgraph(searchTags ...string) tagQuerySubgraph {
	subgraph := tagQuerySubgraph{
		vertices: map[uint64]graph.Vertex{},
		edges:    map[tagQueryEdge]struct{}{},
		matches:  map[uint64][]matchedTag{},
		returned: map[uint64]struct{}{},
	}

	matchedTags := book.MatchedTags(searchTags...)
	for _, matched := range matchedTags {
		tagId := matched.Vertex.Id
		subgraph.vertices[tagId] = matched.Vertex
		subgraph.matches[tagId] = append(subgraph.matches[tagId], matched)

		for noteId := range book.Notes.Adjacency[tagId] {
			note := book.Notes.Vertices[noteId]
			if note.Label != "note" {
				continue
			}
			subgraph.vertices[noteId] = note
			for _, tag := range book.tagsOfNotes([]graph.Vertex{note}) {
				subgraph.vertices[tag.Id] = tag
				subgraph.edges[tagQueryEdge{noteId, tag.Id}] = struct{}{}
			}
		}
	}

	if everyTermMatchesATag(searchTags, matchedTags) {
		for _, note := range book.TagIntersection(matchedTags) {
			subgraph.returned[note.Id] = struct{}{}
		}
	}

	return subgraph
}

func (subgraph tagQuerySubgraph) dot() string {
	var sb strings.Builder
	sb.WriteString("digraph {\n")

	for _, id := range subgraph.sortedVertexIds() {
		fmt.Fprintf(&sb, "  %d [%s];\n", id, subgraph.vertexAttributes(subgraph.vertices[id]))
	}
	for _, edge := range subgraph.sortedEdges() {
		fmt.Fprintf(&sb, "  %d -> %d;\n", edge.note, edge.tag)
	}

	sb.WriteString("}\n")
	return sb.String()
}

func (subgraph tagQuerySubgraph) vertexAttributes(vertex graph.Vertex) string {
	if vertex.Label == "note" {
		attributes := fmt.Sprintf("label=%s shape=box", dotString(fmt.Sprintf("note:%d", vertex.Id)))
		if _, ok := subgraph.returned[vertex.Id]; ok {
			attributes += " style=filled fillcolor=lightgreen"
		}
		return attributes
	}

	lines := []string{fmt.Sprint(vertex.Value)}
	matches := slices.Clone(subgraph.matches[vertex.Id])
	slices.SortFunc(matches, func(a, b matchedTag) int { return cmp.Compare(a.Term, b.Term) })
	for _, matched := range matches {
		lines = append(lines, fmt.Sprintf("term: %s, distance: %d", matched.Term, matched.Distance))
	}

	attributes := fmt.Sprintf("label=%s shape=ellipse", dotString(lines...))
	if len(matches) > 0 {
		attributes += " style=filled fillcolor=lightblue"
	}
	return attributes
}

func (subgraph tagQuerySubgraph) sortedVertexIds() []uint64 {
	ids := make([]uint64, 0, len(subgraph.vertices))
	for id := range subgraph.vertices {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	return ids
}

func (subgraph tagQuerySubgraph) sortedEdges() []tagQueryEdge {
	edges := make([]tagQueryEdge, 0, len(subgraph.edges))
	for edge := range subgraph.edges {
		edges = append(edges, edge)
	}
	slices.SortFunc(edges, func(a, b tagQueryEdge) int {
		return cmp.Or(cmp.Compare(a.note, b.note), cmp.Compare(a.tag, b.tag))
	})
	return edges
}

func dotString(lines ...string) string {
	escaper := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`)
	escaped := make([]string, len(lines))
	for i, line := range lines {
		escaped[i] = escaper.Replace(line)
	}
	return `"` + strings.Join(escaped, `\n`) + `"`
}
