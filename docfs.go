package main

import (
	"errors"
	"mime"
	"os"
	"path"

	"rsc.io/pdf"
)

func parseDocument(filename string) error {
	mime.TypeByExtension(filename)
	// filenames with slashes inthem should be supported
	docdir := path.Join(*mtpt, filename)
	if _, err := os.Stat(docdir); os.IsNotExist(err) {
		os.MkdirAll(docdir), 0755)
	}
	switch doctype {
	case "application/pdf":
		pages, err := pdf.Open(filename)
		if err != nil {
			return err
		}
		title, err := os.Create(path.Join(docdir, "title"))
		defer title.Close()
		if err != nil {
			return err
		}
		title.WriteString(r.Outline.Title)
		// Walk tree, generate sidebar (ToC)
		/* sidebar, err := os.Create(path.Join(docdir, "sidebar"))
		defer sidebar.Close()
		if err != nil {
			return err
		} */
	default:
		errors.New("Unsupported document type requested. PRs welcome!")
	}
}

func createDocument(filename string) {
	logPath := path.Join(*log, filename)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		fp, _ := os.Create(path.Join(*log, filename))
		fp.Close()
	}
	switch runtime.GOOS {
	case "plan9":
		command := exec.Command("/bin/bind", logPath, path.Join(*mtpt, filename, "document"))
		err := command.Run()
		if err != nil {
			return err
		}
	default:
		err := os.Symlink(logPath, path.Join(*mtpt, filename, "document")
	}
}

func newDocument(req chan string, done chan struct{}) {
	for {
		select {
		case filename := <-req:
			createDocument(filename)
			parseDocument(filename)	
		case <-done:
			break
		}
	}
}

func main() {
	// Drink tab, listen to duran duran
	// ctrl

	// event

	// tabs
}