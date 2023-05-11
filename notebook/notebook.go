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
		context := ""

		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			context = util.Context(val.Content, fmt.Sprint(id))
		}

		result := Result{
			Context: context,
			Value:   vertex,
		}
		results = append(results, result)
	}

	return results
}

type Result struct {
	Context string
	Value   interface{}
}

func (book *Notebook) SearchByTag(searchTag string) []Result {
	var tagVertex graph.Vertex
	traversal := graph.Traversal(book.Notes)
	for vertex := range traversal.V().HasLabel("tag").Iterate() {
		tag := fmt.Sprint(vertex.Properties["Value"])
		if util.Matches(searchTag, tag) {
			tagVertex = vertex
			break
		}
	}

	var results []Result
	for id := range book.Notes.Adjacency[tagVertex.Id] {
		vertex := book.Notes.Vertices[id]
		context := ""

		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			tag := fmt.Sprint(tagVertex.Properties["Value"])
			context = util.Context(val.Content, tag)
		}

		result := Result{
			Context: context,
			Value:   vertex,
		}
		results = append(results, result)
	}

	return results
}

type matchedTag struct {
	Distance int
	Tag      string
	Vertex   graph.Vertex
}

func (book *Notebook) SearchByTags(searchTags []string) []Result {
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

	tagVerticesSlice = tagVerticesSlice[:len(searchTags)]

	var results []Result
	var intersection = make(map[uint64]graph.Vertex)
	var mostConnectedVertex = tagVerticesSlice[0]
VERTEX:
	for vertexId := range book.Notes.Adjacency[mostConnectedVertex.Vertex.Id] {
		for _, tagVertex := range tagVerticesSlice[1:] {
			if !book.Notes.IsEdge(tagVertex.Vertex.Id, vertexId) {
				continue VERTEX
			}
		}
		intersection[vertexId] = book.Notes.Vertices[vertexId]
	}

	for _, vertex := range intersection {
		context := ""

		switch val := vertex.Properties["Value"].(type) {
		case parser.Note:
			matchedTag := ""
			for _, tagVertex := range tagVerticesSlice {
				if book.Notes.IsEdge(tagVertex.Vertex.Id, vertex.Id) {
					matchedTag = fmt.Sprint(tagVertex.Vertex.Properties["Value"])
					break
				}
			}
			context = util.Context(val.Content, matchedTag)
		}

		result := Result{
			Context: context,
			Value:   vertex,
		}
		results = append(results, result)
	}

	return results
}

func (book *Notebook) Tags(search string) []string {
	var tags []string
	traversal := graph.Traversal(book.Notes)
	for _, tag := range traversal.V().HasLabel("tag").Values("Value") {
		tagStr := fmt.Sprint(tag)
		if util.Matches(search, tagStr) {
			tags = append(tags, tagStr)
		}
	}
	sort.Strings(tags)
	return tags
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
