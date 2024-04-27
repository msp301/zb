package notebook

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/parser"
	"github.com/msp301/zb/util"
)

type Notebook struct {
	files   []string
	Filters []FilterFunc
	Invalid map[uint64]parser.Note
	Notes   *graph.Graph
	tags    map[string]uint64
}

type FilterFunc func(note parser.Note) bool

func New() *Notebook {
	return &Notebook{
		Invalid: map[uint64]parser.Note{},
		Notes:   graph.New(),
		tags:    map[string]uint64{},
		files:   []string{},
	}
}

func (book *Notebook) AddFilter(filter FilterFunc) {
	book.Filters = append(book.Filters, filter)
}

func (book *Notebook) Load(dir string) error {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		nameWithoutSuffix := strings.TrimSuffix(d.Name(), ".md")
		if !util.IdRegex.MatchString(nameWithoutSuffix) {
			return nil
		}

		fileId, err := strconv.ParseUint(nameWithoutSuffix, 0, 64)
		if err != nil {
			return nil
		}

		book.Notes.AddVertex(graph.Vertex{Id: fileId})
		book.files = append(book.files, path)

		return nil
	})

	return err
}

func (book *Notebook) Read() []parser.Note {
	var filteredNotes []parser.Note

	for _, path := range book.files {
		fileNotes := parser.Parse(path)

	NOTE:
		for _, note := range fileNotes {
			book.Notes.AddVertex(graph.Vertex{Id: note.Id, Label: "note", Properties: map[string]interface{}{"Value": note}})

			for _, link := range note.Links {
				err := book.Notes.AddEdge(graph.Edge{
					Label: "link",
					From:  note.Id,
					To:    link,
				})
				if err != nil {
					book.Invalid[note.Id] = note
				}
			}

			for _, tag := range note.Tags {
				tagId := book.addTag(tag)
				book.Notes.AddVertex(graph.Vertex{Id: tagId, Label: "tag", Properties: map[string]interface{}{"Value": tag}})
				err := book.Notes.AddEdge(graph.Edge{
					Label: "tag",
					From:  note.Id,
					To:    tagId,
				})
				if err != nil {
					book.Invalid[note.Id] = note
				}
			}

			for _, tag := range note.Tags {
				if !book.Notes.IsVertex(book.tags[tag]) {
					book.Notes.AddVertex(graph.Vertex{Id: book.tags[tag], Label: "tag", Properties: map[string]interface{}{"Value": tag}})
				}

				for _, relatedTag := range note.Tags {
					if relatedTag == tag {
						continue
					}
					if !book.Notes.IsVertex(book.tags[relatedTag]) {
						book.Notes.AddVertex(graph.Vertex{Id: book.tags[relatedTag], Label: "tag", Properties: map[string]interface{}{"Value": tag}})
					}
					err := book.Notes.AddEdge(graph.Edge{
						Label: "tag",
						From:  book.tags[tag],
						To:    book.tags[relatedTag],
					})
					if err != nil {
						book.Invalid[note.Id] = note
					}
				}
			}

			for _, filter := range book.Filters {
				filtered := filter(note)
				if filtered {
					continue NOTE
				}
			}
			filteredNotes = append(filteredNotes, note)
		}
	}

	return filteredNotes
}

func (book *Notebook) IsNote(noteId uint64) bool {
	_, ok := book.Notes.Vertices[noteId]
	return ok
}

func (book *Notebook) SearchRelated(id uint64) []Result {
	var results []Result
	for adjId := range book.Notes.Adjacency[id] {
		vertex := book.Notes.Vertices[adjId]
		contexts := []util.ContextMatch{{Text: "", Line: 0}}
		startLine := 0

		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			startLine = val.Start
			matched, ok := util.Context(val.Content, fmt.Sprint(id))
			if ok {
				contexts = matched
			}
		}

		for _, context := range contexts {
			result := Result{
				Context: context.Text,
				Line:    startLine + context.Line,
				Value:   vertex,
			}
			results = append(results, result)
		}
	}

	return results
}

type Result struct {
	Context string
	Line    int
	Value   interface{}
}

type matchedTag struct {
	Distance int
	Tag      string
	Vertex   graph.Vertex
}

func (book *Notebook) SearchByTags(searchTags ...string) []Result {
	var results []Result

	tagVerticesSlice := book.MatchedTags(searchTags...)
	intersection := book.TagIntersection(tagVerticesSlice)
	for _, vertex := range intersection {
		context := []util.ContextMatch{{Text: "", Line: 0}}
		startLine := 0

		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			startLine = val.Start
			matchedTag := ""
			for _, tagVertex := range tagVerticesSlice {
				if book.Notes.IsEdge(tagVertex.Vertex.Id, vertex.Id) {
					matchedTag = fmt.Sprint(tagVertex.Vertex.Properties["Value"])
					break
				}
			}
			matched, ok := util.Context(val.Content, matchedTag)
			if ok {
				context = matched
			}
		}

		for _, context := range context {
			result := Result{
				Context: context.Text,
				Line:    startLine + context.Line,
				Value:   vertex,
			}
			results = append(results, result)
		}
	}

	return results
}

