package main

import (
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

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
		url, err := markup.NewURL(e.url, e.msg)
		if err != nil {
			continue
		}
		aside.WritefList(e.len, "%s\n", url)
	}

}

type docs struct {
	// No-op
}

func (d *docs) Run(c *fs.Control, cmd *fs.Command) error {
	switch cmd.Name {
	case "open":
		// Files can have whitespace
		// Commands sent in may be paths
		newfile := strings.Join(cmd.Args, " ")
		c.CreateBuffer(newfile, "document")

		err := parseDocument(c, newfile)
		if err != nil {
			c.DeleteBuffer(path.Base(newfile), "document")
			return err
		}

		return nil
	case "close":
		return c.DeleteBuffer(strings.Join(cmd.Args, " "), "document")
	default:
		return errors.New("Command not supported")
	}
	return nil
}

func (d *docs) Quit() {
	// No-op
}

func doParse(c *fs.Control, newfile, mime string) error {
	switch mime {
	case "application/pdf":
		return parsePdf(c, newfile)
	case "application/epub", "application/zip":
		return parseEpub(c, newfile)
	}
	return fmt.Errorf("unsupported document type requested: %q. PRs welcome!", mime)
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
