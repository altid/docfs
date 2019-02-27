package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/ubqt-systems/cleanmark"
	"github.com/ubqt-systems/fs"
)

type entry struct {
	len int
	msg []byte
	url []byte
}

func writeOutline(c *fs.Control, docname string, entries chan entry) {
	w := c.SideWriter(docname)
	sidebar := cleanmark.NewCleaner(w)
	defer sidebar.Close()
	for e := range entries {
		url := cleanmark.NewUrl(e.url, e.msg)
		sidebar.WritefList(e.len, "%s\n", url)
	}
}

type docs struct {
}

func newDocs() *docs {
	return &docs{}
}

func (d *docs) Open(c *fs.Control, newfile string) error {
	c.CreateBuffer(path.Base(newfile), "document")
	err := parseDocument(c, newfile)
	if err != nil {
		log.Printf("%v in Open function call", err)
		c.DeleteBuffer(path.Base(newfile), "document")
		log.Print(err)
		return err
	}
	return nil
}

func (d *docs) Close(c *fs.Control, newfile string) error {
	c.DeleteBuffer(path.Base(newfile), "document")
	return nil
}

func (d *docs) Default(c *fs.Control, line string) error {
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
	if err != nil {
		return err
	}
	n, err := fp.Read(buf)
	if err != nil {
		return err
	}
	mime :=  http.DetectContentType(buf[:n])
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
