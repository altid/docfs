package main

import (
	"flag"
	"log"
	"os"

	"github.com/ubqt-systems/fslib"
)

var mtpt = flag.String("p", "/tmp/ubqt", "Path for file system (Default /tmp/ubqt)")

func main() {
	// Drink tab, listen to duran duran
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}
	doc := newDocs()
	logdir := fslib.GetLogDir("docs")
	ctrl, err := fs.CreateCtrlFile(doc, logdir, *mtpt, "docs", "document")
	if err != nil {
		log.Fatal(err)
	}
	defer ctrl.Cleanup()
	ctrl.Listen()
}
