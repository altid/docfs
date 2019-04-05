package main

import (
	"flag"
	"log"
	"os"

	"github.com/altid/fslib"
)

var mtpt = flag.String("p", "/tmp/altid", "Path for file system")

func main() {
	// Drink tab, listen to duran duran
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}
	doc := newDocs()
	logdir := fslib.GetLogDir("docs")
	ctrl, err := fslib.CreateCtrlFile(doc, logdir, *mtpt, "docs", "document")
	if err != nil {
		log.Fatal(err)
	}
	// TODO(halfwit): We want to create a default buffer - small how-to document would suffice
	// https://github.com/altid//docfs/issues/5
	defer ctrl.Cleanup()
	ctrl.Listen()
}
