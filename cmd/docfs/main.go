package main

import (
	"flag"
	"log"
	"os"

	"github.com/altid/docfs"
)

var (
	srv		= flag.String("s", "docs", "name of service")
	addr 	= flag.String("a", "127.0.0.1:12345", "listening address")
	mdns	= flag.Bool("m", false, "enable mDNS broadcast of service")
	debug 	= flag.Bool("d", false, "enable debug printing")
	ldir	= flag.Bool("l", false, "enable logging for main buffers")
	setup	= flag.Bool("conf", false, "run configuration setup")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	if *setup {
		if e := docfs.CreateConfig(*srv, *debug); e != nil {
			log.Fatal(e)
		}
		os.Exit(0)
	}

	docs, err := docfs.Register(*ldir, *addr, *srv, *debug)
	if err != nil {
		log.Fatal(err)
	}

	defer docs.Cleanup()
	if *mdns {
		if e := docs.Broadcast(); e != nil {
			log.Fatal(e)
		}
	}

	if e := docs.Run(); e != nil {
		log.Fatal(e)
	}

	os.Exit(0)
}