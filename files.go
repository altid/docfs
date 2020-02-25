package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/altid/libs/fs"
	"github.com/altid/libs/markup"
)

type entry struct {
	len int
	msg []byte
	url []byte
}

// TODO: Test navi
func writeOutline(c *fs.Control, docname string, entries chan entry) {
	w, err := c.SideWriter(docname)
	if err != nil {
		log.Fatal(err)
	}

	aside := markup.NewCleaner(w)
	defer aside.Close()

	for e := range entries {
		url, err := markup.NewUrl(e.url, e.msg)
		if err != nil {
			continue
		}
		aside.WritefList(e.len, "%s\n", url)
	}

}

type docs struct{}

func newDocs() *docs {
	return &docs{}
}

func (d *docs) Open(c *fs.Control, newfile string) error {
	c.CreateBuffer(path.Base(newfile), "document")

	err := parseDocument(c, newfile)
	if err != nil {
		c.DeleteBuffer(path.Base(newfile), "document")
		return err
	}

	return nil
}

func (d *docs) Close(c *fs.Control, newfile string) error {
	c.DeleteBuffer(path.Base(newfile), "document")
	return nil
}

func (d *docs) Link(c *fs.Control, from, newfile string) error {
	c.DeleteBuffer(path.Base(newfile), "document")
	return d.Open(c, newfile)
}

func (d *docs) Default(c *fs.Control, cmd, from, msg string) error {
	return nil
}

func doParse(c *fs.Control, newfile, mime string) error {
	switch mime {
	case "application/pdf":
		return parsePdf(c, newfile)
	case "application/epub", "application/zip":
		return parseEpub(c, newfile)
	}
	return fmt.Errorf("Unsupported document type requested: %q. PRs welcome!", mime)
}

func mimeFromContents(c *fs.Control, newfile string) error {
	buf := make([]byte, 512)
	fp, err := os.Open(newfile)
	defer fp.Close()
	if err != nil {
		return err
	}
	n, err := fp.Read(buf)
	if err != nil {
		return err
	}
	mime := http.DetectContentType(buf[:n])
	return doParse(c, newfile, mime)
}

func parseDocument(c *fs.Control, newfile string) error {
	ext := filepath.Ext(newfile)
	mime := mime.TypeByExtension(ext)
	if mime != "" {
		return doParse(c, newfile, mime)
	}
	return mimeFromContents(c, newfile)
}
