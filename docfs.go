package docfs

import (
	"context"

	"github.com/altid/docfs/internal/commands"
	"github.com/altid/docfs/internal/session"
	"github.com/altid/libs/config"
	"github.com/altid/libs/mdns"
	"github.com/altid/libs/service"
	"github.com/altid/libs/service/listener"
	"github.com/altid/libs/store"
)

type Docfs struct {
	run     func() error
	session *session.Session
	name    string
	addr    string
	debug   bool
	mdns    *mdns.Entry
	ctx     context.Context
}

var defaults *session.Defaults = &session.Defaults{
	Logdir:  "",
	TLSCert: "",
	TLSKey:  "",
}

func CreateConfig(srv string, debug bool) error {
	return config.Create(defaults, srv, "", debug)
}

func Register(ldir bool, addr, srv string, debug bool) (*Docfs, error) {
	if e := config.Marshal(defaults, srv, "", debug); e != nil {
		return nil, e
	}
	l, err := tolisten(defaults, addr, debug)
	if err != nil {
		return nil, err
	}
	s := tostore(defaults, ldir, debug)
	session := &session.Session{
		Defaults: defaults,
		Verbose:  debug,
	}

	session.Parse()
	ctx := context.Background()

	d := &Docfs{
		session: session,
		ctx:     ctx,
		name:    srv,
		addr:    addr,
		debug:   debug,
	}

	c := service.New(srv, addr, debug)
	c.WithListener(l)
	c.WithStore(s)
	c.WithContext(ctx)
	c.WithCallbacks(session)
	c.WithRunner(session)

	c.SetCommands(commands.Commands)
	d.run = c.Listen

	return d, nil
}

func (doc *Docfs) Run() error {
	return doc.run()
}

func (doc *Docfs) Broadcast() error {
	entry, err := mdns.ParseURL(doc.addr, doc.name)
	if err != nil {
		return err
	}
	if e := mdns.Register(entry); e != nil {
		return e
	}
	doc.mdns = entry
	return nil
}

func (doc *Docfs) Cleanup() {
	if doc.mdns != nil {
		doc.mdns.Cleanup()
	}
	doc.session.Quit()
}

func (doc *Docfs) Session() *session.Session {
	return doc.session
}

func tolisten(d *session.Defaults, addr string, debug bool) (listener.Listener, error) {
	//if ssh {
	//    return listener.NewListenSsh()
	//}

	if d.TLSKey == "none" && d.TLSCert == "none" {
		return listener.NewListen9p(addr, "", "", debug)
	}

	return listener.NewListen9p(addr, d.TLSCert, d.TLSKey, debug)
}

func tostore(d *session.Defaults, ldir, debug bool) store.Filer {
	if ldir {
		return store.NewLogstore(d.Logdir.String(), debug)
	}

	return store.NewRamstore(debug)
}
