package notebook

import (
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/parser"
	"github.com/msp301/zb/util"
)

type Notebook struct {
	Filters  []FilterFunc
	lookup   map[uint64]parser.Note
	Notes    *graph.Graph
	tags     map[string]uint64
	tagGraph *graph.Graph
	Path     string
	notes    []parser.Note
}

type FilterFunc func(note parser.Note) bool

func New(path string) *Notebook {
	return &Notebook{
		Path:     path,
		lookup:   map[uint64]parser.Note{},
		Notes:    graph.New(),
		tags:     map[string]uint64{},
		tagGraph: graph.New(),
	}
}

func (book *Notebook) AddFilter(filter FilterFunc) {
	book.Filters = append(book.Filters, filter)
}

func (book *Notebook) Read() []parser.Note {
	var filteredNotes []parser.Note
	var noteFiles []string

	filepath.WalkDir(book.Path, func(path string, d fs.DirEntry, err error) error {
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

	for _, path := range noteFiles {
		fileNotes := parser.Parse(path)
		book.notes = append(book.notes, fileNotes...)

	NOTE:
		for _, note := range fileNotes {
			book.lookup[note.Id] = note
			book.Notes.AddVertex(graph.Vertex{Id: note.Id, Label: "note", Properties: note})

			for _, link := range note.Links {
				book.Notes.AddEdge(note.Id, link)
			}

			for _, tag := range note.Tags {
				tagId := book.addTag(tag)
				book.Notes.AddVertex(graph.Vertex{Id: tagId, Label: "tag", Properties: tag})
			}

			for _, tag := range note.Tags {
				if !book.Notes.IsVertex(book.tags[tag]) {
					book.Notes.AddVertex(graph.Vertex{Id: book.tags[tag], Label: "tag", Properties: tag})
				}

				for _, relatedTag := range note.Tags {
					if relatedTag == tag {
						continue
					}
					if !book.Notes.IsVertex(book.tags[relatedTag]) {
						book.Notes.AddVertex(graph.Vertex{Id: book.tags[relatedTag], Label: "tag", Properties: tag})
					}
					book.Notes.AddEdge(book.tags[tag], book.tags[relatedTag])
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

	tagKeys := map[uint64]string{}
	for key, val := range book.tags {
		tagKeys[val] = key
	}

	//tagMap := map[string][]string{}
	//for tag, relatedTags := range book.tagGraph.Edges {
	//	for _, relatedTag := range relatedTags {
	//		tagStr, ok := tagKeys[tag]
	//		if ok {
	//			tagMap[tagStr] = append(tagMap[tagStr], tagKeys[relatedTag])
	//		}
	//	}
	//}

	return filteredNotes
}

func (book *Notebook) IsNote(noteId uint64) bool {
	_, ok := book.Notes.Vertices[noteId]
	return ok
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
