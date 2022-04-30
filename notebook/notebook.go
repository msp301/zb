package notebook

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/parser"
)

type Notebook struct {
	Filters   []FilterFunc
	lookup    map[uint64]parser.Note
	linkGraph *graph.Graph
	tags      map[string]uint64
	tagGraph  *graph.Graph
	Path      string
	Notes     []parser.Note
}

type FilterFunc func(note parser.Note) bool

func New(path string) Notebook {
	return Notebook{
		Path:      path,
		lookup:    map[uint64]parser.Note{},
		linkGraph: graph.New(),
		tags:      map[string]uint64{},
		tagGraph:  graph.New(),
	}
}

func (book *Notebook) AddFilter(filter FilterFunc) {
	book.Filters = append(book.Filters, filter)
}

func (book *Notebook) Read() []parser.Note {
	var filteredNotes []parser.Note

	filepath.Walk(book.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		fileNotes := parser.Parse(path)
		book.Notes = append(book.Notes, fileNotes...)

	NOTE:
		for _, note := range fileNotes {
			book.lookup[note.Id] = note
			book.linkGraph.AddVertex(note.Id)

			for _, link := range note.Links {
				book.linkGraph.AddEdge(note.Id, link)
			}

			for _, tag := range note.Tags {
				tagId := book.addTag(tag)
				book.tagGraph.AddVertex(tagId)
			}

			for _, tag := range note.Tags {
				for _, relatedTag := range note.Tags {
					if relatedTag == tag {
						continue
					}
					book.tagGraph.AddEdge(book.tags[tag], book.tags[relatedTag])
				}
			}

			for _, filter := range book.Filters {
				if filter(note) {
					continue NOTE
				}
			}
			filteredNotes = append(filteredNotes, note)
		}

		return nil
	})

	tagKeys := map[uint64]string{}
	for key, val := range book.tags {
		tagKeys[val] = key
	}

	tagMap := map[string][]string{}
	for tag, relatedTags := range book.tagGraph.Edges {
		for _, relatedTag := range relatedTags {
			tagStr, ok := tagKeys[tag]
			if ok {
				tagMap[tagStr] = append(tagMap[tagStr], tagKeys[relatedTag])
			}
		}
	}

	//fmt.Println(book.linkGraph.Edges)
	//fmt.Println(book.tagGraph.Edges)
	fmt.Println(book.tagGraph.Vertices)
	fmt.Println(tagMap)

	return filteredNotes
}

func (book *Notebook) GetNote(noteId uint64) (parser.Note, bool) {
	note, ok := book.lookup[noteId]
	return note, ok
	//note, ok := book.lookup.GetVertex(noteId)
	//return note, ok
}

func (book *Notebook) addTag(tag string) uint64 {
	tagNum := len(book.tags)
	tagId := uint64(tagNum + 1)
	book.tags[tag] = tagId
	return tagId
}
