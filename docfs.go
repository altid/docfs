package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/altid/libs/config"
	"github.com/altid/libs/fs"
)

var mtpt = flag.String("p", "/tmp/altid", "Path for file system")
var debug = flag.Bool("d", false, "enable debug logging")
var srv = flag.String("s", "docs", "Name of service")
var setup = flag.Bool("conf", false, "Set up configuration file")

func main() {
	// Drink tab, listen to duran duran
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	conf := &struct {
		Log  config.Logdir `Use directory for cached files (none to disable)`
		Addr config.ListenAddress
	}{"none", ""}

	if *setup {
		if e := config.Create(conf, *srv, "", *debug); e != nil {
			log.Fatal(e)
		}

		os.Exit(0)
	}

	if e := config.Marshal(conf, *srv, "", *debug); e != nil {
		log.Fatal(e)
	}

	ctx, cancel := context.WithCancel(context.Background())
	doc := &docs{cancel}

	ctrl, err := fs.CreateCtlFile(ctx, doc, string(conf.Log), *mtpt, *srv, "document", *debug)
	if err != nil {
		log.Fatal(err)
	}

	ctrl.CreateBuffer("welcome", "document")
	wc, err := ctrl.MainWriter("welcome", "document")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(wc, welcome)
	wc.Close()

	defer ctrl.Cleanup()
	ctrl.Listen()
}
