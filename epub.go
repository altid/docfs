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
func (h *helper) Nav(u *html.URL) error {
	return nil
}

func (h *helper) Img(img *html.Image) error {
	var rc io.ReadCloser
	var err error

	src := string(img.Src)
	dirs := []string{
		"OEBPS/" + src,
		"OPS/" + src,
		"EPUB/" + src,
		src,
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
	wc, err := h.c.ImageWriter(h.dir, src)
	if err != nil {
		return err
	}

	defer wc.Close()
	_, err = io.Copy(wc, rc)

	return err
}

func parseEpubNavi(c *fs.Control, docname string, r *epubgo.Epub) error {
	var n int
	w, err := c.NavWriter(docname)
	if err != nil {
		return err
	}

	navi := markup.NewCleaner(w)
	defer navi.Close()

	it, err := r.Navigation()
	if err != nil {
		return err
	}

	for {
		url, _ := markup.NewURL([]byte(it.URL()), []byte(it.Title()))
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

func parseEpubTitle(c *fs.Control, docname string, r *epubgo.Epub) error {
	w, err := c.TitleWriter(docname)
	if err != nil {
		return err
	}

	title := markup.NewCleaner(w)
	defer title.Close()
	t, _ := r.Metadata("title")
	title.WriteStringEscaped(t[0])

	return nil
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
	w, err := c.MainWriter(docname, "document")
	if err != nil {
		return err
	}

	body, err := html.NewCleaner(w, h)
	if err != nil {
		return err
	}

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

	sw, err := c.StatusWriter(docname)
	if err != nil {
		log.Fatal(err)
	}

	status := markup.NewCleaner(sw)
	defer status.Close()

	// TODO: EPUB3 currently aren't supported by https://github.com/meskio/epubgo
	// https://github.com/altid/docfs/issues/4
	pages, err := epubgo.Open(newfile)
	if err != nil {
		return err
	}
	
	defer pages.Close()
	status.WriteString("Parsing file...")

	if e := parseEpubTitle(c, docname, pages); e != nil {
		return e
	}

	if e := parseEpubNavi(c, docname, pages); e != nil {
		return e
	}

	if e := parseEpubBody(c, docname, pages); e != nil {
		return e
	}

	return c.Remove(docname, "status")
}
