package session

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/altid/libs/markup"
	"github.com/altid/libs/service/controller"
)

type entry struct {
	len int
	msg []byte
	url []byte
}

func doParse(c controller.Controller, newfile, mime string) error {
	switch mime {
	case "application/pdf":
		return parsePdf(c, newfile)
	case "application/epub", "application/zip":
		return parseEpub(c, newfile)
	}
	return fmt.Errorf("unsupported document type requested: %q", mime)
}

func mimeFromContents(c controller.Controller, newfile string) error {
	buf := make([]byte, 512)
	fp, err := os.Open(newfile)
	if err != nil {
		return err
	}
	defer fp.Close()
	n, err := fp.Read(buf)
	if err != nil {
		return err
	}
	mime := http.DetectContentType(buf[:n])
	return doParse(c, newfile, mime)
}

func parseDocument(c controller.Controller, newfile string) error {
	ext := filepath.Ext(newfile)
	mime := mime.TypeByExtension(ext)
	if mime != "" {
		return doParse(c, newfile, mime)
	}
	return mimeFromContents(c, newfile)
}

func writeOutline(c controller.Controller, docname string, entries chan entry) error {
	w, err := c.SideWriter(docname)
	if err != nil {
		return err
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
	return nil
}