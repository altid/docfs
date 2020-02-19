package main

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/altid/libs/fs"
	"github.com/altid/libs/html"
	"github.com/altid/libs/markup"
	"github.com/meskio/epubgo"
)

type helper struct {
	c   *fs.Control
	r   *epubgo.Epub
	dir string
}

// Maybe EPUB 3 uses <nav> properly?
func (h *helper) Nav(u *markup.Url) error {
	return nil
}

func (h *helper) Img(link string) error {
	log.Println(link)
	var rc io.ReadCloser
	var err error
	dirs := []string{
		"OEBPS/" + link,
		"OPS/" + link,
		"EPUB/" + link,
		link,
	}
	for _, try := range dirs {
		rc, err = h.r.OpenFile(try)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	defer rc.Close()
	wc := h.c.ImageWriter(h.dir, link)
	defer wc.Close()
	_, err = io.Copy(wc, rc)
	return err
}

func parseEpubNavi(c *fs.Control, docname string, r *epubgo.Epub) error {
	var n int
	w := c.NavWriter(docname)
	navi := markup.NewCleaner(w)
	defer navi.Close()
	it, err := r.Navigation()
	if err != nil {
		return err
	}
	for {
		url, _ := markup.NewUrl([]byte(it.URL()), []byte(it.Title()))
		navi.WritefList(n, "%s\n", url)
		switch {
		case it.HasChildren():
			it.In()
			n++
		case it.IsLast():
			it.Out()
			n--
		}
		err = it.Next()
		if err != nil {
			return err
		}
	}
}

func parseEpubTitle(c *fs.Control, docname string, r *epubgo.Epub) {
	w := c.TitleWriter(docname)
	title := markup.NewCleaner(w)
	defer title.Close()
	t, _ := r.Metadata("title")
	title.WriteStringEscaped(t[0])
}

func parseEpubBody(c *fs.Control, docname string, r *epubgo.Epub) error {
	h := &helper{
		c:   c,
		r:   r,
		dir: docname,
	}
	it, err := r.Spine()
	if err != nil {
		return err
	}
	w := c.MainWriter(docname, "document")
	body := html.NewHTMLCleaner(w, h)
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
	status := markup.NewCleaner(w)
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
	parseEpubTitle(c, docname, pages)
	parseEpubNavi(c, docname, pages)
	err = parseEpubBody(c, docname, pages)
	if err != nil {
		status.WriteString("Error parsing file. See log for details.")
		return err
	}
	return c.Remove(docname, "status")
}
