package main

import (
	"fmt"
	"os"
	"path"

	"github.com/ubqt-systems/cleanmark"
	"github.com/ubqt-systems/fs"
	"rsc.io/pdf"
)

func parsePdfBody(c *fs.Control, docname string, r *pdf.Reader) error {
	// test if logged document exists and is not empty
	// exit early
	// loop writetomain.
	// TODO: v2 set up logdir/resources/images/mydoc.pdf/, and link that to our local mydoc.pdf/images/
	// TODO: v2 We want anchor links when they exist
	// We want... ANCHOR LINKS WHEN THEY EXIST
	// Initially, we'll just translate the page to simple text
	// eventually we'll properly convert it to markdown.
	numPages := r.NumPage()
	fmt.Println(numPages)
	w := c.MainWriter(docname, "document")
	body := cleanmark.NewCleaner(w)
	defer body.Close()
	for i := 1; i <= numPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		// Alright buf[y][x] - then we right the buf[y] when y increases
		for _, t := range page.Content().Text {
			fmt.Print(t.S)
		}
		//body.WriteStringEscaped(page.V.RawString())
 	}
	return nil
}

func parsePdfTitle(c *fs.Control, docname string, r *pdf.Reader) {
	w := c.TitleWriter(docname)
	title := cleanmark.NewCleaner(w)
	defer title.Close()
	tstring := r.Outline().Title
	if tstring == "" {
		tstring = docname
	}
	title.WriteStringEscaped(tstring)
}

func parsePdfSidebar(c *fs.Control, docname string, r *pdf.Reader) {
	entries := make(chan entry)
	go func(entries chan entry) {
		w := c.SideWriter(docname)
		sidebar := cleanmark.NewCleaner(w)
		defer sidebar.Close()
		for e := range entries {
			sidebar.WriteList(e.len, e.msg)
		}
	}(entries)
	// Skip the document title - do the first walk here
	for _, item := range r.Outline().Child {
		walkPdfOutline(item, 0, entries)
	}
	close(entries)
}

func walkPdfOutline(r pdf.Outline, n int, entries chan entry) {
	if r.Title != "" {
		entries <- entry{
			len: n,
			msg: []byte(r.Title),
		}
	}
	n++
	for _, item := range r.Child {
		walkPdfOutline(item, n, entries)
	}
}

func parsePdf(c *fs.Control, newfile string) error {
	docname := path.Base(newfile)
	docdir := path.Join(*mtpt, "docs", docname)
	if _, err := os.Stat(docdir); os.IsNotExist(err) {
		os.MkdirAll(docdir, 0755)
	}
	w := c.StatusWriter(docname)
	status := cleanmark.NewCleaner(w)
	defer status.Close()
	pages, err := pdf.Open(newfile)
	if err != nil {
		status.WriteString("Error opening file. See log for details.")
		return err
	}
	status.WriteString("Parsing file...")
	parsePdfTitle(c, docname, pages)
	parsePdfSidebar(c, docname, pages)
	err = parsePdfBody(c, docname, pages)
	if err != nil {
		status.WriteString("Error parsing file. See log for details.")
		return err
	}
	return c.Remove(docname, "status")
}
