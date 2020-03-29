package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/altid/libs/config"
	"github.com/altid/libs/config/types"
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
		Log    types.Logdir        `altid:"directory to log files to,no_prompt"`
		Listen types.ListenAddress `altid:"listen address to use,omit_empty"`
	}{"none", "none"}

	if *setup {
		if e := config.Create(conf, *srv, "", *debug); e != nil {
			log.Fatal(e)
		}

		os.Exit(0)
	}

	if e := config.Marshal(conf, *srv, "", *debug); e != nil {
		log.Fatal(e)
	}

	ctrl, err := fs.New(&docs{}, string(conf.Log), *mtpt, *srv, "document", *debug)
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