func (book *Notebook) Count() int {
	var noteCount int
	traversal := graph.Traversal(book.Notes)
	for range traversal.V().HasLabel("note").Iterate() {
		noteCount++
	}
	return noteCount
}

func (book *Notebook) LinkCount() int {
	return len(book.Notes.Edges) / 2
}

func (book *Notebook) Tags(search string) []string {
	var tags []string
	tagConnections := book.tagConnections(search)
	for _, tagConnection := range tagConnections {
		tags = append(tags, tagConnection.Tag)
	}
	sort.Strings(tags)
	return tags
}

type TagConnection struct {
	Tag         string
	Connections int
}

func (book *Notebook) tagConnections(search string) []TagConnection {
	var tagConnections []TagConnection
	traversal := graph.Traversal(book.Notes)
	for vertex := range traversal.V().HasLabel("tag").Iterate() {
		tag := fmt.Sprint(vertex.Properties["Value"])
		if util.Matches(search, tag) {
			tagConnection := TagConnection{Tag: tag, Connections: len(book.Notes.Adjacency[vertex.Id])}
			tagConnections = append(tagConnections, tagConnection)
		}
	}
	return tagConnections
}

func (book *Notebook) TagConnections(search string) []TagConnection {
	tagConnections := book.tagConnections(search)
	sort.SliceStable(tagConnections, func(i, j int) bool {
		return tagConnections[i].Connections > tagConnections[j].Connections
	})
	return tagConnections
}

func (book *Notebook) MatchedTags(searchTags ...string) []matchedTag {
	tagVertices := make(map[uint64]matchedTag)
	traversal := graph.Traversal(book.Notes)
	for vertex := range traversal.V().HasLabel("tag").Iterate() {
		tag := fmt.Sprint(vertex.Properties["Value"])
		for _, searchTag := range searchTags {
			if util.Matches(searchTag, tag) {
				tagVertices[vertex.Id] = matchedTag{
					Distance: util.Distance(searchTag, tag),
					Tag:      tag,
					Vertex:   vertex,
				}
			}
		}
	}

	if len(tagVertices) == 0 {
		return nil
	}

	tagVerticesSlice := make([]matchedTag, 0, len(tagVertices))
	for _, vertex := range tagVertices {
		tagVerticesSlice = append(tagVerticesSlice, vertex)
	}

	sort.SliceStable(tagVerticesSlice, func(i, j int) bool {
		if tagVerticesSlice[i].Distance == tagVerticesSlice[j].Distance {
			return len(book.Notes.Adjacency[tagVerticesSlice[i].Vertex.Id]) > len(book.Notes.Adjacency[tagVerticesSlice[j].Vertex.Id])
		} else {
			return tagVerticesSlice[i].Distance < tagVerticesSlice[j].Distance
		}
	})

	if len(tagVerticesSlice) > len(searchTags) {
		tagVerticesSlice = tagVerticesSlice[:len(searchTags)]
	}

	return tagVerticesSlice
}

func (book *Notebook) TagIntersection(matchedTags []matchedTag) []graph.Vertex {
	var intersection = make(map[uint64]graph.Vertex)

	if len(matchedTags) == 0 {
		return nil
	}

	var mostConnectedVertex = matchedTags[0]
VERTEX:
	for vertexId := range book.Notes.Adjacency[mostConnectedVertex.Vertex.Id] {
		for _, tagVertex := range matchedTags[1:] {
			if !book.Notes.IsEdge(tagVertex.Vertex.Id, vertexId) {
				continue VERTEX
			}
		}
		intersection[vertexId] = book.Notes.Vertices[vertexId]
	}

	var sortedVertices []uint64
	for vertexId := range intersection {
		sortedVertices = append(sortedVertices, vertexId)
	}
	sort.SliceStable(sortedVertices, func(i, j int) bool {
		return sortedVertices[i] < sortedVertices[j]
	})

	var sortedIntersection []graph.Vertex
	for _, vertexId := range sortedVertices {
		sortedIntersection = append(sortedIntersection, intersection[vertexId])
	}

	return sortedIntersection
}

func (book *Notebook) addTag(tag string) uint64 {
	tagId, ok := book.tags[tag]
	if !ok {
		tagNum := len(book.tags)
		tagId = uint64(tagNum + 1)
		book.tags[tag] = tagId
	}

	return tagId
}
