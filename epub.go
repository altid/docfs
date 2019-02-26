package main

import (
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
	entries := make(chan entry)
	go func(entries chan entry) {
		w := c.SideWriter(docname)
		sidebar := cleanmark.NewCleaner(w)
		defer sidebar.Close()
		for e := range entries {
			sidebar.WriteList(e.len, e.msg)
		}
	}(entries)		
	//it, err := r.Navigation()
	// walk it.Next()'s. We have to check depth here, and In() or Out() accordingly
	// We should check HasChildren to see if we go In, IsLast() to see if we go out, etc.
	// it.Title()
	// it.Next()
	// it.Url()
	// TODO v2: We want to set up proper anchor links here, and in the main body.
	close(entries)	
}

func parseEpubBody(c *fs.Control, docname string, r *epubgo.Epub) error {
	// Then we iterate through spine elements, and read it all in. 
	// It's likely there will be links to `files` in the pdf, which are images, and other shit. So we'll find this too and add it to our resources.
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
