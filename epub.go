package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/meskio/epubgo"
	"github.com/altid/cleanmark"
	fs "github.com/altid/fslib"
)

var findAssets = regexp.MustCompile(`<img\s+[^>]*?src=["']([^"']+)`)

func parseEpubTitle(c *fs.Control, docname string, r *epubgo.Epub) {
	w := c.TitleWriter(docname)
	title := cleanmark.NewCleaner(w)
	defer title.Close()
	t, _ := r.Metadata("title")
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

// Walk the document a second time and pull out any assets linked in
func parseAssets(c *fs.Control, docname string, r *epubgo.Epub) error {
	it, err := r.Spine()
	if err != nil {
		return err
	}
	for {
		content, err := it.Open()
		reader := bufio.NewReader(content)
		if err = it.Next(); err != nil {
			return err
		}
		for {
			buff, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			findPageAssets(buff, r, docname)
		}
	}
	return nil
}

func findPageAssets(s string, r *epubgo.Epub, docname string) {
	matches := findAssets.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		var a io.ReadCloser
		var err error
		dirs := []string{
			"OEBPS/" + match[1],
			"OPS/" + match[1],
			"EPUB/" + match[1],
			match[1],
		}
		for _, try := range dirs{
			a, err = r.OpenFile(try)
			if err == nil {
				break
			}
		}
		if err != nil {
			continue
		}
		defer a.Close()
		fp := path.Join(*mtpt, "docs", docname, match[1])
		os.MkdirAll(path.Dir(fp), 0755)
		output, err := os.Create(fp)
		if err != nil {
			continue
		}
		defer output.Close()
		io.Copy(output, a)
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
		if err = it.Next(); err != nil {
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

	// TODO: EPUB3 currently aren't supported by https://github.com/meskio/epubgo
	// https://github.com/altid/docfs/issues/4
	pages, err := epubgo.Open(newfile)
	if err != nil {
		status.WriteString("Error opening file. See log for details.")
		return err
	}
	defer pages.Close()
	status.WriteString("Parsing file...")
	parseAssets(c, docname, pages)
	parseEpubTitle(c, docname, pages)
	parseEpubSidebar(c, docname, pages)
	err = parseEpubBody(c, docname, pages)
	if err != nil {
		status.WriteString("Error parsing file. See log for details.")
		return err
	}
	return c.Remove(docname, "status")
}
