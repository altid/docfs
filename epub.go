package main

import (
	"io"
	"log"
	"path"
	"os"

	"github.com/meskio/epubgo"
	"github.com/ubqt-systems/cleanmark"
	"github.com/ubqt-systems/fs"
)

func parseEpubTitle(c *fs.Control, docname string, r *epubgo.Epub) {
	w := c.TitleWriter(docname)
	title := cleanmark.NewCleaner(w)
	defer title.Close()
	t, _:= r.Metadata("title")
	title.WriteStringEscaped(t[0])
}

func parseEpubSidebar(c *fs.Control, docname string, r *epubgo.Epub) {
	var n int
	entries := make(chan entry)
	defer close(entries)
	go writeOutline(c, docname, entries)
	it, err := r.Navigation()
	if err != nil {
		log.Print(err)
		return
	}
	for {
		entries <- entry{
			len: n,
			url: []byte(it.URL()),
			msg: []byte(it.Title()),
		}
		switch {
		case it.HasChildren():
			it.In()
			n++
		case it.IsLast():
			err = it.Out()
			n--
		}
		err = it.Next()
		if err != nil {
			break
		}
	}
}

func parseEpubBody(c *fs.Control, docname string, r *epubgo.Epub) error {
	// Iterate through spine elements, and convert html to our markdown
	it, err := r.Spine()
	if err != nil {
		return err
	}
	w := c.MainWriter(docname, "document")
	body := cleanmark.NewHTMLCleaner(w)
	defer body.Close()
	for {
		content, err := it.Open()
		if err != nil {
			log.Print(err)
			continue
		}
		err = body.Parse(content)
		if err != io.EOF {
			return err
		}
		err = it.Next()
		if err != nil {
			break
		}
		body.WriteString("\n")
	}
	return nil
}

func parseEpub(c *fs.Control, newfile string) error {
	docname := path.Base(newfile)
	docdir := path.Join(*mtpt, "docs", docname)
	if _, err := os.Stat(docdir); os.IsNotExist(err) {
		os.MkdirAll(docdir, 0755)
	}
	w := c.StatusWriter(docname)
	status := cleanmark.NewCleaner(w)
	defer status.Close()
	pages, err := epubgo.Open(newfile)
	if err != nil {
		status.WriteString("Error opening file. See log for details.")
		return err
	}
	defer pages.Close()
	status.WriteString("Parsing file...")
	parseEpubTitle(c, docname, pages)
	parseEpubSidebar(c, docname, pages)
	err = parseEpubBody(c, docname, pages)
	if err != nil {
		status.WriteString("Error parsing file. See log for details.")
		return err
	}
	return c.Remove(docname, "status")
}
