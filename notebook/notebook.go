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
	Filters []FilterFunc
	Notes   *graph.Graph
	tags    map[string]uint64
	Path    string
}

type FilterFunc func(note parser.Note) bool

func New(path string) *Notebook {
	return &Notebook{
		Path:  path,
		Notes: graph.New(),
		tags:  map[string]uint64{},
	}
}

func (book *Notebook) AddFilter(filter FilterFunc) {
	book.Filters = append(book.Filters, filter)
}

func (book *Notebook) Read() []parser.Note {
	var filteredNotes []parser.Note
	var noteFiles []string

	err := filepath.WalkDir(book.Path, func(path string, d fs.DirEntry, err error) error {
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
		noteFiles = append(noteFiles, path)

		return nil
	})
	if err != nil {
		return nil
	}

	for _, path := range noteFiles {
		fileNotes := parser.Parse(path)

	NOTE:
		for _, note := range fileNotes {
			book.Notes.AddVertex(graph.Vertex{Id: note.Id, Label: "note", Properties: map[string]interface{}{"Value": note}})

			for _, link := range note.Links {
				book.Notes.AddEdge(graph.Edge{
					Label: "link",
					From:  note.Id,
					To:    link,
				})
			}

			for _, tag := range note.Tags {
				tagId := book.addTag(tag)
				book.Notes.AddVertex(graph.Vertex{Id: tagId, Label: "tag", Properties: map[string]interface{}{"Value": tag}})
				book.Notes.AddEdge(graph.Edge{
					Label: "tag",
					From:  note.Id,
					To:    tagId,
				})
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
					book.Notes.AddEdge(graph.Edge{
						Label: "tag",
						From:  book.tags[tag],
						To:    book.tags[relatedTag],
					})
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

func (book *Notebook) SearchRelated(id uint64) []graph.Vertex {
	var vertices []graph.Vertex
	for vertexId := range book.Notes.Adjacency[id] {
		vertices = append(vertices, book.Notes.Vertices[vertexId])
	}
	return vertices
}

func (book *Notebook) SearchByTag(searchTag string) []graph.Vertex {
	var tagVertex graph.Vertex
	book.Notes.Walk(func(vertex graph.Vertex, depth int) bool {
		if vertex.Label != "tag" {
			return true
		}
		tag := fmt.Sprint(vertex.Properties["Value"])
		if tag == searchTag {
			tagVertex = vertex
		}
		return true
	})

	var results []graph.Vertex
	for id := range book.Notes.Adjacency[tagVertex.Id] {
		vertex := book.Notes.Vertices[id]
		results = append(results, vertex)
	}

	return results
}

func (book *Notebook) Tags() []string {
	var tags []string
	traversal := graph.Traversal(book.Notes)
	for _, tag := range traversal.V().HasLabel("tag").Values("Value") {
		tags = append(tags, fmt.Sprint(tag))
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
